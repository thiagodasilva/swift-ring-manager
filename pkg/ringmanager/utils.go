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
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/lpabon/godbc"
)

// Return a 16-byte uuid
// From http://www.ashishbanerjee.com/home/go/go-generate-uuid
func GenUUID() string {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	godbc.Check(n == len(uuid), n, len(uuid))
	godbc.Check(err == nil, err)

	return hex.EncodeToString(uuid)
}

func jsonFromBody(r io.Reader, v interface{}) error {

	// Check body
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}

// Unmarshal JSON from request
func GetJsonFromRequest(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return jsonFromBody(r.Body, v)
}

// Unmarshal JSON from response
func GetJsonFromResponse(r *http.Response, v interface{}) error {
	defer r.Body.Close()
	return jsonFromBody(r.Body, v)
}

// Check if a sorted string list has a string
func SortedStringHas(s sort.StringSlice, x string) bool {
	index := s.Search(x)
	if index == len(s) {
		return false
	}
	return s[s.Search(x)] == x
}

// Delete a string from a sorted string list
func SortedStringsDelete(s sort.StringSlice, x string) sort.StringSlice {
	index := s.Search(x)
	if len(s) != index && s[index] == x {
		s = append(s[:index], s[index+1:]...)
	}

	return s
}

// FileHash returns the sha256 hash of file
// It expects the file already opened and it will not close it.
func FileHash(file *os.File) (result string, err error) {
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return result, nil
}
