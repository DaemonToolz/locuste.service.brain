package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/zeromq/goczmq"
)

// #region Section DRONES
var zmqDroneDealers map[string]*goczmq.Sock
var zmqDroneAccessMutex sync.Mutex
var zmqDroneCmdMutex sync.Mutex
var droneCmdChannels map[string]chan interface{}

// #endregion Section DRONES

// #region Section MODULES
var zmqModuleDealers map[string]*goczmq.Sock
var zmqModuleAccessMutex sync.Mutex
var zmqModuleCmdMutex sync.Mutex
var moduleCmdChannels map[string]chan interface{}

// #endregion Section MODULES

func initZMQ() {
	zmqDroneDealers = make(map[string]*goczmq.Sock)
	zmqModuleDealers = make(map[string]*goczmq.Sock)

	droneCmdChannels = make(map[string]chan interface{})
	moduleCmdChannels = make(map[string]chan interface{})
}

func createTgtZMQDealer(who *map[string]*goczmq.Sock, how *map[string]chan interface{}, modMutex *sync.Mutex, cmdMutex *sync.Mutex, request *ZMQIdentificationRequest, isExt bool) {
	var err error
	modMutex.Lock()
	defer modMutex.Unlock()

	if _, ok := (*how)[request.Name]; ok {
		cmdMutex.Lock()
		close((*how)[request.Name])
		delete(*how, request.Name)
		cmdMutex.Unlock()
	}

	(*who)[request.Name], err = goczmq.NewDealer(fmt.Sprintf("tcp://127.0.0.1:%d", request.ZMQPort)) // 5555 is Default ZMQ port

	if err != nil {
		failOnError(err, "CreateZMQDealer")
		delete((*who), request.Name)
		if isExt {
			AddOrUpdateExtCompStatus(request.Name, OrderStream, false)
		} else {
			AddOrUpdateStatus(Component(request.Name), false)
		}

	} else {
		cmdMutex.Lock()
		(*how)[request.Name] = make(chan interface{})
		cmdMutex.Unlock()
		if isExt {
			AddOrUpdateExtCompStatus(request.Name, OrderStream, true)
		} else {
			AddOrUpdateStatus(Component(request.Name), false)
		}
		go func(toWhom *map[string]*goczmq.Sock, commChan *map[string]chan interface{}, lock *sync.Mutex, name string, isInternal bool) {
			messageListenerLoop(toWhom, commChan, lock, name, isInternal)
		}(who, how, modMutex, request.Name, !isExt)
	}

	log.Println("Dealer ZeroMQ initialisé")
}

// CreateZMQDealer Création d'un Dealer CZMQ (rattaché au zmqSocket)
func CreateZMQDealer(request ZMQIdentificationRequest, internal bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if internal {
		createTgtZMQDealer(&zmqModuleDealers, &moduleCmdChannels, &zmqModuleAccessMutex, &zmqModuleCmdMutex, &request, !internal)
	} else {
		createTgtZMQDealer(&zmqDroneDealers, &droneCmdChannels, &zmqDroneAccessMutex, &zmqDroneCmdMutex, &request, !internal)
	}

}

func messageListenerLoop(toWhom *map[string]*goczmq.Sock, commChan *map[string]chan interface{}, lock *sync.Mutex, name string, isInternal bool) {
	for data := range (*commChan)[name] {
		SendZMQMessage(toWhom, lock, name, data, isInternal)
	}
}

// SendToZMQMessageChannelAuto Envoi d'une commande au canal de synchronization dédié à ZMQ (moins de paramètres)
func SendToZMQMessageChannelAuto(name string, payload interface{}, internal bool) {

	if internal {
		SendToZMQMessageChannel(&moduleCmdChannels, &zmqModuleCmdMutex, name, payload)
	} else {
		SendToZMQMessageChannel(&droneCmdChannels, &zmqDroneCmdMutex, name, payload)
	}

}

// SendToZMQMessageChannel Envoi d'une commande au canal de synchronization dédié à ZMQ
func SendToZMQMessageChannel(commChan *map[string]chan interface{}, cmdLock *sync.Mutex, name string, payload interface{}) {
	cmdLock.Lock()
	defer cmdLock.Unlock()
	if _, ok := (*commChan)[name]; ok {
		(*commChan)[name] <- payload
	}
}

// SendZMQMessage Envoyer un message via ZMQ (not thread-safe)
func SendZMQMessage(toWhom *map[string]*goczmq.Sock, lock *sync.Mutex, name string, payload interface{}, isInternal bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			if !isInternal {
				AddOrUpdateExtCompStatus(name, OrderStream, false)
			} else {
				AddOrUpdateStatus(Component(name), false)
			}
			DestroyZMQDealer(toWhom, lock, name, isInternal) // On détruit tout, car un bug s'est présenté
		}
	}()

	if _, ok := (*toWhom)[name]; ok {
		jPayload, err := json.Marshal(&payload)
		if err != nil {
			failOnError(err, fmt.Sprintf("SendMessSendZMQMessageage:%s", name))
			return
		}

		lock.Lock()
		_, erro := (*toWhom)[name].Write([]byte(jPayload))
		lock.Unlock()
		if erro != nil {
			failOnError(erro, fmt.Sprintf("SendMessSendZMQMessageage:%s", name))
		}

	}
}

// DestroyZMQRouters Destruction de toutes les connectiques CZMQ
func DestroyZMQRouters() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	for key := range zmqDroneDealers {
		DestroyZMQDealer(&zmqDroneDealers, &zmqDroneAccessMutex, key, false)
	}

	for key := range zmqModuleDealers {
		DestroyZMQDealer(&zmqModuleDealers, &zmqModuleAccessMutex, key, true)
	}
}

// DestroyZMQDealerAuto Destruction d'un dealer CZQM (moins de paramètres)
func DestroyZMQDealerAuto(key string, internal bool) {
	if internal {
		DestroyZMQDealer(&zmqModuleDealers, &zmqModuleAccessMutex, key, internal)
	} else {
		DestroyZMQDealer(&zmqDroneDealers, &zmqDroneAccessMutex, key, internal)
	}
}

// DestroyZMQDealer Destruction d'un dealer CZMQ
func DestroyZMQDealer(who *map[string]*goczmq.Sock, modMutex *sync.Mutex, name string, isInternal bool) {
	if !isInternal {
		AddOrUpdateExtCompStatus(name, OrderStream, false)
	} else {
		AddOrUpdateStatus(Component(name), false)
	}
	modMutex.Lock()
	defer modMutex.Unlock()
	if _, ok := (*who)[name]; ok {
		if (*who)[name] != nil {
			(*who)[name].Destroy()
		}
		delete((*who), name)
	}

}

// StartZMQServer Démarre un serveur ZMQ
func StartZMQServer(name string, port int, internal bool) {
	CreateZMQDealer(ZMQIdentificationRequest{
		Name:     name,
		ZMQPort:  port,
		Scope:    ZMQDrone,
		Internal: internal,
	}, internal)
}
