## terraform-provider-appoptics
This provider lets you save clicking in the AO UI by allowing you to produce AppOptics bits alongside the rest of your cloud infratructure.

You're able to programmatically create:

* Dashboards
* Charts
* Metrics
* Alerts
* Notification Services

### Example usage
See `example.tf` [in this repo](https://github.com/appoptics/terraform-provider-appoptics/blob/master/examples/) to understand how to start using the plugin.

### Installing
* Grab the latest release binary from the [Releases page](https://github.com/appoptics/terraform-provider-appoptics/releases).
* Extract and place the binary into `$HOME/.terraform.d/plugins/solarwinds.com/appopticsprovider/appoptics/<VERSION>/<ARCH>/terraform-provider-appoptics` (Replace `<VERSION>` with the version downloaded and `<ARCH>` with the machine architecture (eg. `darwin_amd64` or `darwin_arm64`)
* Set the execute flag on the binary
```
chmod 755 $HOME/.terraform.d/plugins/solarwinds.com/appopticsprovider/appoptics/<VERSION>/<ARCH>/terraform-provider-appoptics
```
* You should now be able to write TF code for AppOptics alongside the rest of your infrastructure code

### Usage Notes
In order for the provider to work in a module, you need to add a required_providers block in your module as such:
```hcl
terraform {
  required_providers {
    appoptics = {
      source  = "solarwinds.com/appopticsprovider/appoptics"
      version = ">= 0.5.1"
    }
  }
}
```
This needs to be done because this provider has not been published to the Terraform registry, which is the default location that Terraform will look in when searching for providers.

### Issues/Bugs
Please report bugs and request enhancements in the [Issues area](https://github.com/appoptics/terraform-provider-appoptics/issues) of this repo.
