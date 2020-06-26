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
			MaxRotationSpeed:    0.20, // en degrés/s

			MaxTilt:     0, // 5% - sera utilisé lorsque l'on intégrera la version PCMD ~ utilisateur avancé
			MaxThrottle: 0, // 5%
			MaxYaw:      5, // 5%
		})
	}
}

// AddOrUpdateDroneSettings Mise à jour des paramètres de vol d'un drone
func AddOrUpdateDroneSettings(name string, settings DroneControlSettings) {
	var original *DroneControlSettings
	if _, ok := DroneSettings[name]; ok {
		droneSettingMutex.Lock()
		tmp := DroneSettings[name]
		droneSettingMutex.Unlock()
		original = &tmp
	}

	if original != nil {
		if settings.CameraRotationSpeed > 0 {
			original.CameraRotationSpeed = settings.CameraRotationSpeed
		}
		if settings.VerticalSpeed > 0 {
			original.VerticalSpeed = settings.VerticalSpeed
		}
		if settings.HorizontalSpeed > 0 {
			original.HorizontalSpeed = settings.HorizontalSpeed
		}
		if settings.MaxRotationSpeed > 0 {
			original.MaxRotationSpeed = settings.MaxRotationSpeed
		}

		if settings.MaxTilt > 0 {
			original.MaxTilt = settings.MaxTilt
		}
		if settings.MaxThrottle > 0 {
			original.MaxThrottle = settings.MaxThrottle
		}
		if settings.MaxYaw > 0 {
			original.MaxYaw = settings.MaxYaw
		}

		droneSettingMutex.Lock()
		DroneSettings[name] = *original
		droneSettingMutex.Unlock()

	} else {
		droneSettingMutex.Lock()
		DroneSettings[name] = settings
		droneSettingMutex.Unlock()
	}

}

// GetDroneSettings Récupère les paramètres d'un drone
func GetDroneSettings(name string) DroneControlSettings {
	droneSettingMutex.Lock()
	defer droneSettingMutex.Unlock()
	return DroneSettings[name]
}
