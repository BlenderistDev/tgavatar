package telegram

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_GetClient(t *testing.T) {
	err := os.Setenv("APP_ID", "111")
	assert.Nil(t, err)
	err = os.Setenv("APP_HASH", "123")
	assert.Nil(t, err)

	factory := NewFactory()
	client, err := factory.GetClient()

	assert.Nil(t, err)
	assert.IsType(t, &TGClient{}, client)
}

func TestFactory_GetClient_MissingAppID(t *testing.T) {
	os.Clearenv()
	factory := NewFactory()
	client, err := factory.GetClient()

	assert.Nil(t, client)
	assert.Equal(t, "telegram TGClient creating error: APP_ID not set or invalid: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestFactory_GetClient_MissingAppHash(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("APP_ID", "111")
	assert.Nil(t, err)

	factory := NewFactory()
	client, err := factory.GetClient()

	assert.Nil(t, client)
	assert.Equal(t, "telegram TGClient creating error: no APP_HASH provided", err.Error())
}
