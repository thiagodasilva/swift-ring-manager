package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"os"

	"path/filepath"

	"log"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func ClusterCreate(w http.ResponseWriter, r *http.Request) {
	// Create a new ClusterInfo
	entry := NewClusterEntryFromRequest()

	// Add cluster to db
	err := db.Update(func(tx *bolt.Tx) error {
		err := entry.Save(tx)
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
	if err := json.NewEncoder(w).Encode(entry.Info); err != nil {
		panic(err)
	}
}

func ClusterList(w http.ResponseWriter, r *http.Request) {

	var list ClusterListResponse

	// Get all the cluster ids from the DB
	err := db.View(func(tx *bolt.Tx) error {
		var err error

		list.Clusters, err = ClusterEntryList(tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		//logger.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send list back
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		panic(err)
	}
}

func ClusterInfo(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get info from db
	info, err := getClusterInfo(id)
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

func getClusterInfo(id string) (*ClusterInfoResponse, error) {
	var info *ClusterInfoResponse
	err := db.View(func(tx *bolt.Tx) error {

		// Create a db entry from the id
		entry, err := NewClusterEntryFromId(tx, id)
		if err != nil {
			return err
		}

		// Create a response from the db entry
		info, err = entry.NewClusterInfoResponse(tx)
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

func ClusterDelete(w http.ResponseWriter, r *http.Request) {

	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete cluster from db
	err := db.Update(func(tx *bolt.Tx) error {

		// Access cluster entry
		entry, err := NewClusterEntryFromId(tx, id)
		if err == ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return err
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			//return logger.Err(err)
			return err
		}

		err = entry.Delete(tx)
		if err != nil {
			if err == ErrConflict {
				http.Error(w, entry.ConflictString(), http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Update allocator hat the cluster has been removed
	//a.allocator.RemoveCluster(id)

	// Show that the key has been deleted
	//logger.Info("Deleted cluster [%s]", id)

	// Write msg
	w.WriteHeader(http.StatusOK)
}

func BuildRing(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get info from db
	clusterInfo, err := getClusterInfo(id)
	if err == ErrNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//TODO define better location to store rings
	clusterPath := filepath.Join("/tmp", clusterInfo.Id)
	os.Mkdir(clusterPath, 0774)

	for _, ringId := range clusterInfo.Rings {
		ringInfo, err := getRingInfo(ringId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ringbuilderName := ringInfo.Name + ".builder"
		ringBuilderPath := filepath.Join(clusterPath, ringbuilderName)
		if _, err := os.Stat(ringBuilderPath); os.IsNotExist(err) {
			// create builder
			cmdRingBuilder := "/usr/bin/swift-ring-builder"
			cmdArgs := []string{ringBuilderPath, "create", "10", "3", "1"}
			out, err := exec.Command(cmdRingBuilder, cmdArgs...).CombinedOutput()
			if err != nil {
				output := fmt.Sprintf("%s", out)
				http.Error(w, output, http.StatusInternalServerError)
				return
			}

			// Add Nodes/Devices to ring
			for _, nodeId := range ringInfo.Nodes {
				n, err := getNodeInfo(nodeId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				for _, deviceId := range n.Devices {
					d, err := getDeviceInfo(deviceId)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					//swift-ring-builder object.builder add r1z1-127.0.0.1:6010/sdb1 1
					deviceArg := fmt.Sprintf("r%dz%d-%s:%s/%s", n.Region, n.Zone, n.Ip, n.Port, d.Name)
					weightArg := fmt.Sprintf("%d", d.Weight.Target)
					cmdArgs = []string{ringBuilderPath, "add", deviceArg, weightArg}
					out, err = exec.Command(cmdRingBuilder, cmdArgs...).CombinedOutput()
					if err != nil {
						log.Printf("output %s", out)
						output := fmt.Sprintf("%s", out)
						http.Error(w, output, http.StatusInternalServerError)
						return
					}
				}
			}

			// rebalance
			cmdArgs = []string{ringBuilderPath, "rebalance"}
			out, err = exec.Command(cmdRingBuilder, cmdArgs...).CombinedOutput()
			if err != nil {
				output := fmt.Sprintf("%s", out)
				http.Error(w, output, http.StatusInternalServerError)
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
