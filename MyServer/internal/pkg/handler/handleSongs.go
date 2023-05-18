package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) getSong(w http.ResponseWriter, r *http.Request) {
	if h.doLog {
		log.Println("New HTTP connection")
	}
	vars := mux.Vars(r)

	name := vars["name"]
	b, err := h.fmanager.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("name", name)
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func (h *Handler) getAllSongs(w http.ResponseWriter, r *http.Request) {
	if h.doLog {
		log.Println("New HTTP connection")
	}
	str := h.fmanager.GetAll()

	fmt.Println(str)
	b, err := json.Marshal(&str)
	if err != nil {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h *Handler) saveSong(w http.ResponseWriter, r *http.Request) {
	if h.doLog {
		log.Println("New HTTP connection")
	}
	vars := mux.Vars(r)

	name := vars["name"]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error: Cant save file", http.StatusBadRequest)
		return
	}
	err = h.fmanager.Add(name, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteSong(w http.ResponseWriter, r *http.Request) {
	if h.doLog {
		log.Println("New HTTP connection")
	}
	vars := mux.Vars(r)
	name := vars["name"]

	err := h.fmanager.Delete(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
