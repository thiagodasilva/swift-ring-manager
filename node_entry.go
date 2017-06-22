package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/lpabon/godbc"
)

type NodeEntry struct {
	Entry

	Info    NodeInfo
	Devices sort.StringSlice
}

func NewNodeEntry() *NodeEntry {
	entry := &NodeEntry{}
	entry.Devices = make(sort.StringSlice, 0)

	return entry
}

func NewNodeEntryFromRequest(req *NodeAddRequest) *NodeEntry {
	godbc.Require(req != nil)

	node := NewNodeEntry()
	node.Info.Id = GenUUID()
	node.Info.RingId = req.RingId
	node.Info.Ip = req.Ip
	node.Info.Port = req.Port

	// default Region to 1
	node.Info.Region = 1
	if req.Region > 1 {
		node.Info.Region = req.Region
	}

	// default Zone to 1
	node.Info.Zone = 1
	if req.Zone > 1 {
		node.Info.Zone = req.Zone
	}

	// default replication ip to same as node ip
	node.Info.ReplicationIP = req.Ip
	if req.ReplicationIP != "" {
		node.Info.ReplicationIP = req.ReplicationIP
	}

	// default replication port to same as node port
	node.Info.ReplicationPort = req.Port
	if req.ReplicationPort != "" {
		node.Info.ReplicationPort = req.ReplicationPort
	}

	return node
}

func NewNodeEntryFromId(tx *bolt.Tx, id string) (*NodeEntry, error) {
	godbc.Require(tx != nil)

	entry := NewNodeEntry()
	err := EntryLoad(tx, entry, id)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (n *NodeEntry) registerKey() string {
	return "NODE" + n.Info.RingId + n.Info.Id
}

func (n *NodeEntry) Register(tx *bolt.Tx) error {

	val, err := EntryRegister(tx, n, n.registerKey(), []byte(n.Info.Id))

	if err == ErrKeyExists {
		// Now check if the node actually exists.  This only happens
		// when the application crashes and it doesn't clean up stale
		// registrations.
		conflictId := string(val)
		_, err := NewNodeEntryFromId(tx, conflictId)
		if err == ErrNotFound {
			// (stale) There is actually no conflict, we can allow
			// the registration
			return nil
		} else if err != nil {
			//return logger.Err(err)
			return err
		}

		// Return that we found a conflict
		return fmt.Errorf("Node %v already used by ring with id %v\n",
			n.Info.Ip, conflictId)
	} else if err != nil {
		return err
	}

	return nil

}
func (n *NodeEntry) BucketName() string {
	return BOLTDB_BUCKET_NODE
}

func (n *NodeEntry) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(*n)

	return buffer.Bytes(), err
}

func (n *NodeEntry) Unmarshal(buffer []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	err := dec.Decode(n)
	if err != nil {
		return err
	}

	// Make sure to setup arrays if nil
	if n.Devices == nil {
		n.Devices = make(sort.StringSlice, 0)
	}

	return nil
}

func (n *NodeEntry) Save(tx *bolt.Tx) error {
	godbc.Require(tx != nil)
	godbc.Require(len(n.Info.Id) > 0)

	return EntrySave(tx, n, n.Info.Id)

}

func (n *NodeEntry) NewInfoResponse() (*NodeInfoResponse, error) {

	info := &NodeInfoResponse{}
	info.RingId = n.Info.RingId
	info.Id = n.Info.Id
	info.Region = n.Info.Region
	info.Zone = n.Info.Zone
	info.Ip = n.Info.Ip
	info.Port = n.Info.Port
	info.ReplicationIP = n.Info.ReplicationIP
	info.ReplicationPort = n.Info.ReplicationPort
	info.Devices = n.Devices

	return info, nil
}

func (n *NodeEntry) DeviceAdd(id string) {
	godbc.Require(!SortedStringHas(n.Devices, id))

	n.Devices = append(n.Devices, id)
	n.Devices.Sort()
}

func (n *NodeEntry) DeviceDelete(id string) {
	n.Devices = SortedStringsDelete(n.Devices, id)
}
