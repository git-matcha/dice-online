package controller

//controller controls the traffic, so it takes a route, makes the appropriate repo
//the repo returns an object to the controller and the controller passes that to the browser
// in some kind of response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/comethale/dice-online/app/api/shared/repositories/usermanagement/repository"
	"github.com/comethale/dice-online/app/api/shared/session"
	"github.com/comethale/dice-online/app/api/shared/utils"

	"github.com/comethale/dice-online/app/api/shared/database"
)

// UserGet takes a GET request and returns a user from the database
func UserGet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		repo := repository.NewUserRepository(database.POSTGRESQL)

		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		id, err := strconv.Atoi(r.FormValue("id"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		user, err := repo.Get(int64(id))

		json, err := json.Marshal(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)

	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// UserCreate takes a POST request and creates a user if the required information is provided
func UserCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		repo := repository.NewUserRepository(database.POSTGRESQL)

		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.PostFormValue("email")
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		// not a switch because that would be unnecessarily long
		// not indiviual if-statements because all args are required
		if email == "" || username == "" || password == "" {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		user, err := repo.Create(email, password, username)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)

	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// UserLogin logs in the given user
func UserLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		repo := repository.NewUserRepository(database.POSTGRESQL)

		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		hashedInputPassword, err := utils.AuthHashPassword(password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		DBPassword, err := repo.GetPassword(email)

		// create session
		session := session.Instance(r, email)

		loggedIn := utils.AuthVerifyPassword(DBPassword, hashedInputPassword)
		if loggedIn != nil {
			http.Error(w, loggedIn.Error(), http.StatusUnauthorized)
			return
		}

		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Save(r, w)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged In."))

		return
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// UserLogout logs out the user and returns the index page
func UserLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		session := session.Instance(r, email)

		// Revoke users authentication
		session.Values["authenticated"] = false
		session.Save(r, w)
		return

	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}