---
layout: "onelogin"
page_title: "OneLogin: onelogin_scopes"
sidebar_current: "docs-onelogin-resource-onelogin_scopes"
description: |-
  Creates a Scope within an Authentication Server Resource.
---

# onelogin_scope

Creates a Scope Resource.

This resource allows you to create and configure an Scope for an Authentication Server.

## Example Usage

```hcl
resource onelogin_scopes example {
  auth_server_id = 121
  value = "myresource:private"
  description = "Grants private access to myresource"
}
```

## Argument Reference

The following arguments are supported:
* `auth_server_id` - (Required) The ID of the authorization resource.

* `value` - (Required) The value of the scope.

* `description` - (Required) A brief description about the scope.


## Attributes Reference

No further attributes are exported

## Import

A Scope can be imported via the Authorization Server ID and the OneLogin Scope ID.

```
$ terraform import onelogin_scopes.example <auth_server_id>:<scope_id>
```
