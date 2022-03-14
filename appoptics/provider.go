package appoptics

import (
	"os"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for Librato.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPOPTICS_TOKEN", nil),
				Description: "The auth token for the AppOptics account.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"appoptics_dashboard":            resourceAppOpticsSpace(),      // name is legacy from Librato
			"appoptics_dashboard_chart":      resourceAppOpticsSpaceChart(), // name is legacy from Librato
			"appoptics_metric":               resourceAppOpticsMetric(),
			"appoptics_alert":                resourceAppOpticsAlert(),
			"appoptics_notification_service": resourceAppOpticsService(), // changed from API name to differentiate w/ APM Services
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var url string
	if appOpticsURL := os.Getenv("APPOPTICS_URL"); appOpticsURL != "" {
		url = appOpticsURL
	} else {
		url = "https://api.appoptics.com/v1/"
	}
	if do_http_debug := os.Getenv("TF_AO_DEBUG"); do_http_debug != "" {
		return appoptics.NewClient(d.Get("token").(string),
			appoptics.BaseURLClientOption(url),
			appoptics.SetDebugMode(),
		), nil
	} else {
		return appoptics.NewClient(d.Get("token").(string),
			appoptics.BaseURLClientOption(url)), nil
	}
}
