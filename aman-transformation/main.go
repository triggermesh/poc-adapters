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
	http.HandleFunc("/index", Index)
	http.HandleFunc("/bobtom", BobTom)
	port := os.Getenv("PORT")
	sugar.Infof("server started at : %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		// log.Fatal(err)
		sugar.Fatal(err)

	}

}

//this function will reverse the name

func Index(w http.ResponseWriter, r *http.Request) {
	name := &namecard{}
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		log.Fatal(err)
	}
	str := name.Name
	strrev := Reverse(str)
	log.Println("After transformation ", strrev)

	json.NewEncoder(w).Encode(&name)
}

// for transforming bob to tom
func BobTom(w http.ResponseWriter, r *http.Request) {
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

// Reverse accepts a single string input and returns the reverse of the input.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
