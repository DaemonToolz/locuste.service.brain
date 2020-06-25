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
	ArrowUp Key = 38
	// ArrowDown caméra vers le base
	ArrowDown Key = 40
	// R Reset Camera
	R Key = 82
)

// JoystickType Type de joystick (contrôle vertical / rotation ou horizontal)
type JoystickType int

const (
	// SpeedJoystick Joystick dédié au mouvement (Forwards, Back, Left, Right)
	SpeedJoystick JoystickType = iota
	// AltitutdeJoystick Joystick dédié à l'altitude (Up, Down, Rotate left, Rotate right)
	AltitutdeJoystick JoystickType = iota
)

// OnTouchDown Identification d'une touche clavier pressée
type OnTouchDown struct {
	DroneID string `json:"drone_id"`
	KeyDown Key    `json:"key_pressed"`
}

// OnJoystickEvent Evénements joystick
type OnJoystickEvent struct {
	DroneID string       `json:"drone_id"`
	Payload interface{}  `json:"payload"`
	Type    JoystickType `json:"joystick_type"`
}

// SpeedJoystickEvent Evénements joystick SpeedJoystick
type SpeedJoystickEvent struct {
	Yaw  float64 `json:"yaw"`
	Roll float64 `json:"roll"`
}

// AltJoystickEvent Evénements joystick AltitutdeJoystick
type AltJoystickEvent struct {
	Up       float64 `json:"up"`
	Rotation float64 `json:"rotation"`
}

// DroneIdentifier Informations envoyées à la GUI pour reconnaître le drone ciblé
type DroneIdentifier struct {
	Name string `json:"name"`
}

// CommandIdentifier Le "acknowledge" d'un drone pour une commande spécifique
type CommandIdentifier struct {
	Name    string         `json:"name"`
	Command PyDroneCommand `json:"command"`
}

// PyDroneCommand Commande disponible pour le drone
type PyDroneCommand string

// Axis Axe d'un graphe
type Axis string

const (
	// NoCommand Aucune commande
	NoCommand PyDroneCommand = "NoCommand"
	// GoTo Déplacment automatique
	GoTo PyDroneCommand = "AutomaticGoTo"
	// Stop Annulation du déplacment automatique
	Stop PyDroneCommand = "AutomaticCancelGoTo"
	// CamDown Rotation à 180°C de la caméra
	CamDown PyDroneCommand = "AutomaticSetCameraDown"
	// CamStd Remise à 0 zéro de la caméra
	CamStd PyDroneCommand = "AutomaticSetStandardCamera"
	// TiltCamera Changement de l'orientation de la caméra
	TiltCamera PyDroneCommand = "ManualTiltCamera"
	// Move Commande de déplacement manuel
	Move PyDroneCommand = "ManualMove"
	// Tilt Commande de déplacement manuel
	Tilt PyDroneCommand = "ManualTilt"

	// AutomaticTakeOff Ordre de décollage automatique
	AutomaticTakeOff PyDroneCommand = "AutomaticTakeOff"
	// AutomaticGoHome Ordre de retour à la maison automatique
	AutomaticGoHome PyDroneCommand = "AutomaticGoHome"
	// AutomaticLand Ordre d'atterrissage automatique
	AutomaticLand PyDroneCommand = "AutomaticLanding"

	// Section à remplacer

	// TakeOff Décollage
	TakeOff PyDroneCommand = "CommonTakeOff"
	// GoHome ORdre de retour à la maison
	GoHome PyDroneCommand = "CommonGoHome"
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
