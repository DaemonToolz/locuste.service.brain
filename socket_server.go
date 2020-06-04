package main

import (
	"log"
	"strings"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

// Note: Like the RPC section, plan and prepare a refactoring and shared modules
// locust.service.shared
// Divide the server into 2 websockets servers, one specialized in Automaton <=> Brain + Scheduler and Brain + Scheduler <=> UI
// Incoming steps : auto-scale websocket servers

var server *gosocketio.Server

// var droneEventServer map[string]*gosocketio.Server
// var droneControlServer  map[string]*gosocketio.Server

var channelMapping map[string]string

func initSocketServer() {
	channelMapping = make(map[string]string)
	server = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	// droneEventServer = map[string]*gosocketio.Server{ "shared":gosocketio.NewServer(transport.GetDefaultWebsocketTransport()) }
	// droneControlServer = map[string]*gosocketio.Server{ "shared":gosocketio.NewServer(transport.GetDefaultWebsocketTransport()) }

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Channel %s created", c.Id())
	})

	// droneEventServer["shared"]
	server.On("identify", func(c *gosocketio.Channel, request IdentificationRequest) {
		log.Printf("Channel %s identified as %s", c.Id(), request)
		c.Join(request.Name)
		channelMapping[c.Id()] = request.Name // On garde ça pour la déconnexion
		go server.BroadcastTo("operators", "identify", request)
		startVideoServer(request.Name, request.VideoPort)
		startFfmpegStream(request.Name, request.VideoPort)
	})

	//region Refactoring
	server.On("identify_operator", func(c *gosocketio.Channel) {
		log.Println("Nouvel opérateur dans la chambre", c.Id())
		c.Join("operators")

		AddOrUpdateOperator(c.Id(), Operator{
			Name:            "Opérateur anonyme",
			ControlledDrone: "",
			ChannelID:       c.Id(),
			IsAnonymous:     true,
			IsMobile:        false,
		})

		server.BroadcastTo("operators", "operator_update", "")
		for _, name := range ExtractDroneNames() {
			go server.BroadcastTo(name, "identify_operator", "")
			go server.BroadcastTo("operators", "drone_discovery", DroneIdentifier{
				Name: name,
			})
		}
	})

	server.On("identify_mobile", func(c *gosocketio.Channel) {
		log.Println("Nouvel opérateur mobile dans la chambre", c.Id())
		c.Join("mobile_operators")

		AddOrUpdateOperator(c.Id(), Operator{
			Name:            "Opérateur anonyme",
			ControlledDrone: "",
			ChannelID:       c.Id(),
			IsAnonymous:     true,
			IsMobile:        false,
		})
		server.BroadcastTo("operators", "operator_update", "")
		for _, name := range ExtractDroneNames() {
			go server.BroadcastTo(name, "identify_operator", "")
			go server.BroadcastTo("mobile_operators", "drone_discovery", DroneIdentifier{
				Name: name,
			})
		}
	})
	//endregion Refactoring

	server.On("authenticate", func(c *gosocketio.Channel, data OperatorIdentifier) {
		// Get Operator + merge
		AddOrUpdateOperator(c.Id(), Operator{
			Name:            data.Name,
			ControlledDrone: "",
			ChannelID:       c.Id(),
			IsAnonymous:     false,
		})
		server.BroadcastTo("operators", "operator_update", "")
	})

	server.On("release_controls", func(c *gosocketio.Channel) {
		RemoveLead(c.Id())
		server.BroadcastTo("operators", "operator_update", "")
	})

	// droneEventServer["shared"]
	server.On("position_update", func(c *gosocketio.Channel, data interface{}) {
		server.BroadcastTo("operators", "position_update", data)
	})

	// droneEventServer["shared"]
	server.On("acknowledge", func(c *gosocketio.Channel, data DroneIdentifier) {
		server.BroadcastTo("operators", "acknowledge", data)
	})

	// droneEventServer["shared"]
	server.On("flight_status_changed", func(c *gosocketio.Channel, data DroneFlyingStatusMessage) {
		UpdateFlyingStatus(data)
	})

	//region Outdoortest
	// droneEventServer["shared"]
	server.On("navigate_home_status_changed", func(c *gosocketio.Channel, data DroneNavigateHomeStatusMessage) {
		//UpdateFlyingStatus(data)
		server.BroadcastTo("operators", "navigate_home_status_changed", data)
	})

	// droneEventServer["shared"]
	server.On("on_alert_changed", func(c *gosocketio.Channel, data DroneAlertStatusMessage) {
		//UpdateFlyingStatus(data)
		server.BroadcastTo("operators", "on_alert_changed", data)
	})
	//endregion Outdoortest

	// droneEventServer["shared"]
	server.On("internal_status_changed", func(c *gosocketio.Channel, data interface{}) {
		server.BroadcastTo("operators", "internal_status_changed", data)
	})

	server.On("restart_module", func(c *gosocketio.Channel, data Module) {
		ModuleRestart(data)
	})

	// droneEventServer["shared"]
	server.On("on_command_success", func(c *gosocketio.Channel, data CommandIdentifier) {
		if drone, droneOk := AutomatonStatuses[data.Name]; droneOk && !drone.ManualFlight {
			NotifyScheduler(data)
		}
	})

	go StartWatcher()

	// Plan the auto-scaling system in order to split the charge on several sockets, and not one
	server.On("key_pressed", func(c *gosocketio.Channel, pressed_key OnTouchDown) {
		if drone, droneOk := AutomatonStatuses[pressed_key.DroneID]; droneOk && drone.ManualFlight && !drone.SimMode {
			if _, operatorOk := OperatorsInCharge[c.Id()]; operatorOk {
				setLeader := false
				if setLeader = SetLeadingOperator(c.Id(), pressed_key.DroneID); OperatorsInCharge[c.Id()].IsAnonymous || !setLeader {
					return
				}
				if setLeader {
					server.BroadcastTo("operators", "operator_update", "")
				}

			} else {
				return
			}

			command := DefineCommand(pressed_key)
			if command.Name != NoCommand {
				// droneControlServer["shared"] | droneControlServer[pressed_key.DroneID]
				server.BroadcastTo(pressed_key.DroneID, "command", command)
			}
		}
	})

	// droneControlServer["shared"] + droneEventServer["shared"]
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Channel %s disconnected", c.Id())

		if room, ok := channelMapping[c.Id()]; ok && room != "operator" {
			server.BroadcastTo("operators", "relay_endpoint_disconnect", DroneIdentifier{
				Name: channelMapping[c.Id()],
			})
		} else {
			RemoveOperator(c.Id())
			server.BroadcastTo("operators", "operator_update", "")
		}
	})
}

