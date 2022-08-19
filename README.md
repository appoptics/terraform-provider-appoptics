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
* Extract and place the binary into `$HOME/.terraform.d/plugins/registry.terraform.io/hashicorp/appoptics/<VERSION>/<ARCH>/terraform-provider-appoptics` (Replace `<VERSION>` with the version downloaded and `<ARCH>` with the machine architecture (eg. `darwin_amd64` or `darwin_arm64`)
* Set the execute flag on the binary
```
chmod 755 $HOME/.terraform.d/plugins/registry.terraform.io/hashicorp/appoptics/<VERSION>/<ARCH>/terraform-provider-appoptics
```
* You should now be able to write TF code for AppOptics alongside the rest of your infrastructure code

### Issues/Bugs
Please report bugs and request enhancements in the [Issues area](https://github.com/appoptics/terraform-provider-appoptics/issues) of this repo.
