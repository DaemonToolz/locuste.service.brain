package main

// FlightCoordinate struct
type FlightCoordinate struct {
	Name string  `json:"name"`
	Lat  float64 `json:"latitude"`
	Lon  float64 `json:"longitude"`
}

// Boundaries struct
type Boundaries struct {
	MinLat float64 `json:"min_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLat float64 `json:"max_lat"`
	MaxLon float64 `json:"max_lon"`
}

// DroneFlightCoordinates Coordonnées de vol
type DroneFlightCoordinates struct {
	DroneName string            `json:"drone_name"`
	Component *FlightCoordinate `json:"coordinates"`
	Metadata  *NodeMetaData     `json:"metadata"`
}

// NodeMetaData Métadonnées de vol
type NodeMetaData struct {
	Name     string           `json:"street_name"`
	Distance float64          `json:"distance"`
	Altitude float64          `json:"altitude"`
	Previous FlightCoordinate `json:"previous"`
	Next     FlightCoordinate `json:"next"`
}

// SchedulerSummarizedData Informations réduites pour les autopilotes (utilisé en communication Brain <=> Scheduler)
type SchedulerSummarizedData struct {
	DroneName   string `json:"drone_name"`
	IsActive    bool   `json:"is_active"`
	IsManual    bool   `json:"is_manual"`
	IsSimulated bool   `json:"is_simulated"`
	IsRunning   bool   `json:"is_running"`
	IsReady     bool   `json:"is_ready"`
	IsBusy      bool   `json:"is_busy"`
}
