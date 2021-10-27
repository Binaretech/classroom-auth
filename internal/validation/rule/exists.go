package rule

import (
	"context"
	"strings"
	"time"

	ut "github.com/go-playground/universal-translator"

	"github.com/Binaretech/classroom-auth/internal/database"
	"github.com/Binaretech/classroom-auth/internal/lang"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

// exists checks if the field exists in database
func exists() func(validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		params := strings.Split(fl.Param(), ";")

		collection := database.Collection(params[0])

		if collection == nil {
			return false
		}

		var filter bson.M

		if len(params) == 2 {
			filter = bson.M{params[1]: fl.Field().String()}
		} else {
			filter = bson.M{strings.ToLower(fl.FieldName()): fl.Field().String()}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		count, _ := collection.CountDocuments(ctx, filter)

		return count > 0
	}
}

func RegisterExistsRule(validate *validator.Validate) {
	validate.RegisterValidation("exists", exists())

	validate.RegisterTranslation("exists", lang.Translator("es"), func(ut ut.Translator) error {
		return ut.Add("exists", "El {0} no existe.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.StructField())

		return t
	})

	validate.RegisterTranslation("exists", lang.Translator("en"), func(ut ut.Translator) error {
		return ut.Add("exists", "The {0} does not exist.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.StructField())

		return t
	})
}
