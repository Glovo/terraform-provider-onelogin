package authserverschema

import (
	"strconv"

	"github.com/glovo/onelogin-go-sdk/pkg/oltypes"
	authservers "github.com/glovo/onelogin-go-sdk/pkg/services/auth_servers"
	authserverconfigurationschema "github.com/glovo/terraform-provider-onelogin/ol_schema/auth_server/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema returns a key/value map of the various fields that make up an Auth Server at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"configuration": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: authserverconfigurationschema.Schema(),
			},
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin AuthServer.
func Inflate(s map[string]interface{}) (authservers.AuthServer, error) {
	var err error
	authServer := authservers.AuthServer{
		Name:        oltypes.String(s["name"].(string)),
		Description: oltypes.String(s["description"].(string)),
	}
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			authServer.ID = oltypes.Int32(int32(id))
		}
	}
	var cfg []interface{}
	if s["configuration"] != nil {
		cfg = s["configuration"].([]interface{})
	}
	if cfg != nil && len(cfg) > 0 {
		var conf authservers.AuthServerConfiguration
		conf = authserverconfigurationschema.Inflate(cfg[0].(map[string]interface{}))
		authServer.Configuration = &conf
	}
	return authServer, err
}
