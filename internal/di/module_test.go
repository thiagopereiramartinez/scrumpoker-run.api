package di

import (
	"cloud.google.com/go/firestore"
	"github.com/golobby/container"
	Assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupDependencies(t *testing.T) {

	assert := Assert.New(t)
	err := SetupDependencies()

	assert.NoError(err)

	var client = new(firestore.Client)
	container.Make(&client)

	assert.NotNil(client)
}
