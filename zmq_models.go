package main

// ZMQScopeTarget Portée du message
type ZMQScopeTarget int

const (
	// ZMQInternal Communications internes
	ZMQInternal ZMQScopeTarget = iota
	// ZMQDrone Communications externes
	ZMQDrone ZMQScopeTarget = iota
)

// ZMQIdentificationRequest Identification pour la classification CZMQ
type ZMQIdentificationRequest struct {
	Name     string         `json:"name"`
	ZMQPort  int            `json:"zmq_port"`
	Scope    ZMQScopeTarget `json:"scope"`
	Internal bool           `json:"internal"`
}
