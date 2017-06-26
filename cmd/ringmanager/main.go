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

package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/spf13/viper"
	"github.com/thiagodasilva/swift-ring-manager/pkg/ringmanager"
)

func loadConfig() (*viper.Viper, error) {
	v := viper.New()
	viper.SetConfigName("ringmanager")
	viper.AddConfigPath("/etc/ringmanager")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	loadDefaultConfigOptions(v)
	return v, nil
}

func loadDefaultConfigOptions(v *viper.Viper) {
	v.SetDefault("dbfilename", "swift_clusters.db")
	v.SetDefault("ringmanager_dir", "/var/lib/ringmanager")
	v.SetDefault("bind_ip", "127.0.0.1")
	v.SetDefault("bind_port", "8090")

}

func main() {
	v, err := loadConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error loading config file: %s", err))
	}
	addr := v.GetString("bind_ip") + ":" + v.GetString("bind_port")
	router := ringmanager.NewRouter(v)
	log.Fatal(http.ListenAndServe(addr, router))
}
