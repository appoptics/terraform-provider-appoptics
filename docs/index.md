---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appoptics Provider"
subcategory: ""
description: |-
  
---

# AppOptics Provider

## Installation

At the time of writing this documentation, AppOptics provider is not available in the Terraform Registry. You have to:
1. Download provider binary from the [release page on Github](https://github.com/appoptics/terraform-provider-appoptics/releases/tag/v0.5.1)
2. Put the binary in the [local mirror directory](https://www.terraform.io/cli/config/config-file#provider-installation). For example *working-directory*/.terraform/plugins/hashicorp/appoptics/0.5.1/linux_amd64/terraform-provider-appoptics or /usr/local/share/terraform/plugins/hashicorp/appoptics/0.5.1/linux_amd64/terraform-provider-appoptics.

## Authentication

Provider can be authenticated with access token. 
Token can be provided as the environment variable '''APPOPTICS_TOKEN'''
```
export APPOPTICS_TOKEN=..........
```
or with ```provider``` block

```hcl
provider "appoptics" {
  token = var.token
}
```

### Authentication errors

1. Wrong token

Example output:
```
Error: Error reading AppOptics Metric example.metric_one: 401 Unauthorized - {"request":["Authorization Required"]}
```

2. Missing token

If token is missing user will be asked for the input.

```
provider.appoptics.token
  The auth token for the AppOptics account.

  Enter a value:
```

## API Reference

Provider is based on official AppOptics REST API. It's [documentation](https://docs.appoptics.com/api/) can be helpful for using this provider.

## Importing existing resources

The import functionallity is not supported currently. The corresponding issue is [#13](https://github.com/appoptics/terraform-provider-appoptics/issues/13)

## Debugging

For debugging API requests sent by the provider you should set two environment variables:
TF_LOG=debug - generic logging setting for terraform
TF_AO_DEBUG=true - variable specific to this provider
