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

var ts *httptest.Server

func setupDatabase(t *testing.T) (string, func(t *testing.T)) {

	// Setup
	dir, err := ioutil.TempDir("", "ringmanager")
	if err != nil {
		panic("Unable to make temp dir")
	}

	v := viper.New()
	v.Set("dbfilename", "swift_clusters.db")
	v.Set("ringmanager_dir", dir)

	router := NewRouter(v)
	ts = httptest.NewServer(router)

	clusterId := setupCluster(t)

	return clusterId, func(t *testing.T) {
		// teardown
		ts.Close()
		os.RemoveAll(dir) // clean up
	}
}

func setupCluster(t *testing.T) string {

	body := []byte(`{}`)

	// post nothing
	r, err := http.Post(ts.URL+"/clusters", "application/json", bytes.NewBuffer(body))
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
	return msg.Id
}

func TestClusterList(t *testing.T) {

	// setup and teardown test case
	_, tearDown := setupDatabase(t)
	defer tearDown(t)

	r, err := http.Get(ts.URL + "/clusters")
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusOK)
	assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

	// read response
	var msg ClusterListResponse
	err = GetJsonFromResponse(r, &msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(msg.Clusters))
}

func TestClusterInfoIdNotFound(t *testing.T) {
	// setup and teardown test case
	_, tearDown := setupDatabase(t)
	defer tearDown(t)

	r, err := http.Get(ts.URL + "/clusters/12345")
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusNotFound)
}

func TestClusterInfo(t *testing.T) {
	// setup and teardown test case
	id, tearDown := setupDatabase(t)
	defer tearDown(t)

	r, err := http.Get(ts.URL + "/clusters/" + id)
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusOK)
	assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

	var msg ClusterInfoResponse
	err = GetJsonFromResponse(r, &msg)
	assert.Nil(t, err)
	assert.Equal(t, id, msg.Id)
	assert.Zero(t, len(msg.Rings))
}

func TestClusterDelete(t *testing.T) {

	// setup and teardown test case
	id, tearDown := setupDatabase(t)
	defer tearDown(t)

	req, err := http.NewRequest("DELETE", ts.URL+"/clusters/"+id, nil)
	assert.Nil(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	r, err := http.Get(ts.URL + "/clusters")
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusOK)
	assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

	// read response
	var msg ClusterListResponse
	err = GetJsonFromResponse(r, &msg)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(msg.Clusters))
}

func TestClusterDeleteIdNotFound(t *testing.T) {

	// setup and teardown test case
	_, tearDown := setupDatabase(t)
	defer tearDown(t)

	req, err := http.NewRequest("DELETE", ts.URL+"/clusters/12345", nil)
	assert.Nil(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusNotFound)
}

func TestClusterDeleteRingExists(t *testing.T) {

	// setup and teardown test case
	id, tearDown := setupDatabase(t)
	defer tearDown(t)

	body := []byte(`{"name":"account", "cluster":"` + id + `"}`)

	// create ring
	r, err := http.Post(ts.URL+"/rings", "application/json", bytes.NewBuffer(body))
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusCreated)

	r, err = http.Get(ts.URL + "/clusters/" + id)
	assert.Nil(t, err)
	assert.Equal(t, r.StatusCode, http.StatusOK)
	assert.Equal(t, r.Header.Get("Content-Type"), "application/json; charset=UTF-8")

	var msg ClusterInfoResponse
	err = GetJsonFromResponse(r, &msg)
	assert.Nil(t, err)
	assert.Equal(t, id, msg.Id)
	assert.Equal(t, 1, len(msg.Rings))

	req, err := http.NewRequest("DELETE", ts.URL+"/clusters/"+id, nil)
	assert.Nil(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusConflict)
}
