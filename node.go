package main

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func NodeAdd(w http.ResponseWriter, r *http.Request) {
	var msg NodeAddRequest
	err := GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", 422)
		return
	}

	// check information in JSON request
	if len(msg.RingId) == 0 {
		http.Error(w, "Ring id missing", http.StatusBadRequest)
		return
	}

	if len(msg.Ip) == 0 {
		http.Error(w, "Ip missing", http.StatusBadRequest)
		return
	}

	if len(msg.Port) == 0 {
		http.Error(w, "Port missing", http.StatusBadRequest)
		return
	}

	// create a ring entry
	node := NewNodeEntryFromRequest(&msg)

	var ring *RingEntry
	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		ring, err = NewRingEntryFromId(tx, msg.RingId)
		if err == ErrNotFound {
			http.Error(w, "Ring id does not exist", http.StatusNotFound)
			return err
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		// Register node
		err = node.Register(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return err
		}

		// add node to ring
		ring.NodeAdd(node.Info.Id)

		// save ring
		err = ring.Save(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		// save node
		err = node.Save(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if err := json.NewEncoder(w).Encode(node.Info); err != nil {
		panic(err)
	}
}

func NodeInformation(w http.ResponseWriter, r *http.Request) {

	// Get node id from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get Node information
	info, err := getNodeInfo(id)
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

func getNodeInfo(id string) (*NodeInfoResponse, error) {
	var info *NodeInfoResponse
	err := db.View(func(tx *bolt.Tx) error {

		// Create a db entry from the id
		entry, err := NewNodeEntryFromId(tx, id)
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

func NodeDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Implemented yet"}); err != nil {
		panic(err)
	}
}
