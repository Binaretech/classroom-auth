package validation

import (
	"net/http"

	"github.com/Binaretech/classroom-auth/lang"
	"github.com/Binaretech/classroom-auth/utils"
	"github.com/Binaretech/classroom-auth/validation/rule"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

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

type Validator struct {
	validate *validator.Validate
}

func newValidator(validate *validator.Validate) *Validator {
	return &Validator{validate: validate}
}

func (v *Validator) Validate(data any) error {
	errors := map[string][]string{}

	if err := v.validate.Struct(data); err != nil {
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

	return echo.NewHTTPError(http.StatusUnprocessableEntity, &ErrorResponse{ValidationErrors: errors})
}

// SetUpValidator configures and returns an instance of `validator.Validate`
func SetUpValidator(db *mongo.Database) *Validator {
	validate := validator.New()

	rule.RegisterExistsRule(db, validate)
	rule.RegisterUniqueRule(db, validate)

	es_translations.RegisterDefaultTranslations(validate, lang.Translator("es"))
	en_translations.RegisterDefaultTranslations(validate, lang.Translator("en"))

	return newValidator(validate)
}
