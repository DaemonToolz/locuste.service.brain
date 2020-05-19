package main

import (
	"sync"
)

var (
	statusMutex    sync.Mutex
	automatonMutex sync.Mutex
	droneMutex     sync.Mutex
	extCompMutex   sync.Mutex
	autopilotMutex sync.Mutex
)

// GlobalStatuses Récupère l'état de fonctionnement des composants
var GlobalStatuses map[Component]bool

// AutomatonStatuses Représente l'état des automates
var AutomatonStatuses map[string]PyDroneStatus

// DroneStatuses Représente l'état des drones (batterie, etc...)
var DroneStatuses map[string]DroneStatus

// ExternalComponantStatuses Représente l'état composants externes (Format plus poussé comparé à GlobalStatuses, car regroupé par Drone)
var ExternalComponantStatuses map[string](map[ExternalComponent]bool)

// AutopilotStatuses Représente l'état des autopilotes
var AutopilotStatuses map[string]SchedulerSummarizedData

func initHealthMonitor() {
	GlobalStatuses = make(map[Component]bool)
	AutomatonStatuses = make(map[string]PyDroneStatus)
	DroneStatuses = make(map[string]DroneStatus)
	ExternalComponantStatuses = make(map[string](map[ExternalComponent]bool))
	AutopilotStatuses = make(map[string]SchedulerSummarizedData)

	for _, name := range ExtractDroneNames() {
		AddOrUpdateDroneStatus(name, ExtractDroneStatus(name))

		ExternalComponantStatuses[name] = make(map[ExternalComponent]bool)
		AddOrUpdateExtCompStatus(name, VideoServer, false)
		AddOrUpdateExtCompStatus(name, VideoStream, false)

	}
}

// AddOrUpdateStatus Met à jour l'information d'un composant
func AddOrUpdateStatus(component Component, isOnline bool) {
	statusMutex.Lock()
	GlobalStatuses[component] = isOnline
	statusMutex.Unlock()
}

// GetStatus Récupère l'état d'un composant global
func GetStatus(component Component) bool {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	return GlobalStatuses[component]
}

// AddOrUpdateDroneInternalStatus Met à jour l'information d'un statut interne au drone
func AddOrUpdateDroneInternalStatus(name string, stat DroneStatus) {
	droneMutex.Lock()
	DroneStatuses[name] = stat
	droneMutex.Unlock()
}

// GetDroneInternalStatus Récupère un statut interne
func GetDroneInternalStatus(name string) DroneStatus {
	droneMutex.Lock()
	defer droneMutex.Unlock()
	return DroneStatuses[name]
}

// AddOrUpdateDroneStatus Met à jour l'information d'un composant Drone
func AddOrUpdateDroneStatus(name string, status PyDroneStatus) {
	automatonMutex.Lock()
	AutomatonStatuses[name] = status
	automatonMutex.Unlock()
}

// GetDroneStatus Récupère l'état d'un composant dédié
func GetDroneStatuses(name string) PyDroneStatus {
	automatonMutex.Lock()
	defer automatonMutex.Unlock()
	return AutomatonStatuses[name]
}

// AddOrUpdateAutopilotStatus Met à jour l'état du pilote automatique
func AddOrUpdateAutopilotStatus(status SchedulerSummarizedData) {
	autopilotMutex.Lock()
	AutopilotStatuses[status.DroneName] = status
	autopilotMutex.Unlock()
}

// GetAutopilotStatus Récupère le dernier état d'un pilote automatique
func GetAutopilotStatus(name string) SchedulerSummarizedData {
	autopilotMutex.Lock()
	defer autopilotMutex.Unlock()
	return AutopilotStatuses[name]
}

// AddOrUpdateExtCompStatus Met à jour l'information d'un composant externe
func AddOrUpdateExtCompStatus(droneName string, component ExternalComponent, isOnline bool) {
	extCompMutex.Lock()
	ExternalComponantStatuses[droneName][component] = isOnline
	extCompMutex.Unlock()
	NotifyExternalCompChange(droneName)
}

// GetExtCompStatus Récupère l'état d'un composant externe
func GetExtCompStatus(droneName string) map[ExternalComponent]bool {
	extCompMutex.Lock()
	defer extCompMutex.Unlock()
	return ExternalComponantStatuses[droneName]

}
