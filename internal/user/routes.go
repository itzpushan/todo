package user

import (
	"github.com/gorilla/mux"
	"github.com/itzpushan/todo/ent"
)

func Register(r *mux.Router, client *ent.Client) {
	r.HandleFunc("/signup", Signup(client)).Methods("POST")
	r.HandleFunc("/signin", Signin(client)).Methods("POST")
}
