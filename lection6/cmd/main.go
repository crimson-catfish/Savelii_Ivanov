package main

import (
	"net/http"
	"time"

	"entrance/lection6/internal/handlers/auth"
	"entrance/lection6/internal/handlers/chat"
	"entrance/lection6/internal/middlewares"
	"entrance/lection6/internal/storage/mock"

	"github.com/go-chi/chi/v5"
)

const port = ":8080"

func main() {
	r := chi.NewRouter()

	mockRepository := mock.NewMockRepository()

	authService := auth.NewAuthService(mockRepository)
	chatService := chat.DefaultChatService(mockRepository)

	r.Post("/signup", authService.SignUp)
	r.Post("/signin", authService.SignIn)

	r.Get("/publicChats", chatService.ListPublic)
	r.Get("/publicChats/{chatName}", chatService.ReadPublic)

	r.Group(
		func(r chi.Router) {
			r.Use(middlewares.Auth)

			r.Post("/publicChats/{chatName}", chatService.SendToPublic)

			r.Get("/myChats", chatService.ListPrivate)
			r.Get("/myChats/{chatName}", chatService.ReadPrivate)
			r.Post("/myChats/{chatName}", chatService.SendPrivate)
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
