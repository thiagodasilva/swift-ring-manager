package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/heketi/rest"
)

var db *bolt.DB
var dbReadOnly bool
var dbfilename = "anira.db"

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

func main() {

	var err error

	// Setup BoltDB database
	db, err = bolt.Open(dbfilename, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		//logger.Warning("Unable to open database.  Retrying using read only mode")

		// Try opening as read-only
		db, err = bolt.Open(dbfilename, 0666, &bolt.Options{
			ReadOnly: true,
		})
		if err != nil {
			//logger.LogError("Unable to open database: %v", err)
			return
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
			return
		}
	}

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
