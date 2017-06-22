package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/heketi/utils"
	"github.com/lpabon/godbc"
)

type RingEntry struct {
	Entry

	Info  RingInfo
	Nodes sort.StringSlice
}

func NewRingEntry() *RingEntry {
	entry := &RingEntry{}
	entry.Nodes = make(sort.StringSlice, 0)

	return entry
}

func NewRingEntryFromRequest(req *RingAddRequest) *RingEntry {
	godbc.Require(req != nil)

	ring := NewRingEntry()
	ring.Info.Id = GenUUID()
	ring.Info.Name = req.Name
	ring.Info.ClusterId = req.ClusterId

	return ring
}

func NewRingEntryFromId(tx *bolt.Tx, id string) (*RingEntry, error) {
	godbc.Require(tx != nil)

	entry := NewRingEntry()
	err := EntryLoad(tx, entry, id)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *RingEntry) registerKey() string {
	return "RING" + r.Info.ClusterId + r.Info.Name
}

func (r *RingEntry) Register(tx *bolt.Tx) error {

	val, err := EntryRegister(tx, r, r.registerKey(), []byte(r.Info.Id))

	if err == ErrKeyExists {
		// Now check if the ring actually exists.  This only happens
		// when the application crashes and it doesn't clean up stale
		// registrations.
		conflictId := string(val)
		_, err := NewRingEntryFromId(tx, conflictId)
		if err == ErrNotFound {
			// (stale) There is actually no conflict, we can allow
			// the registration
			return nil
		} else if err != nil {
			//return logger.Err(err)
			return err
		}

		// Return that we found a conflict
		return fmt.Errorf("Ring %v already used by cluster with id %v\n",
			r.Info.Name, conflictId)
	} else if err != nil {
		return err
	}

	return nil

}

func (r *RingEntry) Deregister(tx *bolt.Tx) error {

	err := EntryDelete(tx, r, r.registerKey())
	if err != nil {
		return err
	}

	return nil
}

func (r *RingEntry) BucketName() string {
	return BOLTDB_BUCKET_RING
}

func (r *RingEntry) Save(tx *bolt.Tx) error {
	godbc.Require(tx != nil)
	godbc.Require(len(r.Info.Id) > 0)

	return EntrySave(tx, r, r.Info.Id)

}

func (r *RingEntry) IsDeleteOk() bool {
	// Check if the nodes still has drives
	if len(r.Nodes) > 0 {
		return false
	}
	return true
}

func (r *RingEntry) ConflictString() string {
	return fmt.Sprintf("Unable to delete ring [%v] because it contains nodes", r.Info.Id)
}
func (r *RingEntry) Delete(tx *bolt.Tx) error {
	godbc.Require(tx != nil)

	// Check if the nodes still has nodes
	if !r.IsDeleteOk() {
		//logger.Warning(r.ConflictString())
		return ErrConflict
	}

	return EntryDelete(tx, r, r.Info.Id)
}

func (r *RingEntry) NewInfoResponse() (*RingInfoResponse, error) {
	info := &RingInfoResponse{}
	info.ClusterId = r.Info.ClusterId
	info.Id = r.Info.Id
	info.Name = r.Info.Name
	//info.Nodes = make(sort.StringSlice, 0)
	info.Nodes = r.Nodes
	return info, nil
}

func (r *RingEntry) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(*r)

	return buffer.Bytes(), err
}

func (r *RingEntry) Unmarshal(buffer []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	err := dec.Decode(r)
	if err != nil {
		return err
	}

	// Make sure to setup nodes if nil
	if r.Nodes == nil {
		r.Nodes = make(sort.StringSlice, 0)
	}

	return nil
}

func (r *RingEntry) NodeAdd(id string) {
	godbc.Require(!utils.SortedStringHas(r.Nodes, id))

	r.Nodes = append(r.Nodes, id)
	r.Nodes.Sort()
}

func (r *RingEntry) NodeDelete(id string) {
	r.Nodes = SortedStringsDelete(r.Nodes, id)
}
