package main

/*
	Regroupement des éléments purements dédiés à la partie GO
	Inclu les informations envoyées par les systèmes UI
*/

// Key Touche clavier
type Key int

const (
	// Z Devant
	Z Key = 90
	// Q Gauche
	Q Key = 81
	// D Droite
	D Key = 68
	// S Arrière
	S Key = 83

	// A Rotation gauche
	A Key = 65
	// E Rotation droite
	E Key = 69

	// Ctrl vers le bas
	Ctrl Key = 17
	// Space vers le haut
	Space Key = 32

	// G Retour maison
	G Key = 71
	// T Décollage
	T Key = 84

	// ArrowUp caméra vers le haut
	ArrowUp = 38
	// ArrowDown caméra vers le base
	ArrowDown = 40
	// R Reset Camera
	R = 82
)

// OnTouchDown Identification d'une touche clavier pressée
type OnTouchDown struct {
	DroneID string `json:"drone_id"`
	KeyDown Key    `json:"key_pressed"`
}

// DroneIdentifier Informations envoyées à la GUI pour reconnaître le drone ciblé
type DroneIdentifier struct {
	Name string `json:"name"`
}

// PyDroneCommand Commande disponible pour le drone
type PyDroneCommand string

// Axis Axe d'un graphe
type Axis string

const (
	// NoCommand Aucune commande
	NoCommand PyDroneCommand = ""
	// GoTo Déplacment automatique
	GoTo PyDroneCommand = "AutomaticGoTo"
	// Stop Annulation du déplacment automatique
	Stop PyDroneCommand = "AutomaticCancelGoTo"
	// CamDown Rotation à 180°C de la caméra
	CamDown PyDroneCommand = "AutomaticSetCameraDown"
	// CamStd Remise à 0 zéro de la caméra
	CamStd PyDroneCommand = "AutomaticSetStandardCamera"
	// TakeOff Décollage
	TakeOff PyDroneCommand = "CommonTakeOff"
	// GoHome ORdre de retour à la maison
	GoHome PyDroneCommand = "CommonGoHome"
	// Move Commande de déplacement manuel
	Move PyDroneCommand = "ManualMove"
	// TiltCamera Changement de l'orientation de la caméra
	TiltCamera PyDroneCommand = "ManualTiltCamera"
	// ResetCamera Remise à l'état 0
	ResetCamera PyDroneCommand = "CommonSetStandardCamera"
)

const (
	// XAxis Direction devant/derrière
	XAxis Axis = "x"
	// YAxis Direction gauche/droite
	YAxis Axis = "y"
	// ZAxis Direction Bas/haut
	ZAxis Axis = "z"
	// OAxis N'est pas un axe,
	OAxis Axis = "orientation"
	// NoAxis aucun axe
	NoAxis Axis = ""

	// Axes de la caméra

	// Pitch Pitch de la caméra
	Pitch Axis = "pitch"
)

// UserControlCommand Commande disponible pour l'automate
type UserControlCommand string

const (
	// AskForManualFlight Demande d'accès au mode manuel
	AskForManualFlight UserControlCommand = "request_manual"
	// AskForAutomaticFlight Demande d'accès au mode automatique
	AskForAutomaticFlight UserControlCommand = "request_automatic"
	// AskForEmergencyDisconnect Déconnexion d'urgence, reprise en main forcée
	AskForEmergencyDisconnect UserControlCommand = "request_emergency_disconnect"
	// AskForEmergencyReconnect Reconnexion à l'automate
	AskForEmergencyReconnect UserControlCommand = "request_emergency_reconnect"
	// AskTestingMode Mode test / simulation
	AskTestingMode UserControlCommand = "request_simulation"
	// AskNormalMode Mode classique
	AskNormalMode UserControlCommand = "request_normal"
)

// RemoteManualCommand Commande manuelle de bascule
type RemoteManualCommand struct {
	Target  string             `json:"target"`
	Command UserControlCommand `json:"command"`
}
