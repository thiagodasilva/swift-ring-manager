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
	"net/http"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/heketi/rest"
)

var db *bolt.DB
var dbReadOnly bool
var dbfilename = "swift_clusters.db"
var ringmanager_dir = "/tmp/ringmanager"

const (
	ASYNC_ROUTE           = "/queue"
	BOLTDB_BUCKET_CLUSTER = "CLUSTER"
	BOLTDB_BUCKET_RING    = "RING"
	BOLTDB_BUCKET_NODE    = "NODE"
	BOLTDB_BUCKET_DEVICE  = "DEVICE"
)

type App struct {
	asyncManager *rest.AsyncHttpManager
}

func NewRouter() *mux.Router {

	var err error
	dbFilePath := filepath.Join(ringmanager_dir, dbfilename)

	// Setup BoltDB database
	db, err = bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		//logger.Warning("Unable to open database.  Retrying using read only mode")

		// Try opening as read-only
		db, err = bolt.Open(dbfilename, 0666, &bolt.Options{
			ReadOnly: true,
		})
		if err != nil {
			//logger.LogError("Unable to open database: %v", err)
		}
		dbReadOnly = true
	} else {
		err = db.Update(func(tx *bolt.Tx) error {
			// Create Cluster Bucket
			_, err := tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_CLUSTER))
			if err != nil {
				//logger.LogError("Unable to create cluster bucket in DB")
				return err
			}

			// Create Ring Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_RING))
			if err != nil {
				//logger.LogError("Unable to create ring bucket in DB")
				return err
			}
			// Create Node Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_NODE))
			if err != nil {
				//logger.LogError("Unable to create node bucket in DB")
				return err
			}

			// Create Device Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_DEVICE))
			if err != nil {
				//logger.LogError("Unable to create device bucket in DB")
				return err
			}

			return nil

		})
		if err != nil {
			//logger.Err(err)
		}
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
