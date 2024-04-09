package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the http.HandlerFunc() adapter,
	// and then set it as the custom error handler for 404 Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler
	// and set it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Use the requirePermission() middleware on each of the /v1/tasks** endpoints,
	// passing in the required permission code as the first parameter.
       router.HandlerFunc(http.MethodGet, "/v1/tasks", app.listTasksHandler)

// Handle POST requests to create a new task.
       router.HandlerFunc(http.MethodPost, "/v1/tasks", app.createTaskHandler)

// Handle GET requests to fetch a specific task by ID.
       router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.showTaskHandler)

// Handle PATCH requests to update a specific task by ID.
// Note: Use PATCH instead of PUT for partial updates.
       router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", app.updateTaskHandler)

// Handle DELETE requests to delete a specific task by ID.
       router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", app.deleteTaskHandler)

	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// Add the route for the PUT /v1/users/activated endpoint.

	// Add the route for the POST /v1/tokens/authentication endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Add the enableCORS() middleware.
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