// StartWatcher Permet d'écouter les mises à jours de fichiers
func StartWatcher() {
	for {

		select {
		// watch for events
		case event := <-watcher.Events:
			if strings.Contains(event.Name, ".json") {
				droneName := strings.TrimRight(event.Name, ".json")
				if index := strings.LastIndex(droneName, "/"); index != -1 {
					droneName = droneName[index+1:]
				}

				newStatus := ExtractDroneStatus(droneName)
				AddOrUpdateDroneStatus(droneName, newStatus)

				autoPilot := GetAutopilotStatus(droneName)
				if autoPilot == (SchedulerSummarizedData{}) {
					autoPilot.DroneName = droneName
				}
				autoPilot.IsManual = newStatus.ManualFlight
				autoPilot.IsSimulated = newStatus.SimMode

				UpdateAutopilot(autoPilot)
				// Le répertoire Data concerne les infos remontées au Brain et au MapHandler.
				// Imaginer une autre logique pour le MapHandler
				server.BroadcastTo("operators", "automaton_status_changed", DroneIdentifier{
					Name: droneName,
				})
			}

			// Prepare the OSM section

		// watch for errors
		case err := <-watcher.Errors:
			failOnError(err, "Une erreur a été relevée par le Watcher")
		}
	}

}

// RedirectCommand Redirige la commande manuelle
func RedirectCommand(command RemoteManualCommand) {
	server.BroadcastTo(command.Target, string(command.Command), "")
}

// NotifyExternalCompChange Indique un changement dans un des modules externe
func NotifyExternalCompChange(droneName string) {
	if server != nil {
		server.BroadcastTo("operators", "external_module_update", DroneIdentifier{
			Name: droneName,
		})
	}
}

// SendLastCoordinate Envoyer les dernières informations
func SendLastCoordinate(drone DroneFlightCoordinates) {
	// coordinates.latitude,coordinates.longitude,self._drone_coordinates.altitude
	if server != nil {
		go server.BroadcastTo(drone.DroneName, "command", CreateAutomaticGoTo(&drone))
		go server.BroadcastTo("operators", "add_on_schedule", drone)
	}
}

// SendTargetCoordinates Envoi aux opérateurs les dernières coordonnées cibles
func SendTargetCoordinates(input FlightCoordinate) {
	if server != nil {
		go server.BroadcastTo("operators", "target_recalculated", input)
	}
}

// SendNodeLocation Envoi aux opérateurs les dernières coordonnées cibles
func SendNodeLocation(input FlightCoordinate) {
	if server != nil {
		go server.BroadcastTo("operators", "on_location_update", input)
	}
}

// SendAutopilotUpdate Sur mise à jour de l'ordonanceur, on informe les opérateurs
func SendAutopilotUpdate(input SchedulerSummarizedData) {
	if server != nil {
		go server.BroadcastTo("operators", "autopilot_update", input)
	}
}

// SendFlyingStatusUpdate Sur mise à jour de l'ordonanceur, on informe les opérateurs
func SendFlyingStatusUpdate(input DroneSummarizedStatus) {
	if server != nil {
		go server.BroadcastTo("operators", "flying_status_update", input)
	}
}

// SendAutomaticCommand Envoi d'une commande automatique
func SendAutomaticCommand(input PyAutomaticCommand) {
	if server != nil {
		go server.BroadcastTo(input.Target, "command", CreateAutomaticCommand(input))
	}
}
