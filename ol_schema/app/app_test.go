package appschema

import (
	"testing"

	"github.com/glovo/onelogin-go-sdk/pkg/oltypes"
	"github.com/glovo/onelogin-go-sdk/pkg/services/apps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of an App Schema", func(t *testing.T) {
		schema := Schema()
		assert.NotNil(t, schema["name"])
		assert.NotNil(t, schema["visible"])
		assert.NotNil(t, schema["description"])
		assert.NotNil(t, schema["notes"])
		assert.NotNil(t, schema["icon_url"])
		assert.NotNil(t, schema["auth_method"])
		assert.NotNil(t, schema["policy_id"])
		assert.NotNil(t, schema["brand_id"])
		assert.NotNil(t, schema["allow_assumed_signin"])
		assert.NotNil(t, schema["tab_id"])
		assert.NotNil(t, schema["connector_id"])
		assert.NotNil(t, schema["created_at"])
		assert.NotNil(t, schema["updated_at"])
		assert.NotNil(t, schema["provisioning"])
		assert.NotNil(t, schema["parameters"])
	})
}

func mockSetFn(interface{}) int {
	return 0
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput apps.App
	}{
		"creates and returns the address of an App struct with all sub-fields": {
			ResourceData: map[string]interface{}{
				"id":                   "123",
				"name":                 "test",
				"visible":              true,
				"description":          "test",
				"notes":                "test",
				"allow_assumed_signin": true,
				"connector_id":         123,
				"brand_id":             123,
				"parameters": schema.NewSet(mockSetFn, []interface{}{
					map[string]interface{}{
						"param_key_name":             "test",
						"param_id":                   123,
						"label":                      "test",
						"user_attribute_mappings":    "test",
						"user_attribute_macros":      "test",
						"attributes_transformations": "test",
						"default_values":             "test",
						"skip_if_blank":              true,
						"values":                     "test",
						"provisioned_entitlements":   true,
						"safe_entitlements_enabled":  true,
					},
				}),
				"provisioning": map[string]interface{}{
					"enabled": true,
				},
				"configuration": map[string]interface{}{
					"provider_arn":        "test",
					"signature_algorithm": "test",
				},
			},
			ExpectedOutput: apps.App{
				Name:               oltypes.String("test"),
				ID:                 oltypes.Int32(int32(123)),
				Visible:            oltypes.Bool(true),
				Description:        oltypes.String("test"),
				Notes:              oltypes.String("test"),
				AllowAssumedSignin: oltypes.Bool(true),
				ConnectorID:        oltypes.Int32(123),
				BrandID:            oltypes.Int32(123),
				Parameters: map[string]apps.AppParameters{
					"test": apps.AppParameters{
						ID:                        oltypes.Int32(123),
						Label:                     oltypes.String("test"),
						UserAttributeMappings:     oltypes.String("test"),
						UserAttributeMacros:       oltypes.String("test"),
						AttributesTransformations: oltypes.String("test"),
						DefaultValues:             oltypes.String("test"),
						SkipIfBlank:               oltypes.Bool(true),
						Values:                    oltypes.String("test"),
						ProvisionedEntitlements:   oltypes.Bool(true),
						SafeEntitlementsEnabled:   oltypes.Bool(true),
					},
				},
				Provisioning: &apps.AppProvisioning{
					Enabled: oltypes.Bool(true),
				},
				Configuration: &apps.AppConfiguration{
					ProviderArn:        oltypes.String("test"),
					SignatureAlgorithm: oltypes.String("test"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}
