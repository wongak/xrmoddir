package xrmoddir

import (
	"bytes"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"xrmoddir/user"
)

const (
	USER_USERNAME_MINLENGTH = 4
	USER_PASSWORD_MINLENGTH = 6
)

func setHandlers(s *Server) error {
	s.Get("/", index)
	s.Get("/about", about)
	s.Get("/register", register)
	s.Post("/register", handleRegistration)
	return nil
}

func handlePage(templateName string, t *template.Template, l *log.Logger) (respCode int, body string) {
	var buf bytes.Buffer
	respCode = http.StatusOK
	c := NewContent()
	err := t.ExecuteTemplate(&buf, templateName, c)
	if err != nil {
		log.Printf("Template error: %v", err)
		respCode = http.StatusInternalServerError
		return
	}
	body = buf.String()
	return
}

// The Index Page
func index(
	t *template.Template,
	l *log.Logger,
) (int, string) {
	return handlePage("index.tmpl.html", t, l)
}

// About page
func about(
	t *template.Template,
	l *log.Logger,
) (respCode int, body string) {
	return handlePage("about.tmpl.html", t, l)
}

// Registration page
func register(
	t *template.Template,
	l *log.Logger,
) (respCode int, body string) {
	var buf bytes.Buffer
	respCode = http.StatusOK
	c := NewContent()
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

func handleRegistration(
	w http.ResponseWriter,
	r *http.Request,
	t *template.Template,
	l *log.Logger,
	db *sql.DB,
) {
	var buf bytes.Buffer
	c := NewContent()
	c.Data["passwordMinLength"] = USER_PASSWORD_MINLENGTH
	c.Data["usernameMinLength"] = USER_USERNAME_MINLENGTH

	// form handling
	var userId int64
	var err error
	u := new(user.User)
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
			userId, err = user.SQLIdByUsername(db, username)
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
		userId, err = user.SQLIdByEmail(db, email)
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
