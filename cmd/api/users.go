// package main

// import (
// 	"errors"
// 	"net/http"
// 	"time"

// 	"github.com/Zarina/TaskNinja/internal/data"
// 	"github.com/Zarina/TaskNinja/internal/validator"
// )

// func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Name     string `json:"name"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
// 	user := &data.User{
// 		Name:      input.Name,
// 		Email:     input.Email,
// 		Activated: false,
// 	}
// 	err = user.Password.Set(input.Password)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	v := validator.New()
// 	if data.ValidateUser(v, user); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	err = app.models.Users.Insert(user)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrDuplicateEmail):
// 			v.AddError("email", "a user with this email address already exists")
// 			app.failedValidationResponse(w, r, v.Errors)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	// Add the "movies:read" permission for the new user.
// 	err = app.models.Permissions.AddForUser(user.ID, "tasks:read")
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	app.background(func() {
// 		data := map[string]interface{}{
// 			"activationToken": token.Plaintext,
// 			"userID":          user.ID,
// 		}
// 		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
// 		if err != nil {
// 			app.logger.PrintError(err, nil)
// 		}
// 	})
// 	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse the plaintext activation token from the request body.
// 	var input struct {
// 		TokenPlaintext string `json:"token"`
// 	}
// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Validate the plaintext token provided by the client.
// 	v := validator.New()
// 	// if err := v.ValidateTokenPlaintext(input.TokenPlaintext); err != nil {
// 	//     app.failedValidationResponse(w, r, err)
// 	//     return
// 	// }
// 	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	// Retrieve the details of the user associated with the token.
// 	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
// 	if err != nil {
// 		if errors.Is(err, data.ErrRecordNotFound) {
// 			v.AddError("token", "invalid or expired activation token")
// 			app.failedValidationResponse(w, r, v.Errors)
// 		} else {
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Update the user's activation status.
// 	user.Activated = true

// 	// Save the updated user record in the database.
// 	err = app.models.Users.Update(user)
// 	if err != nil {
// 		if errors.Is(err, data.ErrEditConflict) {
// 			app.editConflictResponse(w, r)
// 		} else {
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Delete all activation tokens for the user.
// 	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Send the updated user details to the client in a JSON response.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }


package main

import (
	"errors"
	"net/http"
	"time"
 	"github.com/Zarina/TaskNinja/internal/data"
 	"github.com/Zarina/TaskNinja/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that
	// we set the Activated field to false, which isn't strictly necessary because
	// the Activated field will have the zero-value of false by default. But setting
	// this explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if
	// any of the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get an ErrDuplicateEmail error, use the v.AddError() method to manually add
		// a message to the validator instance, and then call our failedValidationResponse
		// helper().
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")

			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "menus:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *data.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	app.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
}

// activateUserHandler activates a user by setting 'activation = true' using the provided
// activation token in the request body.
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the details of the user associated with the token using the GetForToken() method.
	// If no matching record is found, then we let the client know that the token they provided
	// is not valid.
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in the same
	// way that we did for our move records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully above, then delete all activation tokens for the user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
}
