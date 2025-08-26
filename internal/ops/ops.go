package ops

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func Register(e *echo.Echo) {
	e.Use(middleware.RequestID())

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/healthz", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 500*time.Millisecond)
		defer cancel()

		dbStatus := "unset"
		if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
			dbStatus = "down"
			if db, err := sql.Open("postgres", dsn); err == nil {
				defer db.Close()
				if err := db.PingContext(ctx); err == nil {
					dbStatus = "up"
				}
			}
		}

		redisStatus := "unset"
		if addr := os.Getenv("REDIS_ADDR"); addr != "" {
			redisStatus = "down"
			rdb := redis.NewClient(&redis.Options{Addr: addr})
			defer rdb.Close()
			if err := rdb.Ping(ctx).Err(); err == nil {
				redisStatus = "up"
			}
		}

		status := http.StatusOK
		if dbStatus == "down" || redisStatus == "down" {
			status = http.StatusServiceUnavailable
		}
		return c.JSON(status, map[string]string{
			"status": "ok",
			"db":     dbStatus,
			"redis":  redisStatus,
		})
	})
}
