package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	baseURL   = "https://api.unsplash.com"
	accessKey = "YOUR ACCESS_KEY"
)

type photoResponse struct {
	Results []photo `json:"results"`
}

type photo struct {
	ID     string   `json:"id"`
	Images imageSet `json:"urls"`
}

type imageSet struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

func sendRequest(w http.ResponseWriter, param, key string) (photoResponse, error) {
	url := baseURL + param + key + "&client_id=" + accessKey

	response, err := http.Get(url)

	if err != nil {
		return photoResponse{}, err
	}

	defer response.Body.Close()

	var photo photoResponse

	if err := json.NewDecoder(response.Body).Decode(&photo); err != nil {
		return photoResponse{}, err
	}

	json.NewEncoder(w).Encode(photo)

	return photo, nil
}

func showFeatured(w http.ResponseWriter, req *http.Request) {
	r, err := sendRequest(w, "/photos/featured", "")

	if err != nil {
		println(err)
	}

	jsonObj, err := json.Marshal(r)

	if err != nil {
		println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObj)
}

func search(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	r, err := sendRequest(w, "/search/photos/?query=", params["query"])

	if err != nil {
		println(err)
	}

	fmt.Printf("%+v\n", r)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/photos/{query}", search).Methods("GET")
	r.HandleFunc("/api/photos/featured", showFeatured).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
