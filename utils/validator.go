package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string)
		for _, fe := range ve {
			out[fe.Field()] = messageForTag(fe)
		}
		return out
	}

	return map[string]string{"error": err.Error()}
}

func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Too short"
	default:
		return "Invalid value"
	}
}

func BindAndValidate(ctx *gin.Context, req any) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		errorsMap := FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorsMap})
		return false
	}
	return true
}
