package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"

	"github.com/restuwahyu13/sso-and-cloudstorage/configs"
	"github.com/restuwahyu13/sso-and-cloudstorage/packages"
	"github.com/restuwahyu13/sso-and-cloudstorage/routes"
)

func main() {
	SetupEnvironment()
	db := SetupDatabase()
	rh, nl := SetupServer()
	SetupMiddleware(rh)
	SetupRouter(rh, db)
	SetupGraceFullShutDown(rh, nl)
}

/*
=======================================
# GOLANG AUTOLOAD CONFIG SETUP
=======================================
*/

func SetupEnvironment() {
	if err := packages.ViperLoadConfig(); err != nil {
		logrus.Errorf("Env file config can't load: %v", err.Error())
	}
}

/*
=======================================
# GOLANG DATABASE CONNECTION SETUP
=======================================
*/

func SetupDatabase() *sqlx.DB {
	db, _ := configs.Connection("postgres")

	if err := db.Ping(); err != nil {
		logrus.Errorf("Database connection error: ", err.Error())
		return nil
	}

	logrus.Info("Database connection success")
	return db
}

/*
=======================================
# GOLANG HTPP SERVER SETUP
=======================================
*/

func SetupServer() (*chi.Mux, net.Listener) {
	var lc = net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			err := c.Control(func(fd uintptr) {
				if err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
					logrus.Errorf("Net control error: %v", err.Error())
				}
			})

			if err != nil {
				logrus.Errorf("Net control error: %v", err.Error())
			}

			return nil
		},
	}

	nl, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("127.0.0.1:%s", packages.GetString("GO_PORT")))
	if err != nil {
		logrus.Errorf("Net listen error: %v", err.Error())
	}

	rh := chi.NewRouter()
	logrus.Infof("Server running on port: %s", packages.GetString("GO_PORT"))

	return rh, nl
}

/*
=======================================
# GOLANG HTTP SERVER MIDLEWARE SETUP
=======================================
*/

func SetupMiddleware(rh *chi.Mux) {
	if packages.GetString("GO_ENV") != "production" {
		rh.Use(middleware.Logger)
	}

	rh.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		OptionsPassthrough: true,
		AllowCredentials:   true,
	}))
	rh.Use(middleware.Compress(gzip.BestCompression))
	rh.Use(middleware.ThrottleWithOpts(middleware.ThrottleOpts{Limit: 5, BacklogLimit: 50, BacklogTimeout: time.Duration(5 * time.Minute)}))
	rh.Use(middleware.NoCache)
	rh.Use(middleware.RealIP)
	rh.Use(middleware.RequestID)
}

/*
=======================================
# GOLANG HTTP SERVER ROUTING SETUP
=======================================
*/

func SetupRouter(rh *chi.Mux, db *sqlx.DB) {
	routes.NewPingRoute("/", db, rh).PingRoute()
	routes.NewStudentsRoute("/api/v1/students", db, rh).StudentsRoute()
}

/*
=======================================
# GOLANG GRACE FULL SHUTDOWN SETUP
=======================================
*/

func SetupGraceFullShutDown(rh *chi.Mux, nl net.Listener) {
	var g errgroup.Group
	server := http.Server{
		Addr:           fmt.Sprintf(":%s", packages.GetString("PORT")),
		ReadTimeout:    time.Duration(time.Second) * 60,
		WriteTimeout:   time.Duration(time.Second) * 30,
		IdleTimeout:    time.Duration(time.Second) * 120,
		MaxHeaderBytes: 3145728,
		Handler:        rh,
	}

	g.Go(func() error {
		if err := http.Serve(nl, server.Handler); err != nil {
			logrus.Errorf("Server listening error: %v", err.Error())
			os.Exit(1)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logrus.Errorf("Gorutine deathlock error: %v", err.Error())
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	server.RegisterOnShutdown(func() {
		osSignal := <-ch
		logrus.Infof("HTTP server on going to shutdown & signal received: %s", osSignal.String())
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("HTTP server shutdown error: %s", err.Error())
		return
	}

	logrus.Info("HTTP server shutdown success")
}
