package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/zeromq/goczmq"
)

var zmqDealers map[string]*goczmq.Sock
var zmqAccessMutex sync.Mutex
var zmqCmdMutex sync.Mutex

var commandChannels map[string]chan interface{}

func initZMQ() {
	zmqDealers = make(map[string]*goczmq.Sock)
	commandChannels = make(map[string]chan interface{})
}

// CreateZMQDealer Création d'un Dealer CZMQ (rattaché au zmqSocket)
func CreateZMQDealer(request IdentificationRequest) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var err error
	zmqAccessMutex.Lock()
	defer zmqAccessMutex.Unlock()

	if _, ok := commandChannels[request.Name]; ok {
		zmqCmdMutex.Lock()

		close(commandChannels[request.Name])
		delete(commandChannels, request.Name)
		zmqCmdMutex.Unlock()
	}

	zmqDealers[request.Name], err = goczmq.NewDealer(fmt.Sprintf("tcp://127.0.0.1:%d", request.ZMQPort)) // 5555 is Default ZMQ port

	if err != nil {
		failOnError(err, "CreateZMQDealer")
		delete(zmqDealers, request.Name)
	} else {
		zmqCmdMutex.Lock()
		commandChannels[request.Name] = make(chan interface{})
		zmqCmdMutex.Unlock()
		go func(name string) { messageListenerLoop(name) }(request.Name)
	}

	log.Println("Dealer ZeroMQ initialisé")
}

func messageListenerLoop(name string) {
	for data := range commandChannels[name] {
		SendZMQMessage(name, data)
	}
}

// SendToZMQMessageChannel Envoi d'une commande au canal de synchronization dédié à ZMQ
func SendToZMQMessageChannel(name string, payload interface{}) {
	zmqCmdMutex.Lock()
	defer zmqCmdMutex.Unlock()
	if _, ok := commandChannels[name]; ok {
		commandChannels[name] <- payload
	}
}

// SendZMQMessage Envoyer un message via ZMQ (not thread-safe)
func SendZMQMessage(name string, payload interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if _, ok := zmqDealers[name]; ok {
		jPayload, err := json.Marshal(&payload)
		if err != nil {
			failOnError(err, fmt.Sprintf("SendMessSendZMQMessageage:%s", name))
			return
		}

		zmqAccessMutex.Lock()
		_, erro := zmqDealers[name].Write([]byte(jPayload))
		zmqAccessMutex.Unlock()
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

	for key := range zmqDealers {
		DestroyZMQDealer(key)
	}
}

// DestroyZMQDealer Destruction d'un dealer CZMQ
func DestroyZMQDealer(name string) {

	zmqAccessMutex.Lock()
	defer zmqAccessMutex.Unlock()
	if _, ok := zmqDealers[name]; ok {
		if zmqDealers[name] != nil {
			zmqDealers[name].Destroy()
		}
		delete(zmqDealers, name)
	}
}
