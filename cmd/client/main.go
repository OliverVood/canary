package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

var clients = NewTableClient()

func main() {
	fmt.Println("Launching the client management program...")

	initClients()

	router := chi.NewRouter()

	regRouts(router)
	regFileServer(router)
	runServer(router)
}

func initClients() {
	_, _ = clients.Save(RowSchemeClient{Login: "user_1", Email: "user1@example.com"})
	_, _ = clients.Save(RowSchemeClient{Login: "user_2", Email: "user2@example.com"})
	_, _ = clients.Save(RowSchemeClient{Login: "user_3", Email: "user3@example.com"})
	_, _ = clients.Save(RowSchemeClient{Login: "user_4", Email: "user4@example.com"})
	_, _ = clients.Save(RowSchemeClient{Login: "user_5", Email: "user5@example.com"})
}

func regRouts(router chi.Router) {
	fmt.Println("Router registration...")

	router.Get("/", home)
	router.Route("/clients", func(router chi.Router) {
		router.Get("/", index)
		router.Get("/create", create)
		router.Post("/", store)
		router.Get("/{client}", show)
		router.Get("/{client}/edit", edit)
		router.Patch("/{client}", update)
		router.Delete("/{client}", destroy)
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
