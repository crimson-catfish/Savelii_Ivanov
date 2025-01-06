package main

import (
	"net/http"
	"time"

	"entrance/lection5/handlers"
	"entrance/lection5/middlewares"

	"github.com/go-chi/chi/v5"
)

const port = ":8080"

func main() {
	r := chi.NewRouter()

	r.Post("/signup", handlers.SignUp)
	r.Post("/signin", handlers.SignIn)

	r.Get("/publicChats", handlers.ListPublicChats)
	r.Get("/publicChats/{chatName}", handlers.ReadPublicChat)

	r.Group(
		func(r chi.Router) {
			r.Use(middlewares.Auth)

			r.Post("/publicChats/{chatName}", handlers.SendToPublicChat)

			r.Get("/myChats", handlers.ListPrivateChats)
			r.Get("/myChats/{chatName}", handlers.ReadPrivateChat)
			r.Post("/myChats/{chatName}", handlers.SendPrivate)
		},
	)

	server := &http.Server{
		Addr:              port,
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
