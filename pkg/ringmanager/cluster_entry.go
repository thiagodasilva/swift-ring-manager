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

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/lpabon/godbc"
)

type ClusterEntry struct {
	Info ClusterInfoResponse
}

func ClusterEntryList(tx *bolt.Tx) ([]string, error) {

	list := EntryKeys(tx, BOLTDB_BUCKET_CLUSTER)
	if list == nil {
		return nil, ErrAccessList
	}
	return list, nil
}

func NewClusterEntry() *ClusterEntry {
	entry := &ClusterEntry{}
	entry.Info.Rings = make(sort.StringSlice, 0)

	return entry
}

func NewClusterEntryFromRequest() *ClusterEntry {
	entry := NewClusterEntry()
	entry.Info.Id = GenUUID()

	return entry
}

func NewClusterEntryFromId(tx *bolt.Tx, id string) (*ClusterEntry, error) {

	entry := NewClusterEntry()
	err := EntryLoad(tx, entry, id)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (c *ClusterEntry) BucketName() string {
	return BOLTDB_BUCKET_CLUSTER
}

func (c *ClusterEntry) Save(tx *bolt.Tx) error {
	godbc.Require(tx != nil)
	godbc.Require(len(c.Info.Id) > 0)

	return EntrySave(tx, c, c.Info.Id)
}

func (c *ClusterEntry) NewClusterInfoResponse(tx *bolt.Tx) (*ClusterInfoResponse, error) {

	info := &ClusterInfoResponse{}
	*info = c.Info

	return info, nil
}

func (c *ClusterEntry) Delete(tx *bolt.Tx) error {
	godbc.Require(tx != nil)

	// Check if the cluster still has nodes or volumes
	if len(c.Info.Rings) > 0 {
		// TODO: logger.Warning(c.ConflictString())
		return ErrConflict
	}

	return EntryDelete(tx, c, c.Info.Id)
}

func (c *ClusterEntry) ConflictString() string {
	return fmt.Sprintf("Unable to delete cluster [%v] because it contains rings", c.Info.Id)
}

func (c *ClusterEntry) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(*c)

	return buffer.Bytes(), err
}

func (c *ClusterEntry) Unmarshal(buffer []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	err := dec.Decode(c)
	if err != nil {
		return err
	}

	// Make sure to setup slices if nil
	if c.Info.Rings == nil {
		c.Info.Rings = make(sort.StringSlice, 0)
	}

	return nil
}

func (c *ClusterEntry) RingAdd(id string) {
	c.Info.Rings = append(c.Info.Rings, id)
	c.Info.Rings.Sort()
}

func (c *ClusterEntry) RingDelete(id string) {
	c.Info.Rings = SortedStringsDelete(c.Info.Rings, id)
}
