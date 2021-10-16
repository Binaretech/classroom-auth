package validation

import (
	"github.com/Binaretech/classroom-auth/internal/lang"
	"github.com/Binaretech/classroom-auth/internal/utils"
	"github.com/Binaretech/classroom-auth/internal/validation/rule"

	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	"github.com/spf13/viper"
)

// ErrorResponse store the errors after a validation process
type ErrorResponse struct {
	// ValidationErrors store in the form of name -> []errors
	ValidationErrors map[string][]string `json:"validationErrors"`
}

// SetUpValidator configures and returns an instance of `validator.Validate`
func SetUpValidator() *validator.Validate {
	validate := validator.New()

	rule.RegisterExistsRule(validate)
	rule.RegisterUniqueRule(validate)

	es_translations.RegisterDefaultTranslations(validate, lang.Translator("es"))
	en_translations.RegisterDefaultTranslations(validate, lang.Translator("en"))

	return validate
}

// Struct validate struct and return an ErrorResponse if there are a validation error
func Struct(request interface{}) *ErrorResponse {
	errors := map[string][]string{}
	validate := SetUpValidator()

	if err := validate.Struct(request); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			name := utils.LowerCaseInitial(err.Field())
			if _, ok := errors[name]; !ok {
				errors[name] = []string{}
			}

			errors[name] = append(errors[name], err.Translate(lang.Translator(viper.GetString("locale"))))
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return &ErrorResponse{ValidationErrors: errors}
}
