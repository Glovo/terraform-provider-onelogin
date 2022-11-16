package authserversclientapps

import (
	"github.com/glovo/onelogin-go-sdk/pkg/oltypes"
	clientapps "github.com/glovo/onelogin-go-sdk/pkg/services/auth_servers/client_apps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema returns a key/value map of the various fields that make up a Scope at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"api_auth_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"scope_ids": {
			Type:     schema.TypeSet,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin Scope.
func Inflate(s map[string]interface{}) (clientapps.ClientApp, error) {
	var err error
	clientApp := clientapps.ClientApp{
		APIAuthID: oltypes.Int32(int32(s["api_auth_id"].(int))),
		AppID:     oltypes.Int32(int32(s["app_id"].(int))),
	}
	if val, notNil := s["name"].(string); notNil {
		clientApp.Name = oltypes.String(val)
	}
	if s["scope_ids"] != nil {
		clientApp.ScopeIDs = make([]int32, len(s["scope_ids"].(*schema.Set).List()))
		for i, appID := range s["scope_ids"].(*schema.Set).List() {
			clientApp.ScopeIDs[i] = int32(appID.(int))
		}
	}
	return clientApp, err
}
