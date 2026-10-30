<div align="center">

# 🛠️ terraform-provider-onelogin

**Terraform provider to manage OneLogin apps, users, roles and auth servers as code.**

A maintained Glovo fork (from v0.2.0) of the community OneLogin provider. It is a Terraform plugin, not a service: it runs on the machine that runs `terraform`, and it is published as GitHub release artifacts.

[![Go](https://img.shields.io/badge/Go-1.17-00ADD8?logo=go&logoColor=white)](go.mod)
[![Terraform Plugin SDK](https://img.shields.io/badge/Plugin_SDK-v2.7.0-844FBA?logo=terraform&logoColor=white)](go.mod)
[![Protocol](https://img.shields.io/badge/Plugin_protocol-5.0-844FBA?logo=terraform&logoColor=white)](terraform-registry-manifest.json)
[![Release](https://img.shields.io/badge/Release-GoReleaser-blue)](.goreleaser.yml)
<br/>
[![test files](https://img.shields.io/badge/test_files-36-brightgreen)](#testing)
[![Go Report Card](https://goreportcard.com/badge/github.com/glovo/terraform-provider-onelogin)](https://goreportcard.com/report/github.com/glovo/terraform-provider-onelogin)

</div>

---

## Table of contents

- [Overview](#overview)
- [Adding a resource](#adding-a-resource)
- [Prerequisites & local setup](#prerequisites--local-setup)
- [Testing](#testing)
- [CI/CD](#cicd)
- [Deployment & rollback](#deployment--rollback)
- [Configuration](#configuration)
- [Repository layout](#repository-layout)
- [Contributing & conventions](#contributing--conventions)
- [Ownership & support](#ownership--support)

---

## Overview

The provider lets you describe the final state of your OneLogin organisation (apps, users, role associations and more) in `.tf` files, then calls the OneLogin API to make that state real. It is built on the Terraform Plugin SDK v2 and talks to OneLogin through Glovo's fork of the OneLogin Go SDK (`github.com/glovo/onelogin-go-sdk`, see [`go.mod`](go.mod)).

What this fork adds over the community provider:

- Bugfixes for auth servers
- Compatibility with Terraform 1.x
- Auth server features: **scopes** and **client apps**

| | |
|---|---|
| Language / build | Go 1.17 ([`go.mod`](go.mod)) |
| Registry address | `Glovo/onelogin` ([`docs/index.md`](docs/index.md)) |
| External system | OneLogin API, via `github.com/glovo/onelogin-go-sdk` v1.1.23 |
| Docs | [`docs/`](docs): 1 index page, 2 data sources, 12 resources |

### Resources and data sources

Names registered in [`onelogin/provider.go`](onelogin/provider.go):

| Kind | Names |
|---|---|
| Data sources | `onelogin_user`, `onelogin_users` |
| Resources | `onelogin_apps`, `onelogin_oidc_apps`, `onelogin_saml_apps`, `onelogin_app_rules`, `onelogin_app_role_attachments`, `onelogin_user_mappings`, `onelogin_users`, `onelogin_auth_servers`, `onelogin_scopes`, `onelogin_client_apps`, `onelogin_roles`, `onelogin_smarthooks`, `onelogin_smarthook_environment_variables`, `onelogin_privileges` |

<details><summary><b>Terraform for non-users</b></summary>

&nbsp;

Wildly simplified: Terraform lets you describe the final state of all your OneLogin things (apps, users, associations via roles etc.) in a `.tf` file using HashiCorp Language (HCL). It then fires off the API requests needed to make that desired state real. It also tracks the known state of OneLogin in `.tfstate` and uses that as the source of truth.

```bash
brew install terraform   # install
terraform init           # set up the tfstate and prepare to track
terraform plan           # show the diff between current and desired state
terraform apply          # do the actual work of updating OneLogin
```

</details>

---

## Adding a resource

Adding a new resource usually needs these steps, in order:

1. Make sure the service exists in the [OneLogin SDK](https://github.com/onelogin/onelogin-go-sdk). See `AppsService` for an example.
2. Register the resource in [`onelogin/provider.go`](onelogin/provider.go), the same way the other resources are.
3. Create `onelogin/resource_onelogin_<resource>.go` and `onelogin/resource_onelogin_<resource>_test.go`. This layer is the interface between the Terraform and OneLogin SDKs and makes the actual requests to OneLogin. The test file is used by Terraform for the acceptance tests. Copy an existing pair.
4. Put the pack/unpack logic (JSON to Go struct and back) in `ol_schema/<resource>/<sub-resource>`. See [`ol_schema/`](ol_schema).
5. Add examples `examples/onelogin_<resource>_example.tf` and `examples/onelogin_<resource>_updated_example.tf`. The acceptance tests use them to check that apply works. See [`examples/`](examples).
6. Add a doc page at `docs/resources/onelogin_<resource>.md`. See [`docs/resources/`](docs/resources).

To debug, set `TF_LOG=trace` to see all loggers in the Terraform workflow:

```bash
export TF_LOG=trace
```

---

## Prerequisites & local setup

| Tool | Why |
|---|---|
| Go | Builds the provider (`go 1.17` in [`go.mod`](go.mod); CI reads the version from it) |
| Terraform v0.12.24 or later | Runs the provider (minimum as stated in the previous README) |
| OneLogin admin access | To create API credentials with the **manage all** permission |

If you sideload the provider (not installed from the registry), clone this repository and build it:

```bash
git clone https://github.com/Glovo/terraform-provider-onelogin.git
cd terraform-provider-onelogin
go build .                      # builds the provider binary
go mod download                 # installs dependencies
```

With Terraform v0.13.x or later, use this configuration for the sideloaded provider:

```hcl
terraform {
  required_providers {
    onelogin = {
      source  = "Glovo/onelogin"
      version = "0.2.1"
    }
  }
}

provider "onelogin" {
  # Configuration options
}
```

Then export your credentials (see [Configuration](#configuration)) and use the usual `terraform` commands.

---

## Testing

```bash
go test -v -cover ./...   # what CI runs
```

- There are 36 Go test files (counted as files, not test cases), in [`onelogin/`](onelogin), [`ol_schema/`](ol_schema) and [`utils/`](utils).
- Each resource has a test file next to it in [`onelogin/`](onelogin). These are Terraform acceptance tests that use the `.tf` files in [`examples/`](examples).

> [!WARNING]
> Acceptance testing is switched off in CI: `TF_ACC` is commented out in [`go.yml`](.github/workflows/go.yml), as is the Terraform version matrix. CI runs the tests with dummy credentials, so they do not call OneLogin.

---

## CI/CD

| Workflow | Trigger | Effect |
|---|---|---|
| [`go.yml`](.github/workflows/go.yml) (Tests) | Pull requests, and pushes to `main`, that touch `examples/`, `ol_schema/`, `main.go`, `onelogin/`, `utils/` or the workflow | `go build`, `go generate ./...` plus a clean `git diff` check, then `go test -v -cover ./...` |
| [`release.yml`](.github/workflows/release.yml) (Release) | A tag matching `v*`, or manual run | Imports the GPG key, then runs GoReleaser to build and sign the release artifacts |

> [!NOTE]
> The `paths` list in `go.yml` points at `.github/worflows/go.yml` (typo), so edits to the workflow file itself do not trigger it.

---

## Deployment & rollback

The provider is released by tag; nothing is deployed to a server.

1. Merge your change to `main` once it is approved and CI passes.
2. Smoke test, and make any last-minute updates.
3. Create a version tag `vX.X.X` (for example `v0.2.6`, the latest tag) and push it: `git tag <version> && git push --tags`.
4. The tag triggers the [Release workflow](.github/workflows/release.yml). GoReleaser builds for freebsd, windows, linux and darwin (amd64, 386, arm, arm64; darwin has no 386), zips the binaries, signs the SHA256 checksums with GPG and attaches them, plus `terraform-registry-manifest.json`, to the GitHub release ([`.goreleaser.yml`](.goreleaser.yml)).
5. Open the [Releases](https://github.com/Glovo/terraform-provider-onelogin/releases) page and check that the build artifacts are attached to your version.

Release signing needs the repository secrets `GPG_PRIVATE_KEY` and `PASSPHRASE`.

---

## Configuration

The provider reads its settings from the provider block or from the environment. Never commit credentials.

| Provider argument | Environment variable | Notes |
|---|---|---|
| `client_id` | `ONELOGIN_CLIENT_ID` | Required |
| `client_secret` | `ONELOGIN_CLIENT_SECRET` | Required |
| `url` | `ONELOGIN_OAPI_URL` | Optional. The API URL for your region |
| `region` | none | Optional. Defaults to the SDK's US region |

```bash
export ONELOGIN_CLIENT_ID=<your client id>
export ONELOGIN_CLIENT_SECRET=<your client secret>
export ONELOGIN_OAPI_URL=<the api url for your region>
```

---

## Repository layout

```text
.
├── main.go                         # plugin entry point (plugin.Serve)
├── onelogin/ ⭐                     # provider, resources, data sources, acceptance tests
├── ol_schema/                      # pack/unpack of OneLogin JSON <-> Go structs, per resource
├── utils/                          # shared helpers
├── examples/                       # .tf files used by the acceptance tests (25 files)
├── docs/                           # Terraform registry docs (index, data-sources/, resources/)
├── terraform-registry-manifest.json # registry metadata (protocol 5.0)
├── .goreleaser.yml                 # release build and signing
└── .github/workflows/              # go.yml (tests), release.yml (release)
```

---

## Contributing & conventions

- Open a pull request against `main`. CI runs on it when Go, example or workflow-relevant paths change.
- Keep `go generate ./...` output committed: CI fails on any diff after generation.
- Every new resource needs code, test, examples and a doc page (see [Adding a resource](#adding-a-resource)).
- Dependencies are managed with Go modules.

<details><summary><b>More rules</b></summary>

&nbsp;

Dependency management with `go mod`:

```bash
go get -u "package-name"   # add a package
rm go.sum && go mod download   # re-install dependencies for this project
go get -u ./...            # update dependencies
```

</details>

---

## Ownership & support

- **Team:** no `CODEOWNERS` file or catalog entry exists in this repository.
- **Issues and releases:** [Glovo/terraform-provider-onelogin](https://github.com/Glovo/terraform-provider-onelogin).
- **Code quality:** [Go Report Card](https://goreportcard.com/report/github.com/glovo/terraform-provider-onelogin).
