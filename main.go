package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Entry struct {
	Name        string `json:"name"`
	Coordinates string `json:"coordinates"`
}

var (
	mu      sync.Mutex
	entries []Entry
)

func loadEntries() {
	file, err := os.Open("entries.json")
	if err != nil {
		if os.IsNotExist(err) {
			entries = []Entry{}
			return
		}
		panic(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&entries)
	if err != nil {
		panic(err)
	}
}

func saveEntries() {
	file, err := os.Create("entries.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(entries)
	if err != nil {
		panic(err)
	}
}

func main() {
	loadEntries()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/entries", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(entries)

		case http.MethodPost:
			var entry Entry
			err := json.NewDecoder(r.Body).Decode(&entry)
			if err != nil || entry.Name == "" || entry.Coordinates == "" {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			entries = append(entries, entry)
			saveEntries()
			w.WriteHeader(http.StatusCreated)

		case http.MethodDelete:
			var data struct {
				Index int `json:"index"`
			}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil || data.Index < 0 || data.Index >= len(entries) {
				http.Error(w, "Invalid index", http.StatusBadRequest)
				return
			}
			entries = append(entries[:data.Index], entries[data.Index+1:]...)
			saveEntries()
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server läuft auf http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
