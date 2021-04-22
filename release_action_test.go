package manifest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONMarshal(t *testing.T) {

	t.Run("ReturnsLowercaseName", func(t *testing.T) {

		action := EstafetteReleaseAction{
			Name: "deploy-canary",
		}

		// act
		bytes, err := json.Marshal(action)

		assert.Nil(t, err)
		assert.Equal(t, "{\"name\":\"deploy-canary\"}", string(bytes))
	})
}
