package appoptics

import (
	"os"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
			"appoptics_dashboard":       resourceAppOpticsSpace(),      // name is legacy from Librato
			"appoptics_dashboard_chart": resourceAppOpticsSpaceChart(), // name is legacy from Librato
			"appoptics_metric":          resourceAppOpticsMetric(),
			"appoptics_alert":           resourceAppOpticsAlert(),
			"appoptics_service":         resourceAppOpticsService(),
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

	client := appoptics.NewClient(d.Get("token").(string),
		appoptics.BaseURLClientOption(url),
		appoptics.SetDebugMode(),
	)
	return client, nil
}
