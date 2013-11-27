package user

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"xrmoddir/content"
)

const (
	USER_USERNAME_MINLENGTH = 4
	USER_PASSWORD_MINLENGTH = 6
)

// User handlers for http
type Handler struct {
	UsernamePattern *regexp.Regexp
}

func NewHandler() (*Handler, error) {
	pattern, err := regexp.Compile("/^[a-z][a-z0-9]{3,}$/")
	if err != nil {
		return nil, fmt.Errorf("User handler initialization error: %v", err)
	}
	h := &Handler{
		UsernamePattern: pattern,
	}
	return h, nil
}

func (h *Handler) SetRoutes(m martini.Router) {
	m.Get("/register", h.Register())
	m.Post("/register", h.HandleRegistration())
}

// Registration page
func (h *Handler) Register() func(*template.Template, *log.Logger) (int, string) {
	return func(
		t *template.Template,
		l *log.Logger,
	) (respCode int, body string) {
		var buf bytes.Buffer
		respCode = http.StatusOK
		c := content.NewContent()
		c.Data["passwordMinLength"] = USER_PASSWORD_MINLENGTH
		c.Data["usernameMinLength"] = USER_USERNAME_MINLENGTH
		err := t.ExecuteTemplate(&buf, "register.tmpl.html", c)
		if err != nil {
			log.Printf("Template error: %v", err)
			respCode = http.StatusInternalServerError
			return
		}
		body = buf.String()
		return
	}
}

func (h *Handler) HandleRegistration() func(
	w http.ResponseWriter,
	r *http.Request,
	t *template.Template,
	l *log.Logger,
	db *sql.DB,
) {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		t *template.Template,
		l *log.Logger,
		db *sql.DB,
	) {
		var buf bytes.Buffer
		c := content.NewContent()
		c.Data["passwordMinLength"] = USER_PASSWORD_MINLENGTH
		c.Data["usernameMinLength"] = USER_USERNAME_MINLENGTH

		// form handling
		var userId int64
		var err error
		u := new(User)
		errors := make(map[string]bool)
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")
		if username == "" {
			errors["noUsername"] = true
		} else {
			u.Username = username
			if len(username) < USER_USERNAME_MINLENGTH {
				errors["usernameLen"] = true
			} else {
				userId, err = SQLIdByUsername(db, username)
				if err != nil {
					log.Printf("Error on checking username: %v", err)
					errors["internal"] = true
				} else if userId != 0 {
					errors["usernameInUse"] = true
				}
			}
		}
		if email == "" {
			errors["noEmail"] = true
		} else {
			u.Email = email
			userId, err = SQLIdByEmail(db, email)
			if err != nil {
				log.Printf("Error on checking email: %v", err)
				errors["internal"] = true
			} else if userId != 0 {
				errors["emailInUse"] = true
			}
		}
		if password == "" {
			errors["noPassword"] = true
		} else if len(password) < USER_PASSWORD_MINLENGTH {
			errors["passwordLen"] = true
		} else if password != password2 {
			errors["passwordMismatch"] = true
		} else {
			err := u.SetPassword(password)
			if err != nil {
				log.Printf("Error on setting password: %v", err)
				errors["internal"] = true
			}
		}
		c.Data["errors"] = errors
		c.Data["User"] = u
		err = t.ExecuteTemplate(&buf, "register.tmpl.html", c)
		if err != nil {
			log.Printf("Template error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(errors) > 0 {
			io.Copy(w, &buf)
			return
		}
	}
}
