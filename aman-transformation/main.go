/*
Copyright 2022 TriggerMesh Inc.

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
	"encoding/json"
	"fmt"

	// "fmt"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type namecard struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	// sugar.Infow("failed to fetch URL")

	sugar.Info("TRANSFORMATION APP")
	http.HandleFunc("/", home)
	http.HandleFunc("/index", index)
	http.HandleFunc("/bobtom", bobTom)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	sugar.Infof("server started at : %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		// log.Fatal(err)
		sugar.Fatal(err)

	}

}

// index will reverse the name
func index(w http.ResponseWriter, r *http.Request) {
	name := &namecard{}
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		log.Fatal(err)
	}
	str := name.Name
	strrev := reverse(str)
	log.Println("After transformation ", strrev)

	json.NewEncoder(w).Encode(&name)
}

// bobTom transforms bob to tom
func bobTom(w http.ResponseWriter, r *http.Request) {
	name := &namecard{}
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		log.Fatal(err)
	}

	bob := name.Name
	if bob == "bob" {
		name.Name = "tom"
	} else {
		log.Fatal("Name value is not BOB : Can't transform")
	}
	json.NewEncoder(w).Encode(&name)

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Menue\n1.Goto '/index' for reverse action\n2.Goto '/bobtom' for transformation option%v", r.URL.Path[:1])
}

// reverse accepts a single string input and returns the reverse of the input.
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
