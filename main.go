package main

import (
	"QR_CODE_GO/backend_modules"
	"context"
	"flag"
	"log"
	"net/http"
)

var globalSessions *backend_modules.Manager

func main() {
	globalSessions, _ = backend_modules.NewManager("memory", "gosessionid", 3600)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", backend_modules.Echo)

	http.HandleFunc("/ac", func(w http.ResponseWriter, r *http.Request) {
		backend_modules.Home(w, r, globalSessions)
	})
	http.HandleFunc("/ac/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "globalSessions", globalSessions)
		backend_modules.QrPostHandler(w, r.WithContext(ctx))
	})
	log.Fatal(http.ListenAndServe(*backend_modules.Addr, nil))
}
