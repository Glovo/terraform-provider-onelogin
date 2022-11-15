package appssoschema

import (
	"testing"

	"github.com/glovo/onelogin-go-sdk/pkg/oltypes"
	"github.com/glovo/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
)

func TestFlattenOIDCSSO(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppSso
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SSO fields from an OIDC app": {
			InputData: apps.AppSso{
				ClientID:     oltypes.String("test"),
				ClientSecret: oltypes.String("test"),
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     oltypes.String("test"),
				"client_secret": oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenOIDC(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAMLCert(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppSso
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SAML SSO Certificate fields for the given SAML app": {
			InputData: apps.AppSso{
				MetadataURL: oltypes.String("test"),
				AcsURL:      oltypes.String("test"),
				SlsURL:      oltypes.String("test"),
				Issuer:      oltypes.String("test"),
				Certificate: &apps.AppSsoCertificate{
					Name:  oltypes.String("test"),
					ID:    oltypes.Int32(123),
					Value: oltypes.String("test"),
				},
			},
			ExpectedOutput: map[string]interface{}{
				"name":  oltypes.String("test"),
				"value": oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAMLCert(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAML(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppSso
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SSO fields for a SAML app": {
			InputData: apps.AppSso{
				MetadataURL: oltypes.String("test"),
				AcsURL:      oltypes.String("test"),
				SlsURL:      oltypes.String("test"),
				Issuer:      oltypes.String("test"),
				Certificate: &apps.AppSsoCertificate{
					Name:  oltypes.String("test"),
					ID:    oltypes.Int32(123),
					Value: oltypes.String("test"),
				},
			},
			ExpectedOutput: map[string]interface{}{
				"metadata_url": oltypes.String("test"),
				"acs_url":      oltypes.String("test"),
				"sls_url":      oltypes.String("test"),
				"issuer":       oltypes.String("test"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAML(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
