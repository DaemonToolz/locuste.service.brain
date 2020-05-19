package main


/*
	Regroupement des éléments purements dédiés à la partie Pyton
*/


// IdentificationRequest Requête d'identification
type IdentificationRequest struct{
	Name string `json:"name"`
	VideoPort int `json:"video_port"`
	IP string `json:"ip"`
	Connected bool `json:"connected"`
	ManualFlight bool `json:"manual"`
	SimMode bool `json:"sim"`
	Position interface{} `json:"position"`
}
   
// NavigationType Mode de navigation, si Manual, alors manuel, sinon automatique
type NavigationType struct{
	Manual bool `json:"manual"`
}

// PyDroneCommandMessage Ordre à envoyer aux drones
type PyDroneCommandMessage struct {
	// Name Nom del a commande
	Name PyDroneCommand `json:"name"`
	// Params paramètres de la commande
	Params interface{}  `json:"params"`
}

// PyDroneStatus Statut envoyé par le programme Python
type PyDroneStatus struct {
	Available bool `json:"available"` 
	OnError bool `json:"on_error"`
	OnGoing bool `json:"ongoing"`
	InitializedRelay bool `json:"initialized"`
	Connected bool `json:"connected"`
	ManualFlight bool `json:"manual"`
	SimMode bool `json:"sim"`
}

// PyDroneInternalStatus Etat interne d'un drone - on ne sait pas ce qu'il nous envoie spécifiquement
type PyDroneInternalStatus struct {
	Name string `json:"id"`
	Type string `json:"status"`
	Result interface{} `json:"result"`
}

// DroneStatus Status envoyé par le Drone OLYMPE/ANAFI
type DroneStatus struct {
	Battery bool `json:"battery"` 
}