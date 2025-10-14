package main

import (
	"context"
	"giftlock/internal/assets"
	"giftlock/internal/auth"
	"giftlock/internal/config"
	"giftlock/internal/database"
	"giftlock/internal/gift"
	"giftlock/internal/group"
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
	groupRepo := group.NewPostgresRepository(db.Pool)

	userService := user.NewService(userRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	giftService := gift.NewService(giftRepo)
	groupService := group.NewService(groupRepo)
	htmlService := presentation.NewHtmlPresentationService()

	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService, htmlService)
	giftHandler := gift.NewHandler(giftService)
	groupHandler := group.NewHandler(groupService, htmlService)
	assetsHandler := assets.NewHandler()
	pagesHandler := pages.NewHandler(htmlService)

	middlewareStack := middleware.LoadMiddleware(sessionRepo, cfg)
	server := server.NewServer(
		cfg.ServerAddr,
		middlewareStack,
		authHandler,
		userHandler,
		giftHandler,
		groupHandler,
		assetsHandler,
		pagesHandler,
	)
	server.Run()
}
