package foodlebug

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jasmaa/foodlebug/internal/auth"
	"github.com/jasmaa/foodlebug/internal/store"
)

// Handles user profile
func handleProfile(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := auth.SessionToUser(store, r)

		// Redirect to login on session failure
		if err != nil {
			http.Redirect(w, r, "./login", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusOK)
		var t *template.Template
		t, _ = template.ParseFiles(
			"assets/templates/main.html",
			"assets/templates/profile.html",
			"assets/templates/includes.html",
		)
		t.ExecuteTemplate(w, "main", user)
	})
}

// Account creation page handler
func handleCreateAccount(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		messages := make([]string, 10)

		switch r.Method {
		case "GET":
			displayCreateAccount(w, messages)

		case "POST":
			username := r.FormValue("username")
			password := r.FormValue("password")
			confirm := r.FormValue("confirm-password")

			if username == "" || password == "" {
				messages := append(messages, "Fields cannot be empty.")
				displayCreateAccount(w, messages)
				return
			} else if password != confirm {
				messages := append(messages, "Passwords do not match.")
				displayCreateAccount(w, messages)
				return
			}

			err := auth.CreateNewUser(store, username, password)

			if err != nil {
				messages := append(messages, "Account could not be created.")
				displayCreateAccount(w, messages)
				return
			}

			// go to success page
			displayCreateAccountSuccess(w)
		}
	})
}

// Login page handler
func handleLogin(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		messages := make([]string, 10)

		switch r.Method {
		case "GET":
			displayLogin(w, messages)

		case "POST":
			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := store.GetUser("username", username)
			if err != nil {
				messages := append(messages, "User does not exist.")
				displayLogin(w, messages)
				return
			}

			if !auth.VerifyPassword(user.Password, password) {
				messages := append(messages, "Incorrect password.")
				displayLogin(w, messages)
				return
			}

			// set session cookie
			s, err := auth.NewSession(store, user, time.Time{}, "ip address", "user agent")
			if err != nil {
				messages := append(messages, "Could not login.")
				displayLogin(w, messages)
				return
			}

			cookie := &http.Cookie{
				Name:     "foodlebug",
				Value:    fmt.Sprintf("%s_%s", s.UserKey, s.SessionId),
				HttpOnly: true,
				Path:     "/",
				Expires:  s.Expires,
			}

			http.SetCookie(w, cookie)
			http.Redirect(w, r, "./", http.StatusSeeOther)
		}
	})
}

// == Static display ==

// Display create account screen get
func displayCreateAccount(w http.ResponseWriter, messages []string) {

	w.WriteHeader(http.StatusOK)
	var t *template.Template
	t, _ = template.ParseFiles(
		"assets/templates/main.html",
		"assets/templates/includes.html",
		"assets/templates/createAccount.html",
	)
	t.ExecuteTemplate(w, "main", messages)
}

// Display login screen get
func displayLogin(w http.ResponseWriter, messages []string) {

	w.WriteHeader(http.StatusOK)
	var t *template.Template
	t, _ = template.ParseFiles(
		"assets/templates/main.html",
		"assets/templates/includes.html",
		"assets/templates/login.html",
	)
	t.ExecuteTemplate(w, "main", messages)
}

// Account creation success
func displayCreateAccountSuccess(w http.ResponseWriter) {

	w.WriteHeader(http.StatusOK)
	var t *template.Template
	t, _ = template.ParseFiles(
		"assets/templates/main.html",
		"assets/templates/includes.html",
		"assets/templates/createAccountSuccess.html",
	)
	t.ExecuteTemplate(w, "main", "")
}
