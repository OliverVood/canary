package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OliverVood/courier/cmd/internal/log"
	"github.com/OliverVood/courier/cmd/internal/message"
	"github.com/OliverVood/courier/cmd/internal/queue"
)

func (h *Handler) store(w http.ResponseWriter, r *http.Request) {
	log.Request(r)

	from := r.FormValue("from")
	to := r.FormValue("to")
	body := r.FormValue("body")

	item := message.Format{
		From: from,
		To:   to,
		Body: body,
	}

	itemE, err := json.Marshal(item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = queue.Send(h.Chanel, h.Queue, string(itemE))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Println("Append message...")
}
