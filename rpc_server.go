package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Args Objet RPC envoyé par le module locuste.service.osm
type Args struct {
	PId       int
	Component Component
}

// RPCRegistry Informations en transit
type RPCRegistry struct {
	RPCComponents map[Component]int `json:"rcp_components"`
	MapBoundaries Boundaries        `json:"boundaries"`
}

// Register Enregistre un module qui se connecte à l'unité de contrôle
func (t *RPCRegistry) Register(args *Args, _ *struct{}) error {
	t.RPCComponents[args.Component] = args.PId
	log.Println("Processus RPC ajouté ", string(args.Component))
	AddOrUpdateStatus(args.Component, true)
	return nil
}

// Disconnect Indique qu'un module s'est déconnecté
func (t *RPCRegistry) Disconnect(args *Args, _ *struct{}) error {
	delete(t.RPCComponents, args.Component)
	log.Println("Processus RPC stoppé ", string(args.Component))
	AddOrUpdateStatus(args.Component, false)
	return nil
}

// DefineBoundaries Définir les limites de la carte
func (t *RPCRegistry) DefineBoundaries(args *Boundaries, _ *struct{}) error {
	t.MapBoundaries = *args
	return nil
}

// SendCoordinates Envoi des coordonnées au serveur SocketIO
func (*RPCRegistry) SendCoordinates(args *DroneFlightCoordinates, _ *struct{}) error {
	go SendLastCoordinate(*args)
	return nil
}

// DefineTarget Mise à jour de la cible (déplacement)
func (*RPCRegistry) DefineTarget(args *FlightCoordinate, _ *struct{}) error {
	go SendTargetCoordinates(*args)
	return nil
}

// DefineEdge Deprecated: Envoi des informations du graphe/ville définit dans le module locuste.service.osm
func (*RPCRegistry) DefineEdge(args *FlightCoordinate, _ *struct{}) error {
	go SendNodeLocation(*args)
	return nil
}

// UpdateAutopilot Demande la mise à jour du pilote / ordonanceur d'un drone
func (*RPCRegistry) UpdateAutopilot(args *SchedulerSummarizedData, _ *struct{}) error {
	AddOrUpdateAutopilotStatus(*args)
	go SendAutopilotUpdate(*args)
	return nil
}

// ServerShutdown Arrêt du serveur RPC
func (*RPCRegistry) ServerShutdown(_ *struct{}, _ *struct{}) error {
	AddOrUpdateStatus(SchedulerRPCServer, false)
	return nil
}

// Ping Fonction de ping
func (t *RPCRegistry) Ping(_ *struct{}, reply *string) error {
	*reply = ModuleToRestart
	ModuleToRestart = ""
	return nil
}

// RPCSendCommand Envoi d'une commande [Automatic....]
func (t *RPCRegistry) RPCSendCommand(command *PyAutomaticCommand, _ *struct{}) error {
	go SendAutomaticCommand(*command)
	return nil
}

var flightSchedulerRPC *RPCRegistry

// ModuleToRestart Module à redémarrer (récupéré en même temps que le "PING")
var ModuleToRestart string

// RestartFlightModule Redémarrage du module SchedulerRPCServer de locust.service.osm
func RestartFlightModule() {
	ModuleToRestart = string(SchedulerRPCServer)
}

// RestartFlightSchedulerModule Redémarrage du module Ordonanceur de locust.service.osm
func RestartFlightSchedulerModule() {
	ModuleToRestart = string(SchedulerFlightManager)
}

func initRemoteProcedureCall() {
	listener, err := net.Listen("tcp", appConfig.rpcListenUri())
	if err != nil {
		failOnError(err, "Couldn't initialize the RPC listener")
	}

	ModuleToRestart = ""
	if flightSchedulerRPC == nil {
		flightSchedulerRPC = &RPCRegistry{
			RPCComponents: make(map[Component]int),
			MapBoundaries: Boundaries{},
		}

		rpc.Register(flightSchedulerRPC)
		rpc.HandleHTTP()

	}

	log.Println("Ouverture des ports HTTP pour le processus RPC", listener.Addr().(*net.TCPAddr).Port)
	http.Serve(listener, nil)

}

// RestartRPCServer Redémarrage du serveur RPC local
func RestartRPCServer() {
	if result, ok := GlobalStatuses[BrainRPCServer]; !ok || (ok && !result) {
		initConfiguration()
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
					AddOrUpdateStatus(BrainRPCServer, false)
				}
			}()
			AddOrUpdateStatus(BrainRPCServer, true)
			initRemoteProcedureCall()
			AddOrUpdateStatus(BrainRPCServer, false)
		}()
	}
}
