package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/OliverVood/courier/cmd/internal/log"
	"github.com/OliverVood/courier/cmd/internal/ui/template"
	"github.com/go-chi/chi/v5"
)

func home(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	data := template.NewData(nil, "Home")

	template.Render(w, filepath.Join("pages", "home"), data)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	data := template.NewData(clients.All(), "Clients")

	template.Render(w, filepath.Join("clients", "index"), data)
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	data := template.NewData(nil, "Client creating")

	template.Render(w, filepath.Join("clients", "create"), data)
}

func store(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	email := r.FormValue("email")
	login := r.FormValue("login")

	userData := RowSchemeClient{
		Login: login,
		Email: email,
	}

	data, err := clients.Save(userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Println("id: ", data.ID, "uid: ", data.UID)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func show(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uid := chi.URLParam(r, "client")

	user, ok := clients.Get(uid)
	if ok != nil {
		http.NotFound(w, r)
		return
	}

	data := template.NewData(user, "User")

	template.Render(w, filepath.Join("clients", "show"), data)
}

func edit(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uid := chi.URLParam(r, "client")

	user, ok := clients.Get(uid)
	if ok != nil {
		http.NotFound(w, r)
		return
	}

	data := template.NewData(user, "Client editing")

	template.Render(w, filepath.Join("clients", "edit"), data)
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uid := chi.URLParam(r, "client")

	email := r.FormValue("email")
	login := r.FormValue("login")

	userData := RowSchemeClient{
		Login: login,
		Email: email,
	}

	data, err := clients.Update(uid, userData)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Println("id: ", data.ID, "uid: ", data.UID)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func destroy(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uid := chi.URLParam(r, "client")

	err := clients.Destroy(uid)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
