package ringmanager

import (
	"net/http/httptest"
	"testing"

	"bytes"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestClusterCreate(t *testing.T) {

	// setup the server
	router := NewRouter()
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
