package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CacheGet struct {
	Val int64 `json: "val"`
	Ok  bool  `json: "ok"`
}

type CacheSet struct {
	Evicted bool `json: "evicted"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to homepage!")
}

func addToCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	fmt.Println("adding to cache ..." + key)
	val, err := strconv.ParseInt(vars["val"], 10, 64)
	if err != nil {
		panic(err)
	}
	evicted := lruSet(key, val)
	cacheSet := CacheSet{
		Evicted: evicted,
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(cacheSet)
	json.NewEncoder(w).Encode(cacheSet)
}

func getFromCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val, ok := lruGet(key)
	cacheGet := CacheGet{
		Val: val,
		Ok:  ok,
	}
	json.NewEncoder(w).Encode(cacheGet)
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func handleRequests() {
	myRouter := mux.NewRouter()
	myRouter.Use(corsMiddleware)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/set/{key}/{val}/", addToCache)
	myRouter.HandleFunc("/get/{key}/", getFromCache)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
