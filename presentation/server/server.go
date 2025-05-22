package server

import (
	"context"
	"fmt"
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

func New(addr uint) *Server {
	addrStr := fmt.Sprintf(":%v", addr)

	return &Server{
		addr: addrStr,
		Echo: echo.New(),
	}
}

func (s *Server) MapHandler(usecase registry.UseCase, logger core.Logger) {
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())

	s.GET("/health-check", handler.NewHealthCheck())
	s.POST("/login", handler.NewLogin(usecase.NewLogin(), logger))
	s.POST("/signup", handler.NewSignUp(usecase.NewSignUp(), logger))
	s.POST("/signup/confirm", handler.NewSignUpConfirm(usecase.NewSignUpConfirm(), logger))
	s.POST("/logout", handler.NewLogout())

	// 認証が必要なルートグループ
	auth := s.Group("")
	auth.Use(customMiddleware.Auth(usecase.NewVerifyToken(), logger))
	auth.GET("/me", handler.NewGetAccount(usecase.NewGetAccount(), logger))
	auth.GET("/categories", handler.NewGetCategories(usecase.NewGetCategories()))

	// パーソナルアイテム関連のルート
	personal := auth.Group("/personal")
	personal.GET("/items", handler.NewGetPersonalItems(usecase.NewGetPersonalItems(), logger))
	personal.POST("/items", handler.NewCreatePersonalItem(usecase.NewCreatePersonalItem(), logger))
	personal.PUT("/items/:id", handler.NewUpdatePersonalItem(usecase.NewUpdatePersonalItem(), logger))
	personal.DELETE("/items/:id", handler.NewDeletePersonalItem(usecase.NewDeletePersonalItem(), logger))

	// ファミリー関連のルート
	family := auth.Group("/family")
	family.POST("", handler.NewCreateFamily(usecase.NewCreateFamily(), logger))
	family.GET("/invitation", handler.NewInvitationFamily(usecase.NewInvitationFamily(), logger))
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
