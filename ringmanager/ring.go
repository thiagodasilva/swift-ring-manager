package ringmanager

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func RingAdd(w http.ResponseWriter, r *http.Request) {
	var msg RingAddRequest
	err := GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", 422)
		return
	}

	// check information in JSON request
	if len(msg.ClusterId) == 0 {
		http.Error(w, "Cluster id missing", http.StatusBadRequest)
		return
	}

	if len(msg.Name) == 0 {
		http.Error(w, "Ring name missing", http.StatusBadRequest)
		return
	}

	// create a ring entry
	ring := NewRingEntryFromRequest(&msg)

	var cluster *ClusterEntry
	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		cluster, err = NewClusterEntryFromId(tx, msg.ClusterId)
		if err == ErrNotFound {
			http.Error(w, "Cluster id does not exist", http.StatusNotFound)
			return err
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		// Register ring
		err = ring.Register(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return err
		}

		// add ring to cluster
		cluster.RingAdd(ring.Info.Id)

		// save cluster
		err = cluster.Save(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		// save ring
		err = ring.Save(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	//logger.Info("Added node " + node.Info.Id)
	// Send back we created it (as long as we did not fail)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ring.Info); err != nil {
		panic(err)
	}
}

func RingInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// get ring information
	info, err := getRingInfo(id)
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

func getRingInfo(id string) (*RingInfoResponse, error) {
	var info *RingInfoResponse
	err := db.View(func(tx *bolt.Tx) error {

		// Create a db entry from the id
		entry, err := NewRingEntryFromId(tx, id)
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

func RingDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Implemented yet"}); err != nil {
		panic(err)
	}

}
