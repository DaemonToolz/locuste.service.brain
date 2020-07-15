package main

import (
	"sync"
)

var (
	zmqCallMutex sync.Mutex
)

// ZMQCallMap Carte d'appel ZMQ - Golang func
var ZMQCallMap map[ZMQDefinedFunc]interface{}

func initZMQMapper() {
	ZMQCallMap = make(map[ZMQDefinedFunc]interface{})
}

// AddOrUpdateZMQCaller Met à jour une des fonctions d'appel
func AddOrUpdateZMQCaller(name ZMQDefinedFunc, fn interface{}) {
	zmqCallMutex.Lock()
	ZMQCallMap[name] = fn
	zmqCallMutex.Unlock()
}

// GetZMQCaller Récupère une des fonctions d'appel
func GetZMQCaller(name ZMQDefinedFunc) interface{} {
	zmqCallMutex.Lock()
	defer zmqCallMutex.Unlock()
	return ZMQCallMap[name]
}
