package services

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func translateValidationError(err error) map[string]string {
	errors := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, e := range ve {
		field := strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			errors[field] = "wajib diisi"
		case "min":
			errors[field] = "minimal " + e.Param() + " karakter"
		case "gt":
			errors[field] = "harus lebih besar dari 0"
		case "gte":
			errors[field] = "tidak boleh negatif"
		default:
			errors[field] = "tidak valid"
		}
	}

	return errors
}
