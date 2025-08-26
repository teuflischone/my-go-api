// @title My Go API
// @version 1.0
// @description REST API built with Go (Echo)
// @contact.name Mikhail Kopeikin
// @contact.url https://github.com/MikhailKopeikin/my-go-api

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MikhailKopeikin/my-go-api/utils/crypto"
	"github.com/MikhailKopeikin/my-go-api/utils/jwt"

	_ "github.com/MikhailKopeikin/my-go-api/docs"
	"github.com/MikhailKopeikin/my-go-api/utils"

	"github.com/MikhailKopeikin/my-go-api/config"
	httpDelivery "github.com/MikhailKopeikin/my-go-api/delivery/http"
	appMiddleware "github.com/MikhailKopeikin/my-go-api/delivery/middleware"
	"github.com/MikhailKopeikin/my-go-api/infrastructure/datastore"
	"github.com/MikhailKopeikin/my-go-api/internal/ops"
	pgsqlRepository "github.com/MikhailKopeikin/my-go-api/repository/pgsql"
	redisRepository "github.com/MikhailKopeikin/my-go-api/repository/redis"
	"github.com/MikhailKopeikin/my-go-api/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	configApp := config.LoadConfig()

	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	todoRepo := pgsqlRepository.NewPgsqlTodoRepository(dbInstance)
	userRepo := pgsqlRepository.NewPgsqlUserRepository(dbInstance)

	cryptoSvc := crypto.NewCryptoService()
	jwtSvc := jwt.NewJWTService(configApp.JWTSecretKey)

	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	todoUC := usecase.NewTodoUsecase(todoRepo, redisRepo, ctxTimeout)
	authUC := usecase.NewAuthUsecase(userRepo, cryptoSvc, jwtSvc, ctxTimeout)

	appMiddleware := appMiddleware.NewMiddleware(jwtSvc)

	e := echo.New()
	ops.Register(e)
	e.Use(middleware.CORS())
	e.Use(appMiddleware.Logger(nil))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewTodoHandler(e, appMiddleware, todoUC)
	httpDelivery.NewAuthHandler(e, appMiddleware, authUC)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configApp.ContextTimeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
