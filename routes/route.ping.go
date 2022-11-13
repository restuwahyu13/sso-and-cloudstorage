package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type pingRoute struct {
	prefix string
	db     *sqlx.DB
	router *chi.Mux
}

func NewPingRoute(prefix string, db *sqlx.DB, router *chi.Mux) *pingRoute {
	return &pingRoute{prefix, db, router}
}

func (r *pingRoute) PingRoute() {
	r.router.Route(r.prefix, func(route chi.Router) {
		route.Get(r.prefix, func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")

			res := Response{
				StatCode: http.StatusOK,
				StatMsg:  "Ping Server OK",
			}

			json.NewEncoder(rw).Encode(&res)
		})
	})
}
