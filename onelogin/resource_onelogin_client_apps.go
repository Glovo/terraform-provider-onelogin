package onelogin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/glovo/onelogin-go-sdk/pkg/client"
	clientapps "github.com/glovo/onelogin-go-sdk/pkg/services/auth_servers/client_apps"
	authserversclientapps "github.com/glovo/terraform-provider-onelogin/ol_schema/auth_server/client_apps"
	"github.com/glovo/terraform-provider-onelogin/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ClientApps() *schema.Resource {
	return &schema.Resource{
		Schema:        authserversclientapps.Schema(),
		CreateContext: clientAppsCreate,
		ReadContext:   clientAppsRead,
		UpdateContext: clientAppsUpdate,
		DeleteContext: clientAppsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				sApiAuthID, sAppID, err := utils.ParseNestedResourceImportId(d.Id())
				if err != nil {
					return nil, err
				}
				d.SetId(fmt.Sprintf("%s:%s", sApiAuthID, sAppID))
				appID, err := strconv.Atoi(sAppID)
				if err != nil {
					return nil, err
				}
				d.Set("app_id", int32(appID))
				apiAuthID, err := strconv.Atoi(sApiAuthID)
				if err != nil {
					return nil, err
				}
				d.Set("api_auth_id", int32(apiAuthID))

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func clientAppsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientApp, _ := authserversclientapps.Inflate(map[string]interface{}{
		"app_id":      d.Get("app_id"),
		"api_auth_id": d.Get("api_auth_id"),
		"name":        d.Get("name"),
		"scope_ids":   d.Get("scope_ids"),
	})
	client := m.(*client.APIClient)
	err := client.Services.ClientAppsV2.Create(&clientApp)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem creating the Client App! %v", err)
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "[CREATED] Created Client App ID %d on Auth Server %d", *(clientApp.AppID), *(clientApp.APIAuthID))

	d.SetId(fmt.Sprintf("%d:%d", *(clientApp.APIAuthID), *(clientApp.AppID)))
	return scopesRead(ctx, d, m)
}

func clientAppsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientApp, _ := authserversclientapps.Inflate(map[string]interface{}{
		"app_id":      d.Get("app_id"),
		"api_auth_id": d.Get("api_auth_id"),
		"name":        d.Get("name"),
		"scope_ids":   d.Get("scope_ids"),
	})
	client := m.(*client.APIClient)
	err := client.Services.ClientAppsV2.Update(&clientApp)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem updating the Scope! %v", err)
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "[UPDATED] Updated Client App ID %d on Auth Server %d", *(clientApp.AppID), *(clientApp.APIAuthID))

	d.SetId(fmt.Sprintf("%d:%d", *(clientApp.APIAuthID), *(clientApp.AppID)))
	return scopesRead(ctx, d, m)
}

func clientAppsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.APIClient)
	apiAuthID, ok := d.Get("api_auth_id").(int)
	if !ok {
		tflog.Error(ctx, "[ERROR] Client Api Auth ID must be set!")
		return diag.Errorf("[ERROR] Client Api Auth ID must be set!")
	}
	appID, ok := d.Get("app_id").(int)
	if !ok {
		tflog.Error(ctx, "[ERROR] Client App ID must be set!")
		return diag.Errorf("[ERROR] Client App ID must be set!")
	}

	scopeQuery := clientapps.ClientAppsQuery{
		AuthServerID: strconv.Itoa(apiAuthID),
	}
	clientApps, err := client.Services.ClientAppsV2.Query(&scopeQuery)
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem reading Client App! %v", err)
		return diag.FromErr(err)
	}
	var clientApp *clientapps.ClientApp
	for _, app := range clientApps {
		if *app.AppID == int32(appID) {
			clientApp = &app
			break
		}
	}
	if clientApp == nil {
		d.SetId("")
		return nil
	}
	tflog.Info(ctx, "[READ] Reading Client App with ID %d on Auth Server %D", *(clientApp.APIAuthID), *(clientApp.AppID))
	var scopeIds []int32
	for _, scope := range clientApp.Scopes {
		scopeIds = append(scopeIds, *scope.ID)
	}

	d.Set("name", clientApp.Name)
	d.Set("scope_ids", scopeIds)

	return nil
}

func clientAppsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerID, ok := d.Get("auth_server_id").(int)
	if !ok {
		tflog.Error(ctx, "[ERROR] Client App ID must be set!")
		return diag.Errorf("[ERROR] Client App ID must be set!")
	}
	appID, ok := d.Get("app_id").(int)
	if !ok {
		tflog.Error(ctx, "[ERROR] Client App ID must be set!")
		return diag.Errorf("[ERROR] Client App ID must be set!")
	}

	client := m.(*client.APIClient)

	err := client.Services.ClientAppsV2.Destroy(int32(authServerID), int32(appID))
	if err != nil {
		tflog.Error(ctx, "[ERROR] There was a problem deleting the AuthServer! %v", err)
		return diag.FromErr(err)
	} else {
		tflog.Info(ctx, "[DELETED] Deleted Client App with %d on Auth Server %d", appID, authServerID)
		d.SetId("")
	}

	return nil
}
