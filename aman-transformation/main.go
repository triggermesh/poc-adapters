package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type namecard struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func main() {

	fmt.Println("******************TRANSFORMATION APP*********************")

	http.HandleFunc("/index", Index)

	fmt.Println("Sever started at port 3000\nselect localhost:3000/index")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}

}

func Index(w http.ResponseWriter, r *http.Request) {

	name := &namecard{}

	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		log.Fatal(err)
	}
	str := name.Name
	strrev := Reverse(str)
	fmt.Println("After transformation ", strrev)

	json.NewEncoder(w).Encode(&name)

}

// for reversing string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
