package main

import (
	"fmt"
	"sync"
)

// ZMQDefinedFunc Noms des fonctions échangées Router <=> Dealer
type ZMQDefinedFunc string

// ZMQMessage Message envoyé entre les Dealers ZMQ
type ZMQMessage struct {
	Function ZMQDefinedFunc `json:"function"`
	Params   []interface{}  `json:"params"`
}

// ZMQComponents Composants enregistrés (Avec ProcessID)
var ZMQComponents map[Component]int
var zmqProcessMutex sync.Mutex

// MapBoundaries Limite de la carte
var MapBoundaries Boundaries

func init() {
	ZMQComponents = make(map[Component]int)
	MapBoundaries = Boundaries{}
}

// #region Function Host
const (
	// ZFNRegister Anciennement Register (RPC)
	ZFNRegister ZMQDefinedFunc = "Register"
	// ZFNDisconnect Anciennement Disconnect (RPC)
	ZFNDisconnect ZMQDefinedFunc = "Disconnect"
	// ZFNDisconnect Anciennement DefineBoundaries (RPC)
	ZFNDefineBoundaries ZMQDefinedFunc = "DefineBoundaries"
	// ZFNSendCoordinates Anciennement SendCoordinates (RPC)
	ZFNSendCoordinates ZMQDefinedFunc = "SendCoordinates"
	// ZFNDefineTarget Anciennement DefineTarget (RPC)
	ZFNDefineTarget ZMQDefinedFunc = "DefineTarget"
	// ZFNOnUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNOnUpdateAutopilot ZMQDefinedFunc = "OnUpdateAutopilot"
	// ZFNOnFlyingStatusUpdate Anciennement OnFlyingStatusUpdate (RPC)
	ZFNOnFlyingStatusUpdate ZMQDefinedFunc = "OnFlyingStatusUpdate"
	// ZFNServerShutdown Anciennement ServerShutdown (RPC)
	ZFNServerShutdown ZMQDefinedFunc = "ServerShutdown"
	// ZFNSendCommand Anciennement SendCommand (RPC)
	ZFNSendCommand ZMQDefinedFunc = "SendCommand"
)

func addOrUpdateZMQProcess(cpt Component, pid int) {
	zmqProcessMutex.Lock()
	defer zmqProcessMutex.Unlock()
	ZMQComponents[cpt] = pid
}

func deleteZMQProcess(cpt Component) {
	zmqProcessMutex.Lock()
	defer zmqProcessMutex.Unlock()
	delete(ZMQComponents, cpt)
}

