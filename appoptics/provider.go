package appoptics

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/henrikhodne/go-appoptics/librato"
	"net/url"
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
			"appoptics_space":       resourceAppOpticsSpace(),
			"appoptics_space_chart": resourceAppOpticsSpaceChart(),
			"appoptics_metric":      resourceAppOpticsMetric(),
			"appoptics_alert":       resourceAppOpticsAlert(),
			"appoptics_service":     resourceAppOpticsService(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url, _ := url.Parse("  https://api.appoptics.com/v1/measurements/v1")
	client := appoptics.NewClientWithBaseURL(url, d.Get("email").(string), d.Get("token").(string))
	return client, nil
}
