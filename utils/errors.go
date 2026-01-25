package utils

import "errors"

var (
	ErrDuplicateEmail           = errors.New("email already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrShowerNotFound           = errors.New("shower not found")
	ErrUnauthorizedShowerAccess = errors.New("unauthorized to access this shower")
	ErrCatalogAlreadyExists     = errors.New("catalog already exists for this shower")
	ErrPreferencesAlreadyExist  = errors.New("preferences already exist for this shower")
	ErrProductNotFound          = errors.New("product not found")
	ErrCatalogNotFound          = errors.New("catalog not found")
	ErrProductAlreadyInCatalog  = errors.New("product already exists in this catalog")
	ErrCatalogProductNotFound   = errors.New("catalog product not found")
)
