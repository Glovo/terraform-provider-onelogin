package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/glovo/onelogin-go-sdk/pkg/client"
	"github.com/glovo/onelogin-go-sdk/pkg/services/auth_servers/scopes"
	authserverscopesschema "github.com/glovo/terraform-provider-onelogin/ol_schema/auth_server/scopes"
	"github.com/glovo/terraform-provider-onelogin/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Scopes() *schema.Resource {
	return &schema.Resource{
		Schema:        authserverscopesschema.Schema(),
		CreateContext: scopesCreate,
		ReadContext:   scopesRead,
		UpdateContext: scopesUpdate,
		DeleteContext: scopesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				authServerID, scopeID, err := utils.ParseNestedResourceImportId(d.Id())
				if err != nil {
					return nil, err
				}
				d.SetId(scopeID)
				auth_id, err := strconv.Atoi(authServerID)
				if err != nil {
					return nil, err
				}
				d.Set("auth_server_id", int32(auth_id))

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func scopesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope, _ := authserverscopesschema.Inflate(map[string]interface{}{
		"auth_server_id": d.Get("auth_server_id"),
		"value":          d.Get("value"),
		"description":    d.Get("description"),
	})
	client := m.(*client.APIClient)
	err := client.Services.ScopesV2.Create(&scope)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem creating the Scope! %v", err)
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "[CREATED] Created Scope with %d", *(scope.ID))

	d.SetId(fmt.Sprintf("%d", *(scope.ID)))
	return scopesRead(ctx, d, m)
}

func scopesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope, _ := authserverscopesschema.Inflate(map[string]interface{}{
		"id":             d.Id(),
		"auth_server_id": d.Get("auth_server_id"),
		"value":          d.Get("value"),
		"description":    d.Get("description"),
	})
	client := m.(*client.APIClient)
	err := client.Services.ScopesV2.Update(&scope)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem updating the Scope! %v", err)
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "[UPDATED] Updated Scope with %d", *(scope.ID))

	d.SetId(fmt.Sprintf("%d", *(scope.ID)))
	return scopesRead(ctx, d, m)
}

func scopesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.APIClient)
	uid, _ := strconv.Atoi(d.Id())
	authServerID, ok := d.Get("auth_server_id").(int)
	if !ok {
		tflog.Error(ctx, "[ERROR] Auth Server ID must be set!")
		return diag.Errorf("[ERROR] Auth Server ID must be set!")
	}

	scopeQuery := scopes.ScopesQuery{
		AuthServerID: strconv.Itoa(authServerID),
	}
	serverScopes, err := client.Services.ScopesV2.Query(&scopeQuery)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem reading Scopes! %v", err)
		return diag.FromErr(err)
	}
	var scope *scopes.Scope
	for _, serverScope := range serverScopes {
		if *serverScope.ID == int32(uid) {
			scope = &serverScope
			break
		}
	}
	if scope == nil {
		d.SetId("")
		return nil
	}
	tflog.Info(ctx, "[READ] Reading AuthServer with %d", *(scope.ID))

	d.Set("value", scope.Value)
	d.Set("description", scope.Description)

	return nil
}

func scopesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uid, _ := strconv.Atoi(d.Id())
	authServerID, ok := d.Get("auth_server_id").(int)
	if !ok {
		tflog.Error(ctx, "[ERROR] Auth Server ID must be set!")
		return diag.Errorf("[ERROR] Auth Server ID must be set!")
	}

	client := m.(*client.APIClient)

	err := client.Services.ScopesV2.Destroy(int32(authServerID), int32(uid))
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem deleting the AuthServer! %v", err)
		return diag.FromErr(err)
	} else {
		tflog.Info(ctx, "[DELETED] Deleted AuthServer with %d", uid)
		d.SetId("")
	}

	return nil
}
