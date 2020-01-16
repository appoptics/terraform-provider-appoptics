package appoptics

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	appoptics "github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppOpticsMetric() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppOpticsMetricCreate,
		Read:   resourceAppOpticsMetricRead,
		Update: resourceAppOpticsMetricUpdate,
		Delete: resourceAppOpticsMetricDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"composite": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"color": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"summarize_function": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_max": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"display_min": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"display_units_long": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_units_short": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_stacked": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"created_by_ua": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gap_detection": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"aggregate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAppOpticsMetricCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)
	metric := appoptics.Metric{
		Name: d.Get("name").(string),
		Type: "gauge",
	}
	if a, ok := d.GetOk("display_name"); ok {
		metric.DisplayName = a.(string)
	}
	if a, ok := d.GetOk("description"); ok {
		metric.Description = a.(string)
	}
	if a, ok := d.GetOk("period"); ok {
		metric.Period = a.(int)
	}
	if a, ok := d.GetOk("composite"); ok {
		metric.Composite = a.(string)
		metric.Type = "composite"
	}

	if a, ok := d.GetOk("attributes"); ok {

		attributeData := a.([]interface{})
		attributeDataMap := attributeData[0].(map[string]interface{})
		attributes := appoptics.MetricAttributes{}

		if v, ok := attributeDataMap["color"].(string); ok && v != "" {
			attributes.Color = v
		}
		if v, ok := attributeDataMap["display_max"].(float64); ok && v != 0.0 {
			attributes.DisplayMax = v
		}
		if v, ok := attributeDataMap["display_min"].(float64); ok && v != 0.0 {
			attributes.DisplayMin = v
		}
		if v, ok := attributeDataMap["display_units_long"].(string); ok && v != "" {
			attributes.DisplayUnitsLong = v
		}
		if v, ok := attributeDataMap["display_units_short"].(string); ok && v != "" {
			attributes.DisplayUnitsShort = v
		}
		if v, ok := attributeDataMap["created_by_ua"].(string); ok && v != "" {
			attributes.CreatedByUA = v
		}
		if v, ok := attributeDataMap["summarize_function"].(string); ok && v != "" {
			attributes.SummarizeFunction = v
		}
		if v, ok := attributeDataMap["display_stacked"].(bool); ok {
			attributes.DisplayStacked = v
		}
		if v, ok := attributeDataMap["gap_detection"].(bool); ok {
			attributes.GapDetection = v
		}
		if v, ok := attributeDataMap["aggregate"].(bool); ok {
			attributes.Aggregate = v
		}

		metric.Attributes = attributes
	}

	_, err := client.MetricsService().Create(&metric)
	if err != nil {
		log.Printf("[INFO] ERROR creating Metric: %s", err)
		return fmt.Errorf("Error creating AppOptics metric: %s", err)
	}

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.MetricsService().Retrieve(metric.Name)
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		return fmt.Errorf("Error creating AppOptics metric: %s", retryErr)
	}

	d.SetId(metric.Name)
	return resourceAppOpticsMetricRead(d, meta)
}

func resourceAppOpticsMetricRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	id := d.Id()

	log.Printf("[INFO] Reading AppOptics Metric: %s", id)
	metric, err := client.MetricsService().Retrieve(id)
	if err != nil {
		if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading AppOptics Metric %s: %s", id, err)
	}

	d.Set("name", metric.Name) //nolint
	d.Set("type", metric.Type) //nolint

	if metric.Description != "" {
		d.Set("description", metric.Description) //nolint
	}

	if metric.DisplayName != "" {
		d.Set("display_name", metric.DisplayName) //nolint
	}

	if metric.Period != 0 {
		d.Set("period", metric.Period) //nolint
	}

	if metric.Composite != "" {
		d.Set("composite", metric.Composite) //nolint
	}

	attributes := metricAttributesGather(d, &metric.Attributes)

	// Since attributes isn't a simple terraform type (TypeList), it's best to
	// catch the error returned from the d.Set() function, and handle accordingly.
	if err := d.Set("attributes", attributes); err != nil {
		return err
	}

	return nil
}

func resourceAppOpticsMetricUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	id := d.Id()
	metric, err := client.MetricsService().Retrieve(id)

	if err != nil {
		return err
	}

	if d.HasChange("type") {
		metric.Type = d.Get("type").(string)
	}
	if d.HasChange("description") {
		metric.Description = d.Get("description").(string)
	}
	if d.HasChange("display_name") {
		metric.DisplayName = d.Get("display_name").(string)
	}
	if d.HasChange("period") {
		metric.Period = d.Get("period").(int)
	}
	if d.HasChange("composite") {
		metric.Composite = d.Get("composite").(string)
	}
	if d.HasChange("attributes") {
		attributeData := d.Get("attributes").([]interface{})
		attributeDataMap := attributeData[0].(map[string]interface{})
		attributes := appoptics.MetricAttributes{}

		if v, ok := attributeDataMap["color"].(string); ok && v != "" {
			attributes.Color = v
		}
		if v, ok := attributeDataMap["display_max"].(float64); ok && v != 0.0 {
			attributes.DisplayMax = v
		}
		if v, ok := attributeDataMap["display_min"].(float64); ok && v != 0.0 {
			attributes.DisplayMin = v
		}
		if v, ok := attributeDataMap["display_units_long"].(string); ok && v != "" {
			attributes.DisplayUnitsLong = v
		}
		if v, ok := attributeDataMap["summarize_function"].(string); ok && v != "" {
			attributes.SummarizeFunction = v
		}
		if v, ok := attributeDataMap["display_units_short"].(string); ok && v != "" {
			attributes.DisplayUnitsShort = v
		}
		if v, ok := attributeDataMap["created_by_ua"].(string); ok && v != "" {
			attributes.CreatedByUA = v
		}
		if v, ok := attributeDataMap["display_stacked"].(bool); ok {
			attributes.DisplayStacked = v
		}
		if v, ok := attributeDataMap["gap_detection"].(bool); ok {
			attributes.GapDetection = v
		}
		if v, ok := attributeDataMap["aggregate"].(bool); ok {
			attributes.Aggregate = v
		}
		metric.Attributes = attributes
	}

	log.Printf("[INFO] Updating AppOptics metric: %v", structToString(metric))

	err = client.MetricsService().Update(id, metric)
	if err != nil {
		return fmt.Errorf("Error updating AppOptics metric: %s", err)
	}

	log.Printf("[INFO] Updated AppOptics metric %s", id)

	// Wait for propagation since AppOptics updates are eventually consistent
	wait := resource.StateChangeConf{
		Pending:                   []string{fmt.Sprintf("%t", false)},
		Target:                    []string{fmt.Sprintf("%t", true)},
		Timeout:                   5 * time.Minute,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 5,
		Refresh: func() (interface{}, string, error) {
			log.Printf("[INFO] Checking if AppOptics Metric %s was updated yet", id)
			changedMetric, err := client.MetricsService().Retrieve(id)
			if err != nil {
				return changedMetric, "", err
			}
			return changedMetric, "true", nil
		},
	}

	_, err = wait.WaitForState()
	if err != nil {
		log.Printf("[INFO] ERROR - Failed updating AppOptics Metric %s: %s", id, err)
		return fmt.Errorf("Failed updating AppOptics Metric %s: %s", id, err)
	}

	return resourceAppOpticsMetricRead(d, meta)
}

func resourceAppOpticsMetricDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	id := d.Id()

	log.Printf("[INFO] Deleting Metric: %s", id)
	err := client.MetricsService().Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting Metric: %s", err)
	}

	log.Printf("[INFO] Verifying Metric %s deleted", id)
	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {

		log.Printf("[INFO] Getting Metric %s", id)
		_, err := client.MetricsService().Retrieve(id)
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				log.Printf("[INFO] Metric %s not found, removing from state", id)
				return nil
			}
			log.Printf("[INFO] non-retryable error attempting to Get metric: %s", err)
			return resource.NonRetryableError(err)
		}

		log.Printf("[INFO] retryable error attempting to Get metric: %s", id)
		return resource.RetryableError(fmt.Errorf("metric still exists"))
	})
	if retryErr != nil {
		return fmt.Errorf("Error deleting AppOptics metric: %s", retryErr)
	}

	return nil
}

// Flattens an attributes hash into something that flatmap.Flatten() can handle
func metricAttributesGather(d *schema.ResourceData, attributes *appoptics.MetricAttributes) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if attributes != nil {
		retAttributes := make(map[string]interface{})
		if attributes.Color != "" {
			retAttributes["color"] = attributes.Color
		}
		if attributes.DisplayMax != nil {
			retAttributes["display_max"] = attributes.DisplayMax
		}
		if attributes.DisplayMin != nil {
			retAttributes["display_min"] = attributes.DisplayMin
		}
		if attributes.DisplayUnitsLong != "" {
			retAttributes["display_units_long"] = attributes.DisplayUnitsLong
		}
		if attributes.SummarizeFunction != "" {
			retAttributes["summarize_function"] = attributes.SummarizeFunction
		}
		if attributes.DisplayUnitsShort != "" {
			retAttributes["display_units_short"] = attributes.DisplayUnitsShort
		}
		if attributes.CreatedByUA != "" {
			retAttributes["created_by_ua"] = attributes.CreatedByUA
		}
		retAttributes["display_stacked"] = attributes.DisplayStacked || false
		retAttributes["gap_detection"] = attributes.GapDetection || false
		retAttributes["aggregate"] = attributes.Aggregate || false

		result = append(result, retAttributes)
	}

	return result
}

func structToString(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}
