package handler

import (
	"github.com/gorilla/mux"
)

type Mp3FileManager interface {
	Add(name string, input []byte) error
	Get(name string) ([]byte, error)
	GetAll() []string
	Delete(name string) error
}

type Handler struct {
	fmanager Mp3FileManager
}

func NewHandler(fmanager Mp3FileManager) *Handler {
	return &Handler{fmanager: fmanager}
}

func (h *Handler) InitRouter() *mux.Router {

	router := mux.NewRouter()


	router.HandleFunc("/{name}", h.getSong).Methods("GET")
	router.HandleFunc("/", h.getAllSongs).Methods("GET")
	router.HandleFunc("/{name}", h.saveSong).Methods("POST")
	router.HandleFunc("/{name}", h.deleteSong).Methods("DELETE")

	return router
}