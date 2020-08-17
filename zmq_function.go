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
	// ZFNRequestStatusesReply Réponse des services OSM à ZFNRequestStatuses
	ZFNRequestStatusesReply ZMQDefinedFunc = "RequestStatusesReply"
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

// ZOnUpdateAutopilot Après Mise à jour du pilote auto
func ZOnUpdateAutopilot(params *[]interface{}) {
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

// #endregion Function Host

// #region Function Client
const (
	// ZFNRequestStatuses Anciennement RequestStatuses (RPC)
	ZFNRequestStatuses ZMQDefinedFunc = "RequestStatuses"
	// ZFNNotifyScheduler Anciennement OnCommandSuccess (RPC)
	ZFNNotifyScheduler ZMQDefinedFunc = "OnCommandSuccess"
	// ZFNUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNUpdateAutopilot ZMQDefinedFunc = "UpdateAutopilot"
	// ZFNOnHomeChanged Anciennement OnHomeChanged (RPC)
	ZFNOnHomeChanged ZMQDefinedFunc = "OnHomeChanged"
	// ZFNFetchBoundaries Anciennement GetBoundaries (RPC)
	ZFNFetchBoundaries ZMQDefinedFunc = "GetBoundaries"
	// ZFNUpdateTarget Anciennement UpdateTarget (RPC)
	ZFNUpdateTarget ZMQDefinedFunc = "UpdateTarget"
	// ZFNUpdateFlyingStatus Anciennement FlyingStatusUpdate (RPC)
	ZFNUpdateFlyingStatus ZMQDefinedFunc = "FlyingStatusUpdate"
	// ZFNSendGoHomeCommandTo Anciennement SendGoHomeCommandTo (RPC)
	ZFNSendGoHomeCommandTo ZMQDefinedFunc = "SendGoHomeCommandTo"
	// ZFNSendTakeoffCommandTo Anciennement SendTakeoffCommandTo (RPC)
	ZFNSendTakeoffCommandTo ZMQDefinedFunc = "SendTakeoffCommandTo"

	// Reply sections - From client

	// ZFNRequestStatusReply Fonction réponse de RequestStatuses
	ZFNRequestStatusReply ZMQDefinedFunc = "RequestStatusesReply"
)

// ZRequestStatuses Demande d'envoi des derniers status
func ZRequestStatuses() *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNRequestStatuses,
		Params:   make([]interface{}, 0),
	}
}

// ZNotifyScheduler Notifier le séquenceur OSM
func ZNotifyScheduler(input CommandIdentifier) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNNotifyScheduler,
		Params:   []interface{}{input},
	}
}

// ZUpdateAutopilot Demande de MàJ forcée de l'autopilot
func ZUpdateAutopilot(input SchedulerSummarizedData) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNNotifyScheduler,
		Params:   []interface{}{input},
	}
}

// ZOnHomeChanged MàJ du point de décollage
func ZOnHomeChanged(input FlightCoordinate) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNOnHomeChanged,
		Params:   []interface{}{input},
	}
}

// ZFetchBoundaries Demande les limites de la carte
func ZFetchBoundaries() *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNFetchBoundaries,
		Params:   []interface{}{},
	}
}

// ZUpdateTarget Envoi des instructions pour recalculer la position sur la route
func ZUpdateTarget(input FlightCoordinate) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNUpdateTarget,
		Params:   []interface{}{input},
	}
}

// ZUpdateFlyingStatus mise à jour de l'état du drone (en vol)
func ZUpdateFlyingStatus(input DroneFlyingStatusMessage) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNUpdateFlyingStatus,
		Params:   []interface{}{input},
	}
}

// ZSendGoHomeCommandTo Demander une commande "atterrissage" au drone nommé
func ZSendGoHomeCommandTo(input string) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNSendGoHomeCommandTo,
		Params:   []interface{}{input},
	}
}

// ZSendTakeoffCommandTo Demander une commande "décollage" au drone nommé
func ZSendTakeoffCommandTo(input string) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNSendTakeoffCommandTo,
		Params:   []interface{}{input},
	}
}

// #endregion Function Client
