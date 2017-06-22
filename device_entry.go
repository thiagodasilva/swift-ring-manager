package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/lpabon/godbc"
)

type DeviceEntry struct {
	Entry

	Info   DeviceInfo
	NodeId string
}

func DeviceList(tx *bolt.Tx) ([]string, error) {

	list := EntryKeys(tx, BOLTDB_BUCKET_DEVICE)
	if list == nil {
		return nil, ErrAccessList
	}
	return list, nil
}

func NewDeviceEntry() *DeviceEntry {
	entry := &DeviceEntry{}
	return entry
}

func NewDeviceEntryFromRequest(req *DeviceAddRequest) *DeviceEntry {
	godbc.Require(req != nil)

	device := NewDeviceEntry()
	device.NodeId = req.NodeId
	device.Info.Id = GenUUID()
	device.Info.Name = req.Name
	device.Info.Meta = req.Meta
	device.Info.Weight.Target = req.Weight

	return device
}

func NewDeviceEntryFromId(tx *bolt.Tx, id string) (*DeviceEntry, error) {
	godbc.Require(tx != nil)

	entry := NewDeviceEntry()
	err := EntryLoad(tx, entry, id)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (d *DeviceEntry) registerKey() string {
	return "DEVICE" + d.NodeId + d.Info.Name
}

func (d *DeviceEntry) Register(tx *bolt.Tx) error {
	godbc.Require(tx != nil)

	val, err := EntryRegister(tx,
		d,
		d.registerKey(),
		[]byte(d.Info.Id))
	if err == ErrKeyExists {

		// Now check if the node actually exists.  This only happens
		// when the application crashes and it doesn't clean up stale
		// registrations.
		conflictId := string(val)
		_, err := NewDeviceEntryFromId(tx, conflictId)
		if err == ErrNotFound {
			// (stale) There is actually no conflict, we can allow
			// the registration
			return nil
		} else if err != nil {
			//return logger.Err(err)
			return err
		}

		return fmt.Errorf("Device %v is already used on node %v by device %v",
			d.Info.Name,
			d.NodeId,
			conflictId)

	} else if err != nil {
		return err
	}

	return nil
}

func (d *DeviceEntry) Deregister(tx *bolt.Tx) error {
	godbc.Require(tx != nil)

	err := EntryDelete(tx, d, d.registerKey())
	if err != nil {
		return err
	}

	return nil
}

func (d *DeviceEntry) BucketName() string {
	return BOLTDB_BUCKET_DEVICE
}

func (d *DeviceEntry) Save(tx *bolt.Tx) error {
	godbc.Require(tx != nil)
	godbc.Require(len(d.Info.Id) > 0)

	return EntrySave(tx, d, d.Info.Id)

}

func (d *DeviceEntry) NewInfoResponse() (*DeviceInfoResponse, error) {

	info := &DeviceInfoResponse{}
	info.Id = d.Info.Id
	info.Name = d.Info.Name
	info.Meta = d.Info.Meta
	info.Weight.Current = d.Info.Weight.Current
	info.Weight.Target = d.Info.Weight.Target

	return info, nil
}

func (d *DeviceEntry) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(*d)

	return buffer.Bytes(), err
}

func (d *DeviceEntry) Unmarshal(buffer []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	err := dec.Decode(d)
	if err != nil {
		return err
	}

	return nil
}
