package main

import (
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/julienschmidt/httprouter"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()


r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)


r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

// healthcheck
r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

menu1 := r.PathPrefix("/api/v1").Subrouter()


menu1.HandleFunc("/tasks", app.listTasksHandler).Methods("GET")
// Create a new menu
menu1.HandleFunc("/task", app.createTaskHandler).Methods("POST")
// Get a specific menu
menu1.HandleFunc("/tasks/{id:[0-9]+}", app.showTaskHandler).Methods("GET")
// Update a specific menu
menu1.HandleFunc("/tasks/{id:[0-9]+}", app.updateTaskHandler).Methods("PUT")
// Delete a specific menu
menu1.HandleFunc("/tasks/{id:[0-9]+}", app.deleteTaskHandler).Methods("DELETE")
users1 := r.PathPrefix("/api/v1").Subrouter()
// User handlers with Authentication
users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")
// Add the enableCORS() middleware.
return app.authenticate(r)


	// // Add the route for the POST /v1/users endpoint.
	// router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// // Add the route for the PUT /v1/users/activated endpoint.  activateUserHandler
	// router.HandlerFunc(http.MethodPut, "/v1/activate", app.activateUserHandler)
	// // Add the route for the POST /v1/tokens/authentication endpoint.
	// router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

}
