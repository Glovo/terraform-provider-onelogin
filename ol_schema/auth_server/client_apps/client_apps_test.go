package authserversclientapps

import (
	"testing"

	"github.com/glovo/onelogin-go-sdk/pkg/oltypes"
	clientapps "github.com/glovo/onelogin-go-sdk/pkg/services/auth_servers/client_apps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func mockSetFn(i interface{}) int {
	return i.(int)
}

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of an Scope Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["api_auth_id"])
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["scope_ids"])
	})
}
func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput clientapps.ClientApp
	}{
		"creates an client app on the auth server": {
			ResourceData: map[string]interface{}{
				"app_id":      2,
				"api_auth_id": 1,
				"name":        "value",
				"scope_ids":   schema.NewSet(mockSetFn, []interface{}{2, 3}),
			},
			ExpectedOutput: clientapps.ClientApp{
				AppID:     oltypes.Int32(2),
				APIAuthID: oltypes.Int32(1),
				Name:      oltypes.String("value"),
				ScopeIDs:  []int32{2, 3},
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
