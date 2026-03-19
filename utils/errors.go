package utils

import "errors"

var (
	ErrDuplicateEmail              = errors.New("email already exists")
	ErrUserNotFound                = errors.New("user not found")
	ErrShowerNotFound              = errors.New("shower not found")
	ErrUnauthorizedShowerAccess    = errors.New("unauthorized to access this shower")
	ErrCatalogAlreadyExists        = errors.New("catalog already exists for this shower")
	ErrPreferencesAlreadyExist     = errors.New("preferences already exist for this shower")
	ErrProductNotFound             = errors.New("product not found")
	ErrCatalogNotFound             = errors.New("catalog not found")
	ErrProductAlreadyInCatalog     = errors.New("product already exists in this catalog")
	ErrCatalogProductNotFound      = errors.New("catalog product not found")
	ErrProductIsExclusive          = errors.New("product is exclusive and already assigned to a catalog")
	ErrCatalogNotApproved          = errors.New("catalog has not been approved yet")
	ErrInvalidWebhookExternalID    = errors.New("invalid payment webhook external id")
	ErrCatalogProductAlreadyBought = errors.New("catalog product already bought")
	ErrInvalidProductID            = errors.New("invalid product id")
	ErrInvalidCustomerData         = errors.New("invalid customer data")
	ErrInvalidCurrentURL           = errors.New("invalid current url")
	ErrPaymentConfigurationMissing = errors.New("payment configuration missing")
	ErrPaymentProviderUnavailable  = errors.New("payment provider unavailable")
	ErrInvalidTaxID                = errors.New("invalid tax id")
)
