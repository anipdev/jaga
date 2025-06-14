package routes

import (
	"jaga/controllers"
	"jaga/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", controllers.Login).Methods("POST")

	r.HandleFunc("/users", middleware.RequireRole("super_user", "admin")(controllers.CreateUser)).Methods("POST")

	r.HandleFunc("/assets",
		middleware.RequireRole("super_user", "admin")(controllers.CreateAsset),
	).Methods("POST")

	return r
}
