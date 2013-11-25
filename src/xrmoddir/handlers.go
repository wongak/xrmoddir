package xrmoddir

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
)

func setHandlers(s *Server) error {
	s.Get("/", index)
	s.Get("/about", about)
	return nil
}

// The Index Page
func index(
	w http.ResponseWriter,
	t *template.Template,
	l *log.Logger,
) {
	var buf bytes.Buffer
	c := NewContent()
	err := t.ExecuteTemplate(&buf, "index.tmpl.html", c)
	if err != nil {
		log.Printf("Template error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// About page
func about(
	t *template.Template,
	l *log.Logger,
) (respCode int, body string) {
	var buf bytes.Buffer
	respCode = http.StatusOK
	c := NewContent()
	err := t.ExecuteTemplate(&buf, "about.tmpl.html", c)
	if err != nil {
		log.Printf("Template error: %v", err)
		respCode = http.StatusInternalServerError
		return
	}
	body = buf.String()
	return
}
