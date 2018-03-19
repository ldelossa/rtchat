package userservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	js "github.com/ldelossa/rtchat/jsonerror"
	mw "github.com/ldelossa/rtchat/middleware"
	"golang.org/x/crypto/bcrypt"
)

type HTTPServer struct {
	*http.Server
	DS DataStore
}

func NewHTTPServer(Addr string, DS DataStore) (*HTTPServer, error) {
	// Confirm Addr string is not empty
	if Addr == "" {
		return nil, fmt.Errorf("Addr string supplied is empty")
	}

	// Create server
	s := &HTTPServer{
		Server: &http.Server{
			Addr: Addr,
		},
		DS: DS,
	}

	// Create our mux
	m := http.NewServeMux()

	// Create our routes
	m.HandleFunc("/user", mw.JWTCheck(s.HandleUserCRUD))
	m.HandleFunc("/user/auth", s.HandleAuth)

	// Attach mux to server
	s.Handler = m

	return s, nil
}

func (s *HTTPServer) HandleAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Deserialize post body into auth struct
		auth := &struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := json.NewDecoder(r.Body).Decode(auth)
		if err != nil {
			log.Printf("issue deserializing json to struct: %s", err)
			js.Error(w,
				&js.Response{
					Message: "could not understand your request please check json schema",
				},
				http.StatusBadRequest)
			return
		}

		// Validate auth post
		if (auth.Username == "") || (auth.Password == "") {
			log.Printf("Received auth request with no username or password")
			js.Error(w,
				&js.Response{
					Message: "json did container username or password",
				},
				http.StatusBadRequest)
			return
		}

		// Grab user from DS
		uu, err := s.DS.GetUserByUserName(auth.Username)
		if err != nil {
			log.Printf("username not found in db: %s", err)
			js.Error(w,
				&js.Response{
					Message: "username not found",
				},
				http.StatusBadRequest)
			return
		}

		// Validate password with bycrypt
		err = bcrypt.CompareHashAndPassword([]byte(uu.Password), []byte(auth.Password))
		if err != nil {
			// Authentication failure
			log.Printf("authentication failure with username: %s password: %s", auth.Username, auth.Password)
			js.Error(w,
				&js.Response{
					Message: "authentication failure",
				},
				http.StatusUnauthorized)
			return
		}

		// Create jwt
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": auth.Username,
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString([]byte(mw.SecretKey))
		if err != nil {
			log.Printf("could not create jwt for authenticated user: %s", err)
			js.Error(w,
				&js.Response{Message: "creation of jwt failed, please attempt authentication again"},
				http.StatusInternalServerError)
			return
		}

		// Write token to response
		err = json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{Token: tokenString})
		if err != nil {
			log.Printf("could not serialize jwt: %s", err)
			js.Error(w,
				&js.Response{Message: "could not return token please attempt authentication again"},
				http.StatusInternalServerError)
			return
		}
		return

	case "DEBUG":
		// Create jwt
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": "debug",
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString([]byte(mw.SecretKey))
		if err != nil {
			log.Printf("could not create jwt for authenticated user: %s", err)
			js.Error(w,
				&js.Response{Message: "creation of jwt failed, please attempt authentication again"},
				http.StatusInternalServerError)
			return
		}

		// Write token to response
		err = json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{Token: tokenString})
		if err != nil {
			log.Printf("could not serialize jwt: %s", err)
			js.Error(w,
				&js.Response{Message: "could not return token please attempt authentication again"},
				http.StatusInternalServerError)
			return
		}
		return

	default:
		log.Printf("non supported method request to auth")
		js.Error(w,
			&js.Response{Message: "method not suppored"},
			http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) HandleUserCRUD(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Look for user query param
		r.ParseForm()
		UID := r.FormValue("user")
		if UID == "" {
			log.Printf("GET request made without user parameter")
			js.Error(w,
				&js.Response{
					Message: "request made without user query param",
				},
				http.StatusBadRequest)
			return
		}

		// Attempt to get user from DS
		u, err := s.DS.GetUserByID(UID)
		if err != nil {
			log.Printf("user query failed: %s", err)
			js.Error(w,
				&js.Response{
					Message: "user not foud",
				},
				http.StatusBadRequest)
			return
		}

		// Return user back to caller
		err = json.NewEncoder(w).Encode(&u)
		if err != nil {
			log.Printf("could not serialize response: %s", err)
			js.Error(w,
				&js.Response{
					Message: "issue with response",
				},
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		log.Printf("successfully retrieved user %s", u.ID)
		return

	case "POST":
		// Deserialize body into user struct
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("issue deserializing json to struct: %s", err)
			js.Error(w,
				&js.Response{
					Message: "could not understand your request please check json schema",
				},
				http.StatusBadRequest)
			return
		}

		// Confirm manditory fields are present
		if (u.Email == "") || (u.Username == "") || (u.Password == "") {
			log.Printf("received POST request without manditory field")
			js.Error(w,
				&js.Response{
					Message: "manditory field email, username, or password is missing",
				},
				http.StatusInternalServerError)
			return

		}

		// Create UUID for user
		u.ID = uuid.New().String()

		// Bcrypt provided password before storing
		bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
		if err != nil {
			log.Printf("error creating password hash for user: %s", err)
			js.Error(w,
				&js.Response{
					Message: "issue with persistenting user",
				},
				http.StatusInternalServerError)
			return
		}
		u.Password = string(bytes)

		// Attempt to add user to DS
		err = s.DS.AddUser(u)
		if err != nil {
			log.Printf("could not add user to datastore: %s", err)
			js.Error(w,
				&js.Response{
					Message: fmt.Sprintf("issue with adding user: %s", err),
				},
				http.StatusInternalServerError)
			return
		}

		// Return user ID and log
		json.NewEncoder(w).Encode(&struct {
			ID string `json:"id"`
		}{u.ID})
		w.Header().Set("Content-Type", "application/json")
		log.Printf("successfully added user %s", u.ID)
		return

	case "PUT", "PATCH":
		// Deserialize body into user struct
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("issue deserializing json to struct: %s", err)
			js.Error(w,
				&js.Response{
					Message: "could not understand your request please check json schema",
				},
				http.StatusBadRequest)
			return
		}

		// Confirm manditory fields are present
		if (u.Email == "") || (u.Username == "") || (u.ID == "") {
			log.Printf("received POST request without manditory field")
			js.Error(w,
				&js.Response{
					Message: "manditory field email, username, id, or password is missing",
				},
				http.StatusInternalServerError)
			return

		}

		// Validate user to be updated
		uu, err := s.DS.GetUserByID(u.ID)
		if err != nil {
			log.Printf("issue with getting to be updated user: %s", err)
			js.Error(w,
				&js.Response{
					Message: "issue with validating user",
				},
				http.StatusBadRequest)
			return
		}

		// Handle password update. If the json message provides a password field we consider
		// whether to update the user's password.
		switch {
		case u.Password == "":
			// Password not provided to update. Inject current to not overwrite value
			u.Password = uu.Password
		case u.Password != uu.Password:
			// updated password looks different then stored hash - assuming
			// we are updating the user password. Create bcrypt hash
			bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
			if err != nil {
				log.Printf("error creating password hash for user: %s", err)
				js.Error(w,
					&js.Response{
						Message: "issue with persistenting user",
					},
					http.StatusInternalServerError)
				return
			}
			u.Password = string(bytes)
		}

		// Attempt to update user
		err = s.DS.UpdateUser(u)
		if err != nil {
			log.Printf("issue with updating user: %s", err)
			js.Error(w,
				&js.Response{
					Message: "issue persistening updated user",
				},
				http.StatusInternalServerError)
			return
		}

		// Return 200 OK
		w.WriteHeader(http.StatusOK)
		log.Printf("successfully updated user %s", u.ID)
		return

	case "DELETE":
		// Look for user query param
		r.ParseForm()
		UID := r.FormValue("user")
		if UID == "" {
			log.Printf("GET request made without user parameter")
			js.Error(w,
				&js.Response{
					Message: "request made without user query param",
				},
				http.StatusBadRequest)
			return
		}

		// Attempt to delete UserID
		err := s.DS.DeleteUserByID(UID)
		if err != nil {
			log.Printf("could not delete user %s: %s", UID, err)
			js.Error(w,
				&js.Response{
					Message: "user could not be deleted",
				},
				http.StatusBadRequest)
			return
		}
		// Return 200 OK
		w.WriteHeader(http.StatusOK)
		log.Printf("successfully deleted user %s", UID)
		return

	default:
		log.Printf("non supported method request to auth")
		js.Error(w,
			&js.Response{Message: "method not suppored"},
			http.StatusInternalServerError)
		return

	}
}
