/*
Copyright 2017 The swift-ring-master Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	Route{
		"DownloadRing",
		"GET",
		"/downloadring/{id:[A-Fa-f0-9]+}/{ring}",
		DownloadRing,
	},
}
