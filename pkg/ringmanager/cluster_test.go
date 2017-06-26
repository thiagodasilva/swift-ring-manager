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
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"bytes"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClusterCreate(t *testing.T) {

	// setup config
	dir, err := ioutil.TempDir("", "ringmanager")
	if err != nil {
		assert.FailNow(t, "Unable to make temp dir")
	}
	defer os.RemoveAll(dir) // clean up

	v := viper.New()
	v.Set("dbfilename", "swift_clusters.db")
	v.Set("ringmanager_dir", dir)

	// setup the server
	router := NewRouter(v)
	ts := httptest.NewServer(router)
	defer ts.Close()

	req := []byte(`{}`)

	// post nothing
	r, err := http.Post(ts.URL+"/clusters", "application/json", bytes.NewBuffer(req))
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusCreated)

	// test json response
	var msg ClusterInfoResponse
	err = GetJsonFromResponse(r, &msg)
	assert.Nil(t, err)
	assert.NotEmpty(t, msg.Id)
	assert.Zero(t, len(msg.Rings))

	// Check data in database
	var entry ClusterEntry
	err = db.View(func(tx *bolt.Tx) error {
		return entry.Unmarshal(tx.Bucket([]byte(BOLTDB_BUCKET_CLUSTER)).Get([]byte(msg.Id)))
	})
	assert.Nil(t, err)

	assert.Equal(t, msg.Id, entry.Info.Id)
	assert.Zero(t, len(entry.Info.Rings))
}
