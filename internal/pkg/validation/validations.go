package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"todo/internal/db"
)

func IsUniqueValue(tableName, fieldName, value string) bool {
	var count int64

	err := db.Db.QueryRow("SELECT COUNT(*) FROM "+tableName+" WHERE "+fieldName+" = ?", value).Scan(&count)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return count > 0
}

func IsExistValue(tableName, fieldName string, value interface{}) bool {
	var count int64

	err := db.Db.QueryRow("SELECT COUNT(*) FROM "+tableName+" WHERE "+fieldName+" = ?", value).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
}

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)

	for _, err := range errs {
		fmt.Println()
		switch err.Tag() {
		case "required":
			errorMessages[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "email":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be a valid email address", err.Field())
		case "min":
			errorMessages[err.Field()] = fmt.Sprintf("%s must have at least %s characters", err.Field(), err.Param())
		case "max":
			errorMessages[err.Field()] = fmt.Sprintf("%s must have at most %s characters", err.Field(), err.Param())
		case "gt":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
		case "gte":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
		default:
			errorMessages[err.Field()] = fmt.Sprintf("Validation validations on field %s", err.Field())
		}
	}

	return errorMessages
}
