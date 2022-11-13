package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Response struct {
	StatCode uint        `json:"stat_code"`
	StatMsg  string      `json:"stat_msg"`
	Data     interface{} `json:"data,omitempty"`
}

type studentsRoute struct {
	prefix string
	db     *sqlx.DB
	router *chi.Mux
}

func NewStudentsRoute(prefix string, db *sqlx.DB, router *chi.Mux) *studentsRoute {
	return &studentsRoute{prefix, db, router}
}

func (r *studentsRoute) StudentsRoute() {
	r.router.Route(r.prefix, func(route chi.Router) {

		route.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")

			res := Response{
				StatCode: http.StatusOK,
				StatMsg:  "Building sso and cloud storage using keyloack and minio",
			}

			json.NewEncoder(rw).Encode(&res)
		})

	})
}
