package main

import (
	"context"
	"giftlock/internal/auth"
	"giftlock/internal/config"
	"giftlock/internal/database"
	"giftlock/internal/gift"
	"giftlock/internal/middleware"
	"giftlock/internal/pages"
	"giftlock/internal/presentation"
	"giftlock/internal/server"
	"giftlock/internal/session"
	"giftlock/internal/user"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	db := database.NewPostgresPool(ctx, cfg.DatabaseUrl)

	userRepo := user.NewPostgresRepository(db.Pool)
	sessionRepo := session.NewPostgresRepository(db.Pool)
	giftRepo := gift.NewPostgresRepository(db.Pool)

	userService := user.NewService(userRepo)
	sessionService := session.NewService(sessionRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	giftService := gift.NewService(giftRepo)
	htmlPresenter := presentation.NewHtmlPresenter()

	userHandler := user.NewHandler(userService, htmlPresenter)
	authHandler := auth.NewHandler(authService, htmlPresenter)
	giftHandler := gift.NewHandler(giftService, htmlPresenter)
	webPageHandler := pages.NewHandler(htmlPresenter)

	middlewareStack := middleware.LoadMiddleware(sessionService)
	server := server.NewServer(
		cfg.ServerAddr,
		middlewareStack,
		authHandler,
		userHandler,
		giftHandler,
		webPageHandler,
	)
	server.Run()
}
