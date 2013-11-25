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

var buf bytes.Buffer

func index(
	w http.ResponseWriter,
	t *template.Template,
	l *log.Logger,
) {
	buf.Reset()
	c := NewContent()
	c.Content = "index.tmpl.html"
	err := t.ExecuteTemplate(&buf, "layout_main.tmpl.html", c)
	if err != nil {
		log.Printf("Template error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func about(
	t *template.Template,
	l *log.Logger,
) (respCode int, body string) {
	respCode = http.StatusOK
	buf.Reset()
	c := NewContent()
	c.Content = "about.tmpl.html"
	err := t.ExecuteTemplate(&buf, "layout_main.tmpl.html", c)
	if err != nil {
		log.Printf("Template error: %v", err)
		respCode = http.StatusInternalServerError
		return
	}
	body = buf.String()
	return
}
