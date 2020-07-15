package main

// ZMQDefinedFunc Noms des fonctions échangées Router <=> Dealer
type ZMQDefinedFunc string

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
	// ZFNDefineEdge Anciennement DefineEdge (RPC)
	ZFNDefineEdge ZMQDefinedFunc = "DefineEdge"
	// ZFNOnUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNOnUpdateAutopilot ZMQDefinedFunc = "OnUpdateAutopilot"
	// ZFNOnFlyingStatusUpdate Anciennement OnFlyingStatusUpdate (RPC)
	ZFNOnFlyingStatusUpdate ZMQDefinedFunc = "OnFlyingStatusUpdate"
	// ZFNServerShutdown Anciennement ServerShutdown (RPC)
	ZFNServerShutdown ZMQDefinedFunc = "ServerShutdown"
	// ZFNSendCommand Anciennement SendCommand (RPC)
	ZFNSendCommand ZMQDefinedFunc = "SendCommand"
)

// #endregion Function Host

// #region Function Client
const (
	// ZFNRequestStatuses Anciennement RequestStatuses (RPC)
	ZFNRequestStatuses ZMQDefinedFunc = "RequestStatuses"
	// ZFNNotifySchedulerAnciennement NotifyScheduler (RPC)
	ZFNNotifyScheduler ZMQDefinedFunc = "NotifyScheduler"
	// ZFNUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNUpdateAutopilot ZMQDefinedFunc = "UpdateAutopilot"
)

// #endregion Function Client

// ZMQMessage Message envoyé entre les Dealers ZMQ
type ZMQMessage struct {
	Function ZMQDefinedFunc `json:"function"`
	Params   []interface{}  `json:"params"`
}

/*


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

// OnFlyingStatusUpdate On a une mise à jour côté Scheduler
func (*RPCRegistry) OnFlyingStatusUpdate(args *DroneSummarizedStatus, _ *struct{}) error {
	AddOrUpdateFlyingStatus(*args)
	go SendFlyingStatusUpdate(*args)
	return nil
}

// ServerShutdown Arrêt du serveur RPC
func (*RPCRegistry) ServerShutdown(_ *struct{}, _ *struct{}) error {
	AddOrUpdateStatus(SchedulerRPCServer, false)
	return nil
}
func (t *RPCRegistry) RPCSendCommand(command *PyAutomaticCommand, _ *struct{}) error {
	go SendAutomaticCommand(*command)
	return nil
}
*/

/*


// RequestStatuses Demande le statut des modules côté locuste.service.osm
func RequestStatuses() {
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
