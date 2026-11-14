package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/OliverVood/courier/cmd/internal/db"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	DB *sql.DB
}

func main() {
	fmt.Println("Launching the client management program...")

	var err error

	h := &Handler{}

	h.DB, err = db.Instance()
	if err != nil {
		log.Fatal(err)
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(h.DB)

	router := chi.NewRouter()
	regRouts(router, h)
	regFileServer(router)
	runServer(router)
}

func regRouts(router chi.Router, h *Handler) {

	fmt.Println("Router registration...")

	router.Get("/", h.home)
	router.Route("/clients", func(router chi.Router) {
		router.Get("/", h.index)
		router.Get("/create", h.create)
		router.Post("/", h.store)
		router.Get("/{client}", h.show)
		router.Get("/{client}/edit", h.edit)
		router.Patch("/{client}", h.update)
		router.Delete("/{client}", h.destroy)
	})
}

func regFileServer(router chi.Router) {
	fmt.Println("File server registration...")

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		context := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(context.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir("./static")))
		fs.ServeHTTP(w, r)
	})
}

func runServer(router chi.Router) {
	fmt.Println("Server startup...")
	fmt.Println("Open site: http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
