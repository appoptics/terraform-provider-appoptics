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
* Grab the latest release from the [Releases page](https://github.com/appoptics/terraform-provider-appoptics/releases).
* Run `make build`.
* Place the binary (`terraform-provider-appoptics`) at `~/.terraform.d/plugins/solarwinds.com/appopticsprovider/appoptics/1.0.0/darwin_arm64/terraform-provider-appoptics` (Replace `darwin_arm64` by `HOST_ARCH`) [where Terraform can find it](https://www.terraform.io/language/providers/requirements)
* You should now be able to write TF code for AppOptics alongside the rest of your infrastructure code

### Issues/Bugs
Please report bugs and request enhancements in the [Issues area](https://github.com/appoptics/terraform-provider-appoptics/issues) of this repo.
