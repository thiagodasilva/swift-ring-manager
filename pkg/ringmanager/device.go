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
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func DeviceAdd(w http.ResponseWriter, r *http.Request) {

	var msg DeviceAddRequest
	err := GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", 422)
		return
	}

	// Check the message has devices
	if msg.Name == "" {
		http.Error(w, "Device name missing", http.StatusBadRequest)
		return
	}

	// Create device entry
	device := NewDeviceEntryFromRequest(&msg)

	// Check the node is in the db
	var node *NodeEntry
	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		node, err = NewNodeEntryFromId(tx, msg.NodeId)
		if err == ErrNotFound {
			http.Error(w, "Node id does not exist", http.StatusNotFound)
			return err
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		// Register device
		err = device.Register(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return err
		}

		// Add device to node
		node.DeviceAdd(device.Info.Id)

		// Commit
		err = node.Save(tx)
		if err != nil {
			return err
		}

		// Save drive
		err = device.Save(tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Send back we created it (as long as we did not fail)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(device.Info); err != nil {
		panic(err)
	}

}

func DeviceInformation(w http.ResponseWriter, r *http.Request) {

	// Get device id from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get device information
	info, err := getDeviceInfo(id)
	if err == ErrNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write msg
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		panic(err)
	}

}

func getDeviceInfo(id string) (*DeviceInfoResponse, error) {
	var info *DeviceInfoResponse
	err := db.View(func(tx *bolt.Tx) error {

		// Create a db entry from the id
		entry, err := NewDeviceEntryFromId(tx, id)
		if err != nil {
			return err
		}

		// Create a response from the db entry
		info, err = entry.NewInfoResponse()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return info, nil
}

func DeviceDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Implemented yet"}); err != nil {
		panic(err)
	}
}
