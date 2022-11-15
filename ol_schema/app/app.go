package appschema

import (
	"strconv"

	"github.com/glovo/onelogin-go-sdk/pkg/oltypes"
	"github.com/glovo/onelogin-go-sdk/pkg/services/apps"
	appconfigurationschema "github.com/glovo/terraform-provider-onelogin/ol_schema/app/configuration"
	appparametersschema "github.com/glovo/terraform-provider-onelogin/ol_schema/app/parameters"
	appprovisioningschema "github.com/glovo/terraform-provider-onelogin/ol_schema/app/provisioning"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schema returns a key/value map of the various fields that make up an App at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"visible": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"policy_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"brand_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"allow_assumed_signin": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tab_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"connector_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeBool},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: appparametersschema.Schema(),
			},
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin App.
func Inflate(s map[string]interface{}) (apps.App, error) {
	var err error
	app := apps.App{
		Name:               oltypes.String(s["name"].(string)),
		Description:        oltypes.String(s["description"].(string)),
		Notes:              oltypes.String(s["notes"].(string)),
		ConnectorID:        oltypes.Int32(int32(s["connector_id"].(int))),
		Visible:            oltypes.Bool(s["visible"].(bool)),
		AllowAssumedSignin: oltypes.Bool(s["allow_assumed_signin"].(bool)),
	}
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			app.ID = oltypes.Int32(int32(id))
		}
	}
	if s["brand_id"] != nil {
		brandID := s["brand_id"].(int)
		app.BrandID = oltypes.Int32(int32(brandID))
	}
	if s["parameters"] != nil {
		p := s["parameters"].(*schema.Set).List()
		app.Parameters = make(map[string]apps.AppParameters, len(p))
		for _, val := range p {
			valMap := val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = appparametersschema.Inflate(valMap)
		}
	}
	if s["provisioning"] != nil {
		prov := appprovisioningschema.Inflate(s["provisioning"].(map[string]interface{}))
		app.Provisioning = &prov
	}
	if s["configuration"] != nil {
		var conf apps.AppConfiguration
		conf, err = appconfigurationschema.Inflate(s["configuration"].(map[string]interface{}))
		app.Configuration = &conf
	}
	return app, err
}
