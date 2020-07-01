package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var loadErrorMsg string = "Unable to load the message"

var failStructResponse = struct {
	Success bool `json:"success"`
}{false}

var successStructResponse = struct {
	Success bool `json:"success"`
}{true}

// GetDronesNames retourne le nom de tous les drones disponibles
func GetDronesNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(ExtractDroneNames()); err != nil {
		log.Printf(err.Error())
		panic(err)
	}
}

// GetDroneStatus retourne le dernier état enregistré par un drone
func GetDroneStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := json.NewEncoder(w).Encode(GetDroneStatuses(vars["name"])); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetIntegrity Récupère l'intégrité globale de l'application
func GetIntegrity(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(GlobalStatuses); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetOperators Récupère l'intégralité des opérateurs
func GetOperators(w http.ResponseWriter, r *http.Request) {
	operators := make([]Operator, len(OperatorsInCharge))
	index := 0
	for _, value := range OperatorsInCharge {
		operators[index] = value
		index++
	}

	if err := json.NewEncoder(w).Encode(operators); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// SetCourse Permet de mettre à jour la cible d'un des autopilots
func SetCourse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post FlightCoordinate
	err := decoder.Decode(&post)

	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
	}

	UpdateTarget(post)
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// ExecuteRemoteCommand Permet d'exécuter une commande spécifique
func ExecuteRemoteCommand(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post RemoteManualCommand
	err := decoder.Decode(&post)

	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
	}

	RedirectCommand(post)
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// RestartVideoServer Permet de redémarrer le serveur vidéo dédié à un drone
func RestartVideoServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	conv, err := strconv.Atoi(strings.SplitN(name, ".", -1)[2])
	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
		return
	}
	videoPort := 7000 + conv
	startVideoServer(name, videoPort)
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// RestartVideoStream Permet de redémarrer le stream vidéo dédié à un drone
func RestartVideoStream(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	conv, err := strconv.Atoi(strings.SplitN(name, ".", -1)[2])
	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
		return
	}
	videoPort := 7000 + conv
	startFfmpegStream(name, videoPort)

	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// RestartZMQServer Permet de redémarrer le serveur de messages ZMQ
func RestartZMQServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	conv, err := strconv.Atoi(strings.SplitN(name, ".", -1)[2])
	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
		return
	}
	zmqPort := 7300 + conv
	StartZMQServer(name, zmqPort)
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetDroneHealth Récupère les différents indicateurs relatif au drone (e.g. composants externes)
func GetDroneHealth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := json.NewEncoder(w).Encode(GetExtCompStatus(vars["name"])); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetBoundaries Récupère les limites de la carte (google map)
func GetBoundaries(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(flightSchedulerRPC.MapBoundaries); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetAutopilot Récupère l'éatt d'un pilote automatique
func GetAutopilot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := json.NewEncoder(w).Encode(GetAutopilotStatus(vars["name"])); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// SetAutopilotOn Active le pilote automatique pour le drone
func SetAutopilotOn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	autoPilot := GetAutopilotStatus(vars["name"])

	autoPilot.IsActive = true
	UpdateAutopilot(autoPilot)

	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// SetAutopilotOff Désactive le pilote automatique pour le drone
func SetAutopilotOff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	autoPilot := GetAutopilotStatus(vars["name"])

	autoPilot.IsActive = false
	UpdateAutopilot(autoPilot)

	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// RestartModule Permet de redémarrer un module (version HTTP)
func RestartModule(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post Module
	err := decoder.Decode(&post)

	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
	}

	ModuleRestart(post)
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// UpdateDroneControlSettings Mise à jour des configurations de vol
func UpdateDroneControlSettings(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post DroneControlSettings
	err := decoder.Decode(&post)

	if err != nil {
		if err := json.NewEncoder(w).Encode(failStructResponse); err != nil {
			failOnError(err, loadErrorMsg)
			panic(err)
		}
	}

	AddOrUpdateDroneSettings(post.DroneName, post)
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetOneDroneSettings Récupère les configurations de vol pour un drone
func GetOneDroneSettings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := json.NewEncoder(w).Encode(GetDroneSettings(vars["name"])); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetAllDronesSettings  Récupère les configurations de vol de tous les drones
func GetAllDronesSettings(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(DroneSettings); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// TakeOffAutomated Active le pilote automatique pour le drone
func TakeOffAutomated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	SendTakeoffCommandTo(vars["name"])
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GoHomeAutomated Instruction
func GoHomeAutomated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	SendGoHomeCommandTo(vars["name"])
	if err := json.NewEncoder(w).Encode(successStructResponse); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}
}

// GetDroneFlyingStatus Récupère les états d'un drone'
func GetDroneFlyingStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if err := json.NewEncoder(w).Encode(GetFlyingStatus(vars["name"])); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}

}

// GetDronesFlyingStatus Récupère les états de tous les drones
func GetDronesFlyingStatus(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(FlyingStatuses); err != nil {
		failOnError(err, loadErrorMsg)
		panic(err)
	}

}
