package main

import (
	"net/http"
)

// Route Chemin Web
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes Ensemble de chemins HTTP
type Routes []Route

// A mettre dans un JSON (et charger via Swagger ?)
var routes = Routes{
	Route{
		"Get Drones Names",
		"GET",
		"/drones",
		GetDronesNames,
	},
	Route{
		"Get Drone Status",
		"GET",
		"/drone/{name}",
		GetDroneStatus,
	},
	Route{
		"Get component integrity",
		"GET",
		"/health",
		GetIntegrity,
	},

	Route{
		"Get Operators",
		"GET",
		"/operators",
		GetOperators,
	},
	Route{
		"Get Boundaries",
		"GET",
		"/map/boundaries",
		GetBoundaries,
	},
	Route{
		"Update the control settings for a drone",
		"POST",
		"/drone/settings/controls",
		UpdateDroneControlSettings,
	},
	Route{
		"Get  all drones settings",
		"GET",
		"/drone/settings/controls",
		GetAllDronesSettings,
	},
	Route{
		"Get control settings of one drone",
		"GET",
		"/drone/{name}/settings/controls",
		GetOneDroneSettings,
	},
	Route{
		"Set autopilot target",
		"POST",
		"/drone/{name}/course/set",
		SetCourse,
	},
	Route{
		"Start the maneuvers",
		"GET",
		"/drone/{name}/autopilot/on",
		SetAutopilotOn,
	},

	Route{
		"Stop the maneuvers",
		"GET",
		"/drone/{name}/autopilot/off",
		SetAutopilotOff,
	},
	Route{
		"Get autopilot Status",
		"GET",
		"/drone/{name}/autopilot",
		GetAutopilot,
	},
	Route{
		"Restart video server",
		"GET",
		"/drone/{name}/video/restart",
		RestartVideoServer,
	},

	Route{
		"Restart video stream",
		"GET",
		"/drone/{name}/stream/restart",
		RestartVideoStream,
	},
	Route{
		"Restart Socket Server",
		"POST",
		"/server/module/restart",
		RestartModule,
	},
	Route{
		"Get Drone Health",
		"GET",
		"/drone/{name}/health",
		GetDroneHealth,
	},
	Route{
		"Execute available command",
		"POST",
		"/command",
		ExecuteRemoteCommand,
	},
	Route{
		"Automated Take off of a drone",
		"GET",
		"/drone/{name}/takeoff",
		TakeOffAutomated,
	},
	Route{
		"Automated go home of a drone",
		"GET",
		"/drone/{name}/home",
		GoHomeAutomated,
	},
	Route{
		"Get the flying status of one drone",
		"GET",
		"/drone/{name}/flying_status",
		GetDroneFlyingStatus,
	},
	Route{
		"Get all flying statuses",
		"GET",
		"/drones/flying_status",
		GetDronesFlyingStatus,
	},
}
