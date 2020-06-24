package main

/*
	Regroupement des éléments purements dédiés à la partie Pyton
*/

// IdentificationRequest Requête d'identification
type IdentificationRequest struct {
	Name         string      `json:"name"`
	VideoPort    int         `json:"video_port"`
	ZMQPort      int         `json:"zmq_port"`
	IP           string      `json:"ip"`
	Connected    bool        `json:"connected"`
	ManualFlight bool        `json:"manual"`
	SimMode      bool        `json:"sim"`
	Position     interface{} `json:"position"`
}

// NavigationType Mode de navigation, si Manual, alors manuel, sinon automatique
type NavigationType struct {
	Manual bool `json:"manual"`
}

// PyDroneCommandMessage Ordre à envoyer aux drones
type PyDroneCommandMessage struct {
	// Name Nom del a commande
	Name PyDroneCommand `json:"name"`
	// Params paramètres de la commande
	Params interface{} `json:"params"`
}

// PyDroneStatus Statut envoyé par le programme Python
type PyDroneStatus struct {
	Available        bool `json:"available"`
	OnError          bool `json:"on_error"`
	OnGoing          bool `json:"ongoing"`
	InitializedRelay bool `json:"initialized"`
	Connected        bool `json:"connected"`
	ManualFlight     bool `json:"manual"`
	SimMode          bool `json:"sim"`
}

// PyDroneInternalStatus Etat interne d'un drone - on ne sait pas ce qu'il nous envoie spécifiquement
type PyDroneInternalStatus struct {
	Name   string      `json:"id"`
	Type   string      `json:"status"`
	Result interface{} `json:"result"`
}

// DroneStatus Status envoyé par le Drone OLYMPE/ANAFI
type DroneStatus struct {
	Battery bool `json:"battery"`
}

// PyManualCommand Commande à transmettre
type PyManualCommand struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

// PyAutomaticCommand Commande automatique à transmettre
type PyAutomaticCommand struct {
	Name   PyDroneCommand `json:"command"`
	Target string         `json:"name"`
}

// PyDroneFlyingStatus Etat de vol remonté par la partie Python
type PyDroneFlyingStatus int

const (
	// Landed Etat
	Landed PyDroneFlyingStatus = iota
	// TakingOff Etat
	TakingOff PyDroneFlyingStatus = iota
	// Hovering Etat
	Hovering PyDroneFlyingStatus = iota //
	// Flying Etat
	Flying PyDroneFlyingStatus = iota
	// Emergency Etat
	Emergency PyDroneFlyingStatus = iota
	// UserTakeOff Etat
	UserTakeOff PyDroneFlyingStatus = iota
	// MotorRamping Etat
	MotorRamping PyDroneFlyingStatus = iota
	// EmergencyLanding Etat
	EmergencyLanding PyDroneFlyingStatus = iota
)

// PyDroneNavigateHomeStatus Status liés au retour "maison"
type PyDroneNavigateHomeStatus int

const (
	// Available Retour disponible
	Available PyDroneNavigateHomeStatus = iota
	// InProgress En cours de retour
	InProgress PyDroneNavigateHomeStatus = iota
	// Unavailable Indisponible
	Unavailable PyDroneNavigateHomeStatus = iota //
	// Pending Reçu mais en attente
	Pending PyDroneNavigateHomeStatus = iota
)

// PyDroneAlertStatus Alertes remontées par le drone
type PyDroneAlertStatus int

const (
	// None Aucune alerte
	None PyDroneAlertStatus = iota
	// User Alerte utilisateur
	User PyDroneAlertStatus = iota
	// CutOut Alerte "cut-out"
	CutOut PyDroneAlertStatus = iota
	// CriticalBattery Niveau de batterie critique
	CriticalBattery PyDroneAlertStatus = iota
	// LowBattery Niveau de batterie basse
	LowBattery PyDroneAlertStatus = iota
	// TooMuchAngle Trop d'angle (PCMD)
	TooMuchAngle PyDroneAlertStatus = iota
	// AlmostEmtpyBattery Batterie presque vide
	AlmostEmtpyBattery PyDroneAlertStatus = iota
)

// DroneControlSettings Variables de contrôle pour un drone donné
type DroneControlSettings struct {
	DroneName           string  `json:"drone_name"`
	VerticalSpeed       float64 `json:"vertical_speed"`
	HorizontalSpeed     float64 `json:"horizontal_speed"`
	CameraRotationSpeed float64 `json:"camera_speed"`

	MaxTilt          int     `json:"max_tilt"`
	MaxRotationSpeed float64 `json:"max_rotation_speed"`
}

// DroneFlyingStatusMessage Message en provenance de l'unité de contrôle / Automtate Python
type DroneFlyingStatusMessage struct {
	Name   string              `json:"drone_name"`
	Status PyDroneFlyingStatus `json:"status"`
}

// DroneNavigateHomeStatusMessage Message en provenance de l'unité de contrôle / Automtate Python
type DroneNavigateHomeStatusMessage struct {
	Name   string                    `json:"drone_name"`
	Status PyDroneNavigateHomeStatus `json:"status"`
}

// DroneAlertStatusMessage Message en provenance de l'unité de contrôle / Automtate Python
type DroneAlertStatusMessage struct {
	Name   string             `json:"drone_name"`
	Status PyDroneAlertStatus `json:"status"`
}
