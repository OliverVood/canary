package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/OliverVood/courier/cmd/internal/db"
	"github.com/OliverVood/courier/cmd/internal/log"
	"github.com/OliverVood/courier/cmd/internal/ui/template"
	"github.com/go-chi/chi/v5"
)

var clients = NewTableClient()

func home(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	templateData := template.NewData(nil, "Home")

	template.Render(w, filepath.Join("pages", "home"), templateData)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	base, err := db.Instance()
	if err != nil {
		panic(err)
	}
	defer db.Close(base)

	items, err := clients.All(base)
	if err != nil {
		panic(err)
	}

	templateData := template.NewData(items, "Clients")

	template.Render(w, filepath.Join("clients", "index"), templateData)
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	templateData := template.NewData(nil, "Client creating")

	template.Render(w, filepath.Join("clients", "create"), templateData)
}

func store(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	login := r.FormValue("login")
	email := r.FormValue("email")

	item := Client{
		Login: login,
		Email: email,
	}

	base, err := db.Instance()
	if err != nil {
		panic(err)
	}
	defer db.Close(base)

	item, err = clients.Save(item, base)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func show(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	base, err := db.Instance()
	if err != nil {
		panic(err)
	}
	defer db.Close(base)

	item, ok := clients.Get(uuid, base)
	if ok != nil {
		http.NotFound(w, r)
		return
	}

	templateData := template.NewData(item, "User")

	template.Render(w, filepath.Join("clients", "show"), templateData)
}

func edit(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	base, err := db.Instance()
	if err != nil {
		panic(err)
	}
	defer db.Close(base)

	item, ok := clients.Get(uuid, base)
	if ok != nil {
		http.NotFound(w, r)
		return
	}

	templateData := template.NewData(item, "Client editing")

	template.Render(w, filepath.Join("clients", "edit"), templateData)
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	email := r.FormValue("email")
	login := r.FormValue("login")

	item := Client{
		Login: login,
		Email: email,
	}

	base, err := db.Instance()
	if err != nil {
		panic(err)
	}
	defer db.Close(base)

	item, err = clients.Update(uuid, item, base)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func destroy(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	base, err := db.Instance()
	if err != nil {
		panic(err)
	}
	defer db.Close(base)

	err = clients.Destroy(uuid, base)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
