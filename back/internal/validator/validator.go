package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func New() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := fld.Tag.Get("json")
		if tag == "" {
			tag = fld.Tag.Get("query")
		}
		if tag == "" {
			return ""
		}

		name := strings.SplitN(tag, ",", 2)[0] //nolint:mnd // это не магическое число

		// Пропустить, если ключ тега говорит, что его следует игнорировать
		if name == "-" {
			return ""
		}

		return name
	})

	return validate
}

// ParseError преобразует ошибку валидации в HTTP-ошибку.
func ParseError(err error) error {
	var invalidErr *validator.InvalidValidationError
	if errors.As(err, &invalidErr) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return fmt.Errorf("validate: %v", err)
	}

	return fiber.NewError(fiber.StatusBadRequest, `invalid field "`+validationErrs[0].Field()+`"`)
}
