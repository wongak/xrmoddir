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
	return nil
}

func index(
	w http.ResponseWriter,
	c *XRModDirContent,
	t *template.Template,
	l *log.Logger,
) {
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, "index.tmpl.html", c)
	if err != nil {
		log.Printf("Template error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}
