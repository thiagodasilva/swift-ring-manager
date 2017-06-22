package ringmanager

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	// cluster
	Route{
		"ClusterCreate",
		"POST",
		"/clusters",
		ClusterCreate,
	},
	Route{
		"ClusterInfo",
		"GET",
		"/clusters/{id:[A-Fa-f0-9]+}",
		ClusterInfo,
	},
	Route{
		"ClusterList",
		"GET",
		"/clusters",
		ClusterList,
	},

	// Ring
	Route{
		"RingAdd",
		"POST",
		"/rings",
		RingAdd,
	},
	Route{
		"RingInfo",
		"GET",
		"/rings/{id:[A-Fa-f0-9]+}",
		RingInformation,
	},
	Route{
		"RingDelete",
		"DELETE",
		"/rings/{id:[A-Fa-f0-9]+}",
		RingDelete,
	},

	// Node
	Route{
		"NodeAdd",
		"POST",
		"/nodes",
		NodeAdd,
	},
	Route{
		"NodeInfo",
		"GET",
		"/nodes/{id:[A-Fa-f0-9]+}",
		NodeInformation,
	},
	Route{
		"NodeDelete",
		"DELETE",
		"/nodes/{id:[A-Fa-f0-9]+}",
		NodeDelete,
	},

	// Device
	Route{
		"DeviceAdd",
		"POST",
		"/devices",
		DeviceAdd,
	},
	Route{
		"DeviceInfo",
		"GET",
		"/devices/{id:[A-Fa-f0-9]+}",
		DeviceInformation,
	},
	Route{
		"DeviceDelete",
		"DELETE",
		"/devices/{id:[A-Fa-f0-9]+}",
		DeviceDelete,
	},

	// Actions/Tasks
	Route{
		"BuildRing",
		"POST",
		"/buildring/{id:[A-Fa-f0-9]+}",
		BuildRing,
	},
}
