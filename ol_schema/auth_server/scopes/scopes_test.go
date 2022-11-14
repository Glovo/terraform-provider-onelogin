package authserverscopesschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers/scopes"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of an Scope Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["auth_server_id"])
		assert.NotNil(t, provSchema["description"])
		assert.NotNil(t, provSchema["value"])
	})
}
func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput scopes.Scope
	}{
		"creates an auth server": {
			ResourceData: map[string]interface{}{
				"auth_server_id": 1,
				"value":          "value",
				"description":    "description",
			},
			ExpectedOutput: scopes.Scope{
				AuthServerID: oltypes.Int32(1),
				Value:        oltypes.String("value"),
				Description:  oltypes.String("description"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
