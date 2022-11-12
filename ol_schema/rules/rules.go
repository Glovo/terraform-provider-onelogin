package apprulesschema

import (
	"strconv"

	appruleactionsschema "github.com/glovo/terraform-provider-onelogin/ol_schema/rules/actions"
	appruleconditionsschema "github.com/glovo/terraform-provider-onelogin/ol_schema/rules/conditions"
	"github.com/glovo/terraform-provider-onelogin/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	apprules "github.com/onelogin/onelogin-go-sdk/pkg/services/apps/app_rules"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"match": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validMatch,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"position": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"conditions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: appruleconditionsschema.Schema(),
			},
		},
		"actions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: appruleactionsschema.Schema(),
			},
		},
	}
}

func validMatch(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"all", "any"})
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) apprules.AppRule {
	out := apprules.AppRule{}
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			out.ID = oltypes.Int32(int32(id))
		}
	}
	if s["app_id"] != nil {
		if appId, err := strconv.Atoi(s["app_id"].(string)); err == nil {
			out.AppID = oltypes.Int32(int32(appId))
		}
	}
	if name, notNil := s["name"].(string); notNil {
		out.Name = oltypes.String(name)
	}
	if match, notNil := s["match"].(string); notNil {
		out.Match = oltypes.String(match)
	}
	if pos, notNil := s["position"].(int); notNil {
		out.Position = oltypes.Int32(int32(pos))
	}
	if pos, notNil := s["enabled"].(bool); notNil {
		out.Enabled = oltypes.Bool(pos)
	}
	if s["conditions"] != nil {
		out.Conditions = []apprules.AppRuleConditions{}
		for _, val := range s["conditions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := appruleconditionsschema.Inflate(valMap)
			out.Conditions = append(out.Conditions, cond)
		}
	}
	if s["actions"] != nil {
		for _, val := range s["actions"].([]interface{}) {
			valMap := val.(map[string]interface{})
			cond := appruleactionsschema.Inflate(valMap)
			out.Actions = append(out.Actions, cond)
		}
	}
	return out
}
