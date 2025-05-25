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

func (s *Server) MapHandler(usecase *registry.UseCase, logger core.Logger) {
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(customMiddleware.Error(logger))

	s.GET("/health-check", handler.NewHealthCheck())
	s.POST("/login", handler.NewLogin(usecase.Login))
	s.POST("/signup", handler.NewSignUp(usecase.SignUp))
	s.POST("/signup/confirm", handler.NewSignUpConfirm(usecase.SignUpConfirm))
	s.POST("/logout", handler.NewLogout())

	// 認証が必要なルートグループ
	auth := s.Group("")
	auth.Use(customMiddleware.Auth(usecase.VerifyToken))
	auth.GET("/me", handler.NewGetAccount(usecase.GetAccount))
	auth.GET("/categories", handler.NewGetCategories(usecase.GetCategories))

	// 個人買い物関連のルート
	personal := auth.Group("/personal")
	personal.GET("/items", handler.NewGetPersonalItems(usecase.GetPersonalItems))
	personal.POST("/items", handler.NewCreatePersonalItem(usecase.CreatePersonalItem))
	personal.PUT("/items/:id", handler.NewUpdatePersonalItem(usecase.UpdatePersonalItem))
	personal.DELETE("/items/:id", handler.NewDeletePersonalItem(usecase.DeletePersonalItem))

	// 家族関連のルート
	family := auth.Group("/family")
	family.POST("", handler.NewCreateFamily(usecase.CreateFamily))
	family.GET("/invitation", handler.NewInvitationFamily(usecase.InvitationFamily))
	family.POST("/join/{token}", handler.NewJoinFamily(usecase.JoinFamily))
	family.POST("/items", handler.NewCreateFamilyItem(usecase.CreateFamilyItem))
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
