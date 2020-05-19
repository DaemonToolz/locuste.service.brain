package main

import (
	"log"
	"net/rpc"
	"reflect"
	"time"
)

// Note : Create a new shared project to regroup all duplicated code, structs and enums.
// Plan the refactoring in the next release / version

var client *rpc.Client
var myself Args
var pulse *time.Ticker
var stopCondition chan bool
var lastStatuses map[Component]bool

// NullArgType Type NIL à envoyer en paramètre
type NullArgType struct{}

// RPCNullArg
var RPCNullArg NullArgType

func initRPCClient() {
	pulse = time.NewTicker(1 * time.Second)
	lastStatuses = make(map[Component]bool)
	stopCondition = make(chan bool)
	RPCNullArg = NullArgType{}
	openConnection()
	go func() {
		for {
			select {
			case <-stopCondition:
				log.Println("Connectiques RPC arrêtées")
				AddOrUpdateStatus(SchedulerRPCServer, false)
				return
			case <-pulse.C:
				if client != nil {
					accessCall := client.Go("RPCRegistry.RequestStatuses", &RPCNullArg, &lastStatuses, nil)
					replyCall := <-accessCall.Done
					if client == nil {
						log.Println("La connexion n'était pas initialisée")
						openConnection()
					} else if replyCall.Error == rpc.ErrShutdown || reflect.TypeOf(replyCall.Error) == reflect.TypeOf((*rpc.ServerError)(nil)).Elem() {
						log.Println("Une erreur liée au serveur a été remonté")
						log.Println(replyCall.Error)
						openConnection()
					} else {
						FetchBoundaries() // On force la MàJ
					}
				} else {
					openConnection()
				}

				if lastStatuses != nil {
					for key := range lastStatuses {
						AddOrUpdateStatus(key, lastStatuses[key])
					}
				}
			}
		}
	}()

	log.Println("Connectiques RPC initialisés")
}

func openConnection() *rpc.Client {
	initConfiguration()
	var err error
	client, err = rpc.DialHTTP("tcp", appConfig.rpcSchedulerPort())
	if err != nil {
		AddOrUpdateStatus(BrainSchedulerRPC, false)
		failOnError(err, "couldn't connect to remote RPC server")
	} else {
		AddOrUpdateStatus(BrainSchedulerRPC, true)
	}
	return client
}

// RequestStatuses Demande le statut des modules côté locuste.service.osm
func RequestStatuses() {
	if client != nil {
		client.Go("RPCRegistry.RequestStatuses", &RPCNullArg, &lastStatuses, nil)

	}
}

// Unregister Désenregistre un module connecté via RPC
func Unregister() {
	if client != nil {
		defer client.Close()

	}
	AddOrUpdateStatus(BrainSchedulerRPC, false)
}

// NotifyScheduler Notification de l'ordonanceur
func NotifyScheduler(data DroneIdentifier) {
	if client != nil {
		client.Go("RPCRegistry.OnCommandSuccess", &data, nil, nil)
	}
}

// UpdateAutopilot Mise à jour d'un ordonanceur de vol
func UpdateAutopilot(input SchedulerSummarizedData) {
	if client != nil && input.DroneName != "" {
		client.Go("RPCRegistry.UpdateAutopilot", &input, &RPCNullArg, nil)
	}
}

// OnHomeChanged Dès le décollage
func OnHomeChanged(output FlightCoordinate) {
	if client != nil {
		client.Go("RPCRegistry.OnHomeChanged", &output, &RPCNullArg, nil)
	}
}

// FetchBoundaries Récupère les limites de la carte
func FetchBoundaries() {
	if client != nil { // && flightSchedulerRPC.MapBoundaries == (Boundaries{}) {
		client.Call("RPCRegistry.GetBoundaries", &RPCNullArg, &flightSchedulerRPC.MapBoundaries)
	}
}

// UpdateTarget Envoi des instructions pour recalculer la position sur la route
func UpdateTarget(input FlightCoordinate) {
	if client != nil && input != (FlightCoordinate{}) {
		client.Go("RPCRegistry.UpdateTarget", &input, &RPCNullArg, nil)
	}
}
