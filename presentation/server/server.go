package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sharebasket/core"
	"sharebasket/presentation/handler"
	customMiddleware "sharebasket/presentation/middleware"
	"sharebasket/registry"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	addr string
	*echo.Echo
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
		Echo: echo.New(),
	}
}

func (s *Server) MapHandler(usecase registry.UseCase, logger core.Logger) {
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(customMiddleware.Error(logger))

	s.GET("/health-check", handler.NewHealthCheck())
	s.POST("/login", handler.NewLogin(usecase.NewLogin()))
	s.POST("/signup", handler.NewSignUp(usecase.NewSignUp()))
	s.POST("/signup/confirm", handler.NewSignUpConfirm(usecase.NewSignUpConfirm()))
	s.POST("/logout", handler.NewLogout())

	// 認証が必要なルートグループ
	auth := s.Group("")
	auth.Use(customMiddleware.Auth(usecase.NewVerifyToken()))
	auth.GET("/me", handler.NewGetAccount(usecase.NewGetAccount()))
	auth.GET("/categories", handler.NewGetCategories(usecase.NewGetCategories()))

	// パーソナルアイテム関連のルート
	personal := auth.Group("/personal")
	personal.GET("/items", handler.NewGetPersonalItems(usecase.NewGetPersonalItems()))
	personal.POST("/items", handler.NewCreatePersonalItem(usecase.NewCreatePersonalItem()))
	personal.PUT("/items/:id", handler.NewUpdatePersonalItem(usecase.NewUpdatePersonalItem()))
	personal.DELETE("/items/:id", handler.NewDeletePersonalItem(usecase.NewDeletePersonalItem()))

	// ファミリー関連のルート
	family := auth.Group("/family")
	family.POST("", handler.NewCreateFamily(usecase.NewCreateFamily()))
	family.GET("/invitation", handler.NewInvitationFamily(usecase.NewInvitationFamily()))
	family.POST("/join/{token}", handler.NewJoinFamily(usecase.NewJoinFamily()))
	family.POST("/items", handler.NewCreateFamilyItem(usecase.NewCreateFamilyItem()))
}

func (s *Server) Run() {
	go func() {
		if err := s.Start(s.addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("HTTP server shutdown completed")
}
