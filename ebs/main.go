/**
 *    Copyright 2018 Amazon.com, Inc. or its affiliates
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	configurationClient *ConfigurationServiceClient
	factsClient         *FactsServiceClient
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	configurationClient = &ConfigurationServiceClient{
		client:   &http.Client{},
		clientID: os.Getenv("EXT_CLIENT_ID"),
	}
	factsClient = &FactsServiceClient{}

	// Set a default animal fact for all channels
	configurationClient.SetGlobalSegment(factsClient.GetDefaultFact())

	var r = mux.NewRouter()
	r.HandleFunc("/api/randomfact", randomFactHandler).Methods("GET")
	r.Use(VerifyJWT)

	fmt.Println("Starting server on https://localhost:8081/")
	log.Fatal(http.ListenAndServeTLS(":8081", "../conf/server.crt", "../conf/server.key", handlers.CORS(handlers.AllowedHeaders([]string{"Authorization"}))(r)))
}

// Handle the broadcaster requesting a new animal fact for their channel.
func randomFactHandler(w http.ResponseWriter, r *http.Request) {
	if channelID := r.Context().Value(ChannelIDKey); channelID != nil {
		animalType := configurationClient.GetBroadcasterSegment(channelID.(string))
		animalFact := factsClient.GetRandomFact(AnimalType(animalType))

		configurationClient.SetDeveloperSegment(channelID.(string), animalFact)

		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Channel ID is missing in the request context")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
