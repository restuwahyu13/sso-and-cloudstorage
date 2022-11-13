package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	config "github.com/restuwahyu13/sso-and-cloudstorage/configs"
	pkg "github.com/restuwahyu13/sso-and-cloudstorage/packages"
	"github.com/restuwahyu13/sso-and-cloudstorage/routes"
)

func main() {
	SetupEnvironment()
	db := SetupDatabase()
	rh, nl := SetupServer()
	SetupRouter(rh, db)
	SetupGraceFullShutDown(rh, nl)
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
# GOLANG AUTOLOAD CONFIG SETUP
=======================================
*/

func SetupEnvironment() {
	if err := pkg.ViperLoadConfig(); err != nil {
		logrus.Errorf("Env file config can't load: %v", err.Error())
	}
}

/*
=======================================
# GOLANG DATABASE CONNECTION SETUP
=======================================
*/

func SetupDatabase() *sqlx.DB {
	db, _ := config.Connection("postgres")

	if err := db.Ping(); err != nil {
		logrus.Errorf("Database connection error: ", err.Error())
		return nil
	}

	logrus.Info("Database connection success")
	return db
}

/*
=======================================
# GOLANG HTPP SERVER  SETUP
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

	nl, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("127.0.0.1:%s", pkg.GetString("GO_PORT")))
	if err != nil {
		logrus.Errorf("Net listen error: %v", err.Error())
	}

	rh := chi.NewRouter()

	if pkg.GetString("GO_ENV") == "development" {
		logrus.Infof("Server running on port: %s", pkg.GetString("GO_PORT"))
	}

	return rh, nl
}

/*
=======================================
# GOLANG GRACE FULL SHUTDOWN SETUP
=======================================
*/

func SetupGraceFullShutDown(rh *chi.Mux, nl net.Listener) {
	server := http.Server{
		Addr:           fmt.Sprintf(":%s", pkg.GetString("PORT")),
		ReadTimeout:    time.Duration(time.Second) * 60,
		WriteTimeout:   time.Duration(time.Second) * 30,
		IdleTimeout:    time.Duration(time.Second) * 120,
		MaxHeaderBytes: 3145728,
		Handler:        rh,
	}

	if err := http.Serve(nl, server.Handler); err != nil {
		logrus.Errorf("Server listening error: %v", err.Error())
		os.Exit(1)
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
	}

	logrus.Info("HTTP server shutdown success")
}
