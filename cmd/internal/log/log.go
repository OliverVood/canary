package log

import (
	"fmt"
	"net/http"
)

func Request(r *http.Request) {
	fmt.Printf("REQUEST [METHOD: %s, URL: %s]\n", r.Method, r.URL.Path)
}
