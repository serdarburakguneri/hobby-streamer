package http

import (
	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/asset"
)

func NewRouter(repo *asset.Repository) *mux.Router {
	r := mux.NewRouter()

	h := &Handler{Repo: repo}

	r.HandleFunc("/assets", h.ListAssets).Methods("GET")
	r.HandleFunc("/assets/{id}", h.GetAsset).Methods("GET")
	r.HandleFunc("/assets", h.SaveAsset).Methods("POST")

	return r
}