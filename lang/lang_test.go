package lang

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTranslator(t *testing.T) {
	viper.Set("lang", "es")
	assert.Equal(t, "Ha ocurrido un error en el servidor.", Trans("internal error"))

	viper.Set("lang", "en")
	assert.Equal(t, "An unexpected error has occurred.", Trans("internal error"))
}
