package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

/*  DefineCommand permet de définir l'ordre à envoyer au code Python par le biais des websockets
params : OnTouchDown : touche pressée
return : PyDromeCommandMessage : ordre à envoyer
*/
func DefineCommand(keyPressed OnTouchDown) PyDroneCommandMessage {
	var defaultSpeed float64
	var command PyDroneCommand = NoCommand
	multiplier := 0.0
	var finalOrder PyDroneCommandMessage = PyDroneCommandMessage{}
	var axis Axis = NoAxis

	switch keyPressed.KeyDown {
	case ArrowUp:
		multiplier = 1
		fallthrough
	case ArrowDown:
		if multiplier == 0 {
			multiplier = -1
		}
		axis = Pitch
		command = TiltCamera

	case Z:
		axis = XAxis
		fallthrough
	case D:
		if axis == NoAxis {
			axis = YAxis
		}
		fallthrough
	case E:
		if axis == NoAxis {
			axis = OAxis
		}
		fallthrough
	case Ctrl:
		if axis == NoAxis {
			axis = ZAxis
		}
		command = Move
		multiplier = 1

	case S:
		axis = XAxis
		fallthrough

	case Q:
		if axis == NoAxis {
			axis = YAxis
		}
		fallthrough
	case A:
		if axis == NoAxis {
			axis = OAxis
		}
		fallthrough
	case Space:
		if axis == NoAxis {
			axis = ZAxis
		}
		command = Move
		multiplier = -1
	}

	tempParams := make(map[string]float64) // On caste en int car la SDK Olympe ARM supported mal les flottant sur

	switch command {
	case NoCommand:
		// La section ci-dessous est traitée à par car elle n'utilise pas de paramètres
		switch keyPressed.KeyDown {
		case G:
			command = GoHome
		case T:
			command = TakeOff
		case R:
			command = ResetCamera
		}
		finalOrder.Params = nil

	case Move:
		defaultSpeed = 0.20 * multiplier
		tempParams[string(XAxis)] = 0
		tempParams[string(YAxis)] = 0
		tempParams[string(ZAxis)] = 0
		tempParams[string(OAxis)] = 0
		tempParams[string(axis)] = defaultSpeed
		finalOrder.Params = tempParams

	case TiltCamera:
		defaultSpeed = 0.2 * multiplier
		tempParams[string(axis)] = defaultSpeed
		finalOrder.Params = tempParams
	}

	finalOrder.Name = command
	return finalOrder
}

/*  CreateAutomaticGoTo Créer un ordre reçu par l'autopilote

 */
func CreateAutomaticGoTo(input *DroneFlightCoordinates) PyDroneCommandMessage {

	var finalOrder PyDroneCommandMessage = PyDroneCommandMessage{}

	tempParams := make(map[string]float64) // On caste en int car la SDK Olympe ARM supported mal les flottant sur

	tempParams["latitude"] = input.Component.Lat
	tempParams["longitude"] = input.Component.Lon

	finalOrder.Name = GoTo
	finalOrder.Params = tempParams
	return finalOrder
}

// ExtractDroneNames Récupère les noms des drones
func ExtractDroneNames() []string {
	droneNames := make([]string, 0)

	for _, drone := range drones.Drones {
		droneNames = append(droneNames, fmt.Sprintf("ANAFI_%s", drone.IpAddress))
	}

	return droneNames
}

// ExtractDroneStatus Récupère le dernier status enregistré
func ExtractDroneStatus(name string) PyDroneStatus {
	var pyDroneStatus PyDroneStatus
	statusFile, err := os.Open(fmt.Sprintf("/home/pi/project/locuste/data/%s.json", name))
	defer statusFile.Close()
	if err != nil {
		return PyDroneStatus{Available: false}
	}
	jsonParser := json.NewDecoder(statusFile)
	jsonParser.Decode(&pyDroneStatus)
	pyDroneStatus.Available = true
	return pyDroneStatus
}

// ModuleRestart Redémarrage d'un module
func ModuleRestart(module Module) {
	log.Println("Module à redémarrer : ", module)
	CallModuleRestart(Component(module.System + "." + module.SubSystem))
}
