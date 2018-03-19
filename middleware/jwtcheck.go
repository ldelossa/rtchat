package middleware

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ldelossa/rtchat/jsonerror"
)

var SecretKey = "#UrMQ585~usd{NwS"

// jwtCheck is middleware which checks JWT validity before calling the next http handler
func JWTCheck(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse forms for token value
		r.ParseForm()
		t := r.FormValue("token")
		if t == "" {
			log.Printf("request made without JWT token query parameter")
			jsonerror.Error(w,
				&jsonerror.Response{Message: "please provide token query parameter with jwt value"},
				http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		if err == nil && token.Valid {
			h.ServeHTTP(w, r)
		} else {
			log.Printf("JWT not validated: %s", err)
			jsonerror.Error(w,
				&jsonerror.Response{Message: "could not validate jwt"},
				http.StatusUnauthorized)
			return
		}

	}
}
