package main

import "sync"

var droneSettingMutex sync.Mutex

// DroneSettings Configurations de vol pour chaque drone
var DroneSettings map[string]DroneControlSettings

// https://forum.developer.parrot.com/t/work-with-ardrone3-piloting-pcmd/10186
// https://developer.parrot.com/docs/olympe/arsdkng_ardrone3_piloting.html

func initDroneSettings() {
	DroneSettings = make(map[string]DroneControlSettings)
	for _, name := range ExtractDroneNames() {
		AddOrUpdateDroneSettings(name, DroneControlSettings{
			DroneName:           name,
			CameraRotationSpeed: 0.20, // en degrés/s
			VerticalSpeed:       0.20, // m/s
			HorizontalSpeed:     0.20, // m/s
			MaxTilt:             15,   // 15% - sera utilisé lorsque l'on intégrera la version PCMD ~ utilisateur avancé
			MaxRotationSpeed:    0.20, // en degrés/s
		})
	}
}

// AddOrUpdateDroneSettings Mise à jour des paramètres de vol d'un drone
func AddOrUpdateDroneSettings(name string, settings DroneControlSettings) {
	droneSettingMutex.Lock()
	DroneSettings[name] = settings
	droneSettingMutex.Unlock()
}

// GetDroneSettings Récupère les paramètres d'un drone
func GetDroneSettings(name string) DroneControlSettings {
	droneSettingMutex.Lock()
	defer droneSettingMutex.Unlock()
	return DroneSettings[name]
}
