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

import "sort"

// TODO: not sure we need this yet
type EntryState string

type Entry struct {
	State EntryState
}

type ClusterInfoResponse struct {
	Id    string           `json:"id"`
	Rings sort.StringSlice `json:"rings"`
}

type ClusterListResponse struct {
	Clusters []string `json:"clusters"`
}

type RingAddRequest struct {
	ClusterId string `json:"cluster"`
	Name      string `json:"name"`
}

type RingInfo struct {
	RingAddRequest
	Id string `json:"id"`
}

type RingInfoResponse struct {
	RingInfo
	Nodes sort.StringSlice `json:"nodes"`
}

type NodeAddRequest struct {
	RingId          string `json:"ring"`
	Region          int    `json:"region"`
	Zone            int    `json:"zone"`
	Ip              string `json:"ip"`
	ReplicationIP   string `json:"replicationIP"`
	Port            string `json:"port"`
	ReplicationPort string `json:"replicationPort"`
}

type NodeInfo struct {
	NodeAddRequest
	Id string `json:"id"`
}

type NodeInfoResponse struct {
	NodeInfo
	Devices sort.StringSlice `json:"devices"`
}

type Device struct {
	Name string `json:"name"`
	Meta string `json:"meta"`
}

type DeviceAddRequest struct {
	Device
	Weight uint64 `json:"weight"`
	NodeId string `json:"node"`
}

type DeviceInfo struct {
	Device
	Weight DeviceWeight `json:"weight"`
	Id     string       `json:"id"`
}

type DeviceInfoResponse struct {
	DeviceInfo
}

type DeviceWeight struct {
	Current uint64 `json:"current"`
	Target  uint64 `json:"target"`
}
