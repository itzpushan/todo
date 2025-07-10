package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/itzpushan/todo/ent"
	"github.com/itzpushan/todo/ent/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("thisIsMyJwtSecret")

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		user, err := client.User.Create().
			SetID(uuid.New()).
			SetName(req.Name).
			SetEmail(req.Email).
			SetPassword(string(hashedPassword)).
			Save(r.Context())

		if err != nil {
			http.Error(w, "Signup failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func Signin(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SigninRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := client.User.Query().
			Where(user.EmailEQ(req.Email)).
			Only(r.Context())

		if err != nil {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID.String(),
			"exp":     time.Now().Add(time.Hour * 168).Unix(),
		})

		tokenString, _ := token.SignedString(jwtSecret)

		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
	}
}
