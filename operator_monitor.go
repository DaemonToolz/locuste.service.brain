package main

import "sync"

var operatorMutex sync.Mutex

// OperatorsInCharge Opérateurs identifiés dans l'application
var OperatorsInCharge map[string]Operator

func initOperatorDictionary() {
	OperatorsInCharge = make(map[string]Operator)
}

// AddOrUpdateOperator Met à jour l'information relatif aux opérateurs connectés
func AddOrUpdateOperator(channelID string, operator Operator) {
	operatorMutex.Lock()
	OperatorsInCharge[channelID] = operator
	operatorMutex.Unlock()
}

// RemoveOperator Retire un opérateur
func RemoveOperator(channelID string) {
	RemoveLead(channelID)
	operatorMutex.Lock()
	delete(OperatorsInCharge, channelID)
	operatorMutex.Unlock()
}

// RemoveLead Retire retire le LeadingOperator
func RemoveLead(channelID string) {
	operatorMutex.Lock()
	defer operatorMutex.Unlock()
	if op, ok := OperatorsInCharge[channelID]; ok {
		op.ControlledDrone = ""
		OperatorsInCharge[channelID] = op
	}
}

// SetLeadingOperator Donne le lead à un opérateur
func SetLeadingOperator(channelID string, drone string) bool {
	operatorMutex.Lock()
	defer operatorMutex.Unlock()
	if OperatorsInCharge[channelID].IsAnonymous || (OperatorsInCharge[channelID].ControlledDrone != "") { // Ce canal a le lead sur une instance
		if OperatorsInCharge[channelID].ControlledDrone != drone {
			return false
		}
		return true
	}

	for k := range OperatorsInCharge {
		if OperatorsInCharge[k].ControlledDrone == drone {
			return false
		}
	}
	op := OperatorsInCharge[channelID]
	op.ControlledDrone = drone
	OperatorsInCharge[channelID] = op
	return true
}

// IsLeadingDrone Est-ce que l'opérateur contrôle le drone
func IsLeadingDrone(channelID string, drone string) bool {
	operatorMutex.Lock()
	defer operatorMutex.Unlock()
	if OperatorsInCharge[channelID].IsAnonymous {
		return false
	}
	return OperatorsInCharge[channelID].ControlledDrone == drone
}
