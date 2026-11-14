package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/OliverVood/courier/cmd/internal/log"
	"github.com/OliverVood/courier/cmd/internal/ui/template"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	templateData := template.NewData(nil, "Home")

	template.Render(w, filepath.Join("pages", "home"), templateData)
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	items, err := clients.All(r.Context(), h.DB)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}

	templateData := template.NewData(items, "ClientItems")

	template.Render(w, filepath.Join("clients", "index"), templateData)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	templateData := template.NewData(nil, "ClientItem creating")

	template.Render(w, filepath.Join("clients", "create"), templateData)
}

func (h *Handler) store(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	login := r.FormValue("login")
	email := r.FormValue("email")

	var err error

	item := ClientItem{
		Login: login,
		Email: email,
	}

	item, err = clients.Save(r.Context(), h.DB, item)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	item, err := clients.Get(r.Context(), h.DB, uuid)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}

	templateData := template.NewData(item, "User")

	template.Render(w, filepath.Join("clients", "show"), templateData)
}

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	item, err := clients.Get(r.Context(), h.DB, uuid)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}

	templateData := template.NewData(item, "ClientItem editing")

	template.Render(w, filepath.Join("clients", "edit"), templateData)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	isset, err := clients.Isset(r.Context(), h.DB, uuid)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}
	if !isset {
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}

	email := r.FormValue("email")
	login := r.FormValue("login")

	item := ClientItem{
		Login: login,
		Email: email,
	}

	item, err = clients.Update(r.Context(), h.DB, uuid, item)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) destroy(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	uuid := chi.URLParam(r, "client")

	isset, err := clients.Isset(r.Context(), h.DB, uuid)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}
	if !isset {
		http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		return
	}

	err = clients.Destroy(uuid, h.DB)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
