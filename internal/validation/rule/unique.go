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

// unique checks if the field doesn't exists in database
func unique(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ";")

	collection := database.Collection(params[0])

	if collection == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M

	if len(params) == 2 {
		filter = bson.M{
			params[1]: fl.Field().Interface(),
		}
	} else {
		filter = bson.M{
			strings.ToLower(fl.FieldName()): fl.Field().Interface(),
		}
	}

	count, _ := collection.CountDocuments(ctx, filter)

	return count == 0
}

func RegisterUniqueRule(validate *validator.Validate) {
	validate.RegisterValidation("unique", unique)

	validate.RegisterTranslation("unique", lang.Translator("es"), func(ut ut.Translator) error {
		return ut.Add("unique", "{0} debe ser unico.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique", fe.StructField())

		return t
	})

	validate.RegisterTranslation("unique", lang.Translator("en"), func(ut ut.Translator) error {
		return ut.Add("unique", "The {0} must be unique.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique", fe.StructField())

		return t
	})
}