// ZRegister Enregistrer le processus ZMQ
func ZRegister(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(Args); len(*params) > 0 && ok {
			addOrUpdateZMQProcess(input.Component, input.PId)
			trace(fmt.Sprintf("%s : %s", callSuccess, string(input.Component)))
			AddOrUpdateStatus(input.Component, true)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZRDisconnect Désenregistre le process associé à une file ZMQ
func ZRDisconnect(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(Args); ok {
			deleteZMQProcess(input.Component)
			trace(fmt.Sprintf("%s : %s", callSuccess, string(input.Component)))
			AddOrUpdateStatus(input.Component, false)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZDefineBoundaries Enregistre les limites de la carte
func ZDefineBoundaries(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(Boundaries); ok {
			MapBoundaries = input
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZSendCoordinates Partage les dernières coordonnées avec les clients Angular / Mobile
func ZSendCoordinates(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(DroneFlightCoordinates); ok {
			go SendLastCoordinate(input)
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZDefineTarget Partage les coordonnées cibles avec les clients Angular / Mobile
func ZDefineTarget(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(FlightCoordinate); len(*params) > 0 && ok {
			go SendTargetCoordinates(input)
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZUpdateAutopilot Mise à jour du pilote auto
func ZUpdateAutopilot(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(SchedulerSummarizedData); len(*params) > 0 && ok {
			AddOrUpdateAutopilotStatus(input)
			go SendAutopilotUpdate(input)
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZOnFlyingStatusUpdate Réception des nouvelles infos de vol
func ZOnFlyingStatusUpdate(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(DroneSummarizedStatus); len(*params) > 0 && ok {
			AddOrUpdateFlyingStatus(input)
			go SendFlyingStatusUpdate(input)
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// ZServerShutdown Indicateur d'arrêt du endpoint / service ZMQ distant
func ZServerShutdown(params *[]interface{}) {
	AddOrUpdateStatus(SchedulerRPCServer, false)
}

// ZSendCommand Transmission d'une commande remontée par le service OSM
func ZSendCommand(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(PyAutomaticCommand); len(*params) > 0 && ok {
			go SendAutomaticCommand(input)
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// #endregion Function Host

// #region Function Client
const (
	// ZFNRequestStatuses Anciennement RequestStatuses (RPC)
	ZFNRequestStatuses ZMQDefinedFunc = "RequestStatuses"
	// ZFNNotifySchedulerAnciennement NotifyScheduler (RPC)
	ZFNNotifyScheduler ZMQDefinedFunc = "NotifyScheduler"
	// ZFNUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNUpdateAutopilot ZMQDefinedFunc = "UpdateAutopilot"
	// ZFNOnHomeChanged Anciennement OnHomeChanged (RPC)
	ZFNOnHomeChanged ZMQDefinedFunc = "OnHomeChanged"
	// ZFNFetchBoundaries Anciennement FetchBoundaries (RPC)
	ZFNFetchBoundaries ZMQDefinedFunc = "FetchBoundaries"
	// ZFNUpdateTarget Anciennement UpdateTarget (RPC)
	ZFNUpdateTarget ZMQDefinedFunc = "UpdateTarget"
	// ZFNUpdateFlyingStatus Anciennement UpdateFlyingStatus (RPC)
	ZFNUpdateFlyingStatus ZMQDefinedFunc = "UpdateFlyingStatus"
	// ZFNSendGoHomeCommandTo Anciennement SendGoHomeCommandTo (RPC)
	ZFNSendGoHomeCommandTo ZMQDefinedFunc = "SendGoHomeCommandTo"
	// ZFNSendTakeoffCommandTo Anciennement SendTakeoffCommandTo (RPC)
	ZFNSendTakeoffCommandTo ZMQDefinedFunc = "SendTakeoffCommandTo"

	// Reply sections

	// ZFNRequestStatusReply Fonction réponse de RequestStatuses
	ZFNRequestStatusReply ZMQDefinedFunc = "RequestStatusesReply"
)

// ZRequestStatuses Demande d'envoi des derniers status
func ZRequestStatuses() {

}

// ZRequestStatusesReply  Réception des derniers status
func ZRequestStatusesReply(params *[]interface{}) {
	if params != nil && len(*params) > 0 {
		if input, ok := (*params)[0].(map[Component]bool); len(*params) > 0 && ok {
			lastStatuses = input
			trace(callSuccess)
		} else {
			trace(callFailure)
		}
	} else {
		trace(callFailure)
	}
}

// #endregion Function Client

/*


// RequestStatuses Demande le statut des modules côté locuste.service.osm
func () {
	if client != nil {
		client.Go("RPCRegistry.RequestStatuses", &RPCNullArg, &lastStatuses, nil)

	}
}

// NotifyScheduler Notification de l'ordonanceur
func NotifyScheduler(data CommandIdentifier) {
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

// UpdateFlyingStatus mise à jour de l'état du drone (en vol)
func UpdateFlyingStatus(data DroneFlyingStatusMessage) {
	if client != nil {
		client.Go("RPCRegistry.FlyingStatusUpdate", &data, &RPCNullArg, nil)
	}
}

// SendGoHomeCommandTo Demander une commande "atterrissage" au drone nommé
func SendGoHomeCommandTo(name string) {
	if client != nil {
		client.Go("RPCRegistry.SendGoHomeCommandTo", &name, &RPCNullArg, nil)
	}
}

// SendTakeoffCommandTo Demander une commande "décollage" au drone nommé
func SendTakeoffCommandTo(name string) {
	if client != nil {
		client.Go("RPCRegistry.SendTakeoffCommandTo", &name, &RPCNullArg, nil)
	}
}
*/
