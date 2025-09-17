package main

import (
	"context"
	"giftlock/internal/api"
	"giftlock/internal/auth"
	"giftlock/internal/config"
	"giftlock/internal/database"
	"giftlock/internal/middleware"
	"giftlock/internal/pages"
	"giftlock/internal/presentation"
	"giftlock/internal/session"
	"giftlock/internal/user"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbUrl := config.MustLoadEnv("DATABASE_URL")
	apiAddr := config.MustLoadEnv("API_ADDR")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	db := database.NewPostgresPool(ctx, dbUrl)

	userRepo := user.NewPostgresRepository(ctx, db.Pool)
	sessionRepo := session.NewPostgresRepository(ctx, db.Pool)

	userService := user.NewService(userRepo)
	sessionService := session.NewService(sessionRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	htmlPresenter := presentation.NewHtmlPresenter()

	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	webPageHandler := pages.NewHandler(htmlPresenter)

	middlewareStack := middleware.LoadMiddleware(sessionService)
	server := api.NewServer(
		apiAddr,
		middlewareStack,
		authHandler,
		userHandler,
		webPageHandler,
	)
	server.Run()
}
