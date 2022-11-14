package authserverscopesschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers/scopes"
)

// Schema returns a key/value map of the various fields that make up a Scope at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auth_server_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin Scope.
func Inflate(s map[string]interface{}) (scopes.Scope, error) {
	var err error
	scope := scopes.Scope{
		AuthServerID: oltypes.Int32(int32(s["auth_server_id"].(int))),
		Value:        oltypes.String(s["value"].(string)),
		Description:  oltypes.String(s["description"].(string)),
	}

	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			scope.ID = oltypes.Int32(int32(id))
		}
	}
	return scope, err
}
