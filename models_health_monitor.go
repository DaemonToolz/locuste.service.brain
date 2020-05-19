package main

/*
	Regroupement des éléments purements dédiés à la partie GO
*/
import "sync"

var moduleMutex sync.Mutex

// Component Module applicatif interactif
type Component string

const (
	// BrainWatcher Composant de surveillance de mise à jour des fichiers locaux JSON (OSM to be included)
	BrainWatcher Component = "Brain.FileWatcher"
	// BrainSocketServer Module SocketIO utilisé
	BrainSocketServer Component = "Brain.SocketServer"
	//BrainSocketHandler Deprecated
	BrainSocketHandler Component = "Brain.SocketHandler"
	// BrainHttpServer Module Serveur HTTP MUX utilisé
	BrainHttpServer Component = "Brain.HttpServer"
	// BrainMainRunner Thread principal
	BrainMainRunner Component = "Brain.Runner"
	// BrainRPCServer Module Serveur RPC utilisé
	BrainRPCServer Component = "Brain.RPCServer"
	// BrainSchedulerRPC Connexion à l'ordonanceur / locuste.service.osm
	BrainSchedulerRPC Component = "Brain.SchedulerConnection"

	// SchedulerRPCServer Serveur RPC de l'ordonanceur / locuste.service.osm
	SchedulerRPCServer Component = "Scheduler.RPCServer"
	// SchedulerRPC Connexion RPC vers l'unité de contrôle
	SchedulerRPC Component = "Scheduler.BrainConnection"
	// SchedulerMapHandler Gestionnaire de carte, module chargé de créer le graphe à partir d'une carte OSM
	SchedulerMapHandler Component = "Scheduler.MapHandler"
	// SchedulerFlightManager Module de pilotage automatique
	SchedulerFlightManager Component = "Scheduler.FlightManager"
)

// ExternalComponent Composant externe lié au drone - exemple : serveur vidéo
type ExternalComponent string

const (
	// VideoServer Serveur vidéo (JSMPEG par WebSocket)
	VideoServer ExternalComponent = "External.VideoServer"
	// VideoStream Flux vidéo (FFMPEG vers Websocket)
	VideoStream ExternalComponent = "External.VideoStream"
)

// Module System - Sous-sytème
type Module struct {
	System    string `json:"system"`
	SubSystem string `json:"subsystem"`
}

// ModuleRestartMapper Mappeur global
var ModuleRestartMapper map[Component]interface{}

func initModuleRestartMapper() {
	ModuleRestartMapper = make(map[Component]interface{})
	AddOrUpdateModuleMapper(BrainSocketServer, RestartSocketServer)
	AddOrUpdateModuleMapper(BrainHttpServer, RestartHTTPServer)
	AddOrUpdateModuleMapper(BrainRPCServer, RestartRPCServer)
	AddOrUpdateModuleMapper(SchedulerRPCServer, RestartFlightModule)
	AddOrUpdateModuleMapper(SchedulerFlightManager, RestartFlightSchedulerModule)

}

// AddOrUpdateModuleMapper Ajout d'une fonction de mise à jour
func AddOrUpdateModuleMapper(comp Component, function interface{}) {
	moduleMutex.Lock()
	ModuleRestartMapper[comp] = function
	moduleMutex.Unlock()
}

// CallModuleRestart Fonction pour redémarrer un module
func CallModuleRestart(comp Component) {
	moduleMutex.Lock()
	if _, ok := ModuleRestartMapper[comp]; ok {
		if result, ok := GlobalStatuses[comp]; !ok || (ok && !result) {
			ModuleRestartMapper[comp].(func())()
		}
	}
	moduleMutex.Unlock()
}
