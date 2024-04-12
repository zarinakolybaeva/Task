package data

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

// Define a custom ErrRecordNotFound error.
// We'll return this from our Get() method when looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Tasks       TaskModel
	Permissions PermissionModel // Add a new Permissions field.
	Tokens      TokenModel      // Add a new Tokens field.
	Users       UserModel       // Add a new Users field.
}

// For ease of use, we also add a New() method which returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// return Models{
	// 	Tasks:       TaskModel{
	// 		DB: db,
	// 		InfoLog  *log.Logger,
	// 		ErrorLog *log.Logger,
	// 	},
	// 	Permissions: PermissionModel{DB: db}, // Initialize a new PermissionModel instance.
	// 	Tokens:      TokenModel{DB: db},      // Initialize a new TokenModel instance.
	// 	Users:       UserModel{
	// 		DB: db,
	// 		InfoLog  *log.Logger,
	// 		ErrorLog *log.Logger,
	// 		},       // Initialize a new UserModel instance.
	// }
	return Models{
		Tasks:       TaskModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Permissions: PermissionModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Users: UserModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Tokens: TokenModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
