package main

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
