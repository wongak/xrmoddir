// The main X Rebirth Mod Directory server
package xrmoddir

import (
	"database/sql"
	"github.com/codegangsta/martini"
	"net/http"
)

type Server struct {
	*martini.Martini
	martini.Router
}

func NewServer(db *sql.DB) (*Server, error) {
	// martini initialization
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Action(r.Handle)

	// default mapping
	m.Map(db)

	return &Server{m, r}, nil
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}
