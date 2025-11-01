package main

import (
	"context"
	"giftlock/internal/admin"
	"giftlock/internal/assets"
	"giftlock/internal/auth"
	"giftlock/internal/claim"
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
	log.SetFlags(log.LstdFlags | log.Llongfile)
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	db := database.NewPostgresPool(ctx, cfg.DatabaseUrl)

	userRepo := user.NewPostgresRepository(db.Pool)
	sessionRepo := session.NewPostgresRepository(db.Pool)
	giftRepo := gift.NewPostgresRepository(db.Pool)
	groupRepo := group.NewPostgresRepository(db.Pool)
	claimRepo := claim.NewPostgresRepository(db.Pool)

	adminService := admin.NewService(userRepo)
	userService := user.NewService(userRepo)
	authService := auth.NewService(userRepo, sessionRepo)
	giftService := gift.NewService(giftRepo)
	groupService := group.NewService(groupRepo, giftRepo)
	claimService := claim.NewService(claimRepo)
	htmlService := presentation.NewHtmlPresentationService()

	adminHandler := admin.NewHandler(adminService, htmlService)
	userHandler := user.NewHandler(userService, htmlService)
	authHandler := auth.NewHandler(authService, htmlService)
	giftHandler := gift.NewHandler(giftService, htmlService)
	groupHandler := group.NewHandler(groupService, htmlService)
	claimHandler := claim.NewHandler(claimService, htmlService)
	assetsHandler := assets.NewHandler()
	pagesHandler := pages.NewHandler(htmlService)

	middlewareStack := middleware.LoadMiddleware(sessionRepo, cfg)
	server := server.NewServer(
		cfg.ServerAddr,
		middlewareStack,
		authHandler,
		adminHandler,
		userHandler,
		giftHandler,
		groupHandler,
		claimHandler,
		assetsHandler,
		pagesHandler,
	)
	server.Run()
}
