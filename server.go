package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bstick12/goflake"
	"github.com/gorilla/mux"
)

var generator = goflake.GoFlakeInstanceUsingUnique("D01Z01")
var started = time.Now()

func main() {
	startServer()
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ids", Id)
	router.Queries("count", "{count:[0-9]+}")
	router.HandleFunc("/healthz", Health)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Id(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	countVar := values["count"]
	var count = 2
	if len(countVar) == 1 {
		count, _ = strconv.Atoi(countVar[0])
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	ids := []string{}
	for i := 0; i < count; i++ {
		ids = append(ids, generator.GetBase64UUID())
	}

	json.NewEncoder(w).Encode(ids)
}

func Health(w http.ResponseWriter, r *http.Request) {
	duration := time.Since(started)

	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
