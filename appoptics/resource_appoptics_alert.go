package appoptics

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppOpticsAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppOpticsAlertCreate,
		Read:   resourceAppOpticsAlertRead,
		Update: resourceAppOpticsAlertUpdate,
		Delete: resourceAppOpticsAlertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"rearm_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"services": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"condition": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"grouped": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"dynamic": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"detect_reset": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"summary_function": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceAppOpticsAlertConditionsHash,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceAppOpticsAlertConditionsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["metric_name"].(string)))

	tags, present := m["tag"].([]interface{})
	if present && len(tags) > 0 {
		buf.WriteString(fmt.Sprintf("%d-", alertConditionsTagsHash(tags)))
	}

	detectReset, present := m["detect_reset"]
	if present {
		buf.WriteString(fmt.Sprintf("%t-", detectReset.(bool)))
	}

	duration, present := m["duration"]
	if present {
		buf.WriteString(fmt.Sprintf("%d-", duration.(int)))
	}

	threshold, present := m["threshold"]
	if present {
		buf.WriteString(fmt.Sprintf("%f-", threshold.(float64)))
	}

	summaryFunction, present := m["summary_function"]
	if present {
		buf.WriteString(fmt.Sprintf("%s-", summaryFunction.(string)))
	}

	return hashcode.String(buf.String())
}

func alertConditionsTagsHash(tags []interface{}) int {
	var buf bytes.Buffer
	for _, v := range tags {
		m := v.(map[string]interface{})
		buf.WriteString(fmt.Sprintf("%s-", m["name"]))
		buf.WriteString(fmt.Sprintf("%s-", m["grouped"]))
		buf.WriteString(fmt.Sprintf("%d-", alertConditionsTagsValuesHash(m["values"].([]interface{}))))
	}

	return hashcode.String(buf.String())
}

func alertConditionsTagsValuesHash(s []interface{}) int {
	var buf bytes.Buffer
	for _, v := range s {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	return hashcode.String(buf.String())
}

func resourceAppOpticsAlertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	alert := appoptics.AlertRequest{
		Name: d.Get("name").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		alert.Description = v.(string)
	}
	// GetOK returns not OK for false boolean values, use Get
	alert.Active = d.Get("active").(bool)
	if v, ok := d.GetOk("rearm_seconds"); ok {
		alert.RearmSeconds = v.(int)
	}
	if _, ok := d.GetOk("services"); ok {
		vs := d.Get("services").(*schema.Set)
		services := make([]int, vs.Len())
		for i, serviceID := range vs.List() {
			var err error
			services[i], err = strconv.Atoi(serviceID.(string))
			if err != nil {
				return fmt.Errorf("Error: %s", err)
			}
		}
		alert.Services = services
	}
	if v, ok := d.GetOk("condition"); ok {
		vs := v.(*schema.Set)
		conditions := make([]*appoptics.AlertCondition, vs.Len())

		for i, conditionDataM := range vs.List() {
			conditionData := conditionDataM.(map[string]interface{})
			condition := appoptics.AlertCondition{}

			if v, ok := conditionData["type"].(string); ok && v != "" {
				condition.Type = v
			}
			if v, ok := conditionData["threshold"].(float64); ok && !math.IsNaN(v) {
				condition.Threshold = v
			}
			if v, ok := conditionData["metric_name"].(string); ok && v != "" {
				condition.MetricName = v
			}
			if v, ok := conditionData["tag"].([]interface{}); ok {
				tags := make([]*appoptics.Tag, len(v))
				for i, tagData := range v {
					tag := appoptics.Tag{}
					tag.Grouped = tagData.(map[string]interface{})["grouped"].(bool)
					tag.Dynamic = tagData.(map[string]interface{})["dynamic"].(bool)
					tag.Name = tagData.(map[string]interface{})["name"].(string)
					values := tagData.(map[string]interface{})["values"].([]interface{})
					valuesInStrings := make([]string, len(values))
					for i, v := range values {
						valuesInStrings[i] = v.(string)
					}
					tag.Values = valuesInStrings
					tags[i] = &tag
				}

				condition.Tags = tags
			}
			if v, ok := conditionData["duration"].(int); ok {
				condition.Duration = v
			}
			if v, ok := conditionData["summary_function"].(string); ok && v != "" {
				condition.SummaryFunction = v
			}
			conditions[i] = &condition
		}

		alert.Conditions = conditions
	}
	if v, ok := d.GetOk("attributes"); ok {
		attributeData := v.(map[string]interface{})
		if len(attributeData) > 1 {
			return fmt.Errorf("Only one set of attributes per alert is supported")
		} else if len(attributeData) == 1 {
			// The only attribute here should be the runbook_url
			alert.Attributes = attributeData
		}
	}

	alertResult, err := client.AlertsService().Create(&alert)

	if err != nil {
		return fmt.Errorf("Error creating AppOptics alert %s: %s", alert.Name, err)
	}
	log.Printf("[INFO] Created AppOptics alert: %s", alertResult.Name)

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.AlertsService().Retrieve(alertResult.ID)
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		return fmt.Errorf("Error creating AppOptics alert %s: %s", alert.Name, err)
	}

	d.SetId(strconv.FormatUint(uint64(alertResult.ID), 10))

	return resourceAppOpticsAlertRead(d, meta)
}

func resourceAppOpticsAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)
	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading AppOptics Alert: %d", id)
	alert, err := client.AlertsService().Retrieve(int(id))
	if err != nil {
		if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading AppOptics Alert %s: %s", d.Id(), err)
	}
	log.Printf("[INFO] Received AppOptics Alert: %s", alert.Name)

	d.Set("name", alert.Name) //nolint

	if err := d.Set("description", alert.Description); err != nil {
		return err
	}

	if err := d.Set("active", alert.Active); err != nil {
		return err
	}

	if err := d.Set("rearm_seconds", alert.RearmSeconds); err != nil {
		return err
	}

	// Since the following aren't simple terraform types (TypeList), it's best to
	// catch the error returned from the d.Set() function, and handle accordingly.
	services := flattenServices(d, alert.Services)
	// TODO: does this need `schema.NewSet(...)`?
	if err := d.Set("services", schema.NewSet(schema.HashString, services)); err != nil {
		return err
	}

	conditions := flattenCondition(d, alert.Conditions)
	if err := d.Set("condition", conditions); err != nil {
		return err
	}

	if err := d.Set("attributes", alert.Attributes); err != nil {
		return err
	}

	return nil
}

func flattenServices(d *schema.ResourceData, services []*appoptics.Service) []interface{} {
	retServices := make([]interface{}, 0, len(services))

	for _, serviceData := range services {
		retServices = append(retServices, fmt.Sprintf("%.d", serviceData.ID))
	}

	return retServices
}

func flattenCondition(d *schema.ResourceData, conditions []*appoptics.AlertCondition) []interface{} {
	out := make([]interface{}, 0, len(conditions))
	for _, c := range conditions {
		condition := make(map[string]interface{})
		condition["type"] = c.Type
		condition["threshold"] = c.Threshold
		condition["metric_name"] = c.MetricName
		condition["tag"] = flattenConditionTags(c.Tags)
		// TODO: once we upgrade the appoptics-api-go dependency,
		// we need to add a `condition["detect_reset"] = c.DetectReset` below
		// SEE: https://github.com/appoptics/terraform-provider-appoptics/issues/12
		// condition["detect_reset"] = c.DetectReset
		condition["duration"] = int(c.Duration)
		condition["summary_function"] = c.SummaryFunction
		out = append(out, condition)
	}

	return out
}

func flattenConditionTags(in []*appoptics.Tag) []interface{} {
	var out = make([]interface{}, 0, len(in))
	for _, v := range in {
		m := make(map[string]interface{})
		m["name"] = v.Name
		m["grouped"] = v.Grouped
		m["dynamic"] = v.Dynamic
		if len(v.Values) > 0 {
			m["values"] = flattenConditionTagsValues(v.Values)
		}
		out = append(out, m)
	}

	return out
}

func flattenConditionTagsValues(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

func resourceAppOpticsAlertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	theAlert, err := client.AlertsService().Retrieve(int(id))
	if err != nil {
		return err
	}
	alert := alertToAlertRequest(theAlert)
	alert.ID = int(id)

	if d.HasChange("name") {
		alert.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		alert.Description = d.Get("description").(string)
	}
	if d.HasChange("active") {
		alert.Active = d.Get("active").(bool)
	}
	if d.HasChange("rearm_seconds") {
		alert.RearmSeconds = d.Get("rearm_seconds").(int)
	}
	if d.HasChange("services") {
		vs := d.Get("services").(*schema.Set)
		services := make([]int, vs.Len())
		for i, serviceID := range vs.List() {
			var err error
			services[i], err = strconv.Atoi(serviceID.(string))
			if err != nil {
				return fmt.Errorf("Error: %s", err)
			}
		}
		alert.Services = services
	}

	// We always have to send the conditions hash, from the API docs:
	//
	// NOTE: This method requires the conditions hash.
	// If conditions is not included in the payload, the alert conditions will be removed.
	vs := d.Get("condition").(*schema.Set)
	conditions := make([]*appoptics.AlertCondition, vs.Len())

	for i, conditionDataM := range vs.List() {
		conditionData := conditionDataM.(map[string]interface{})
		condition := appoptics.AlertCondition{}

		if v, ok := conditionData["type"].(string); ok && v != "" {
			condition.Type = v
		}
		if v, ok := conditionData["threshold"].(float64); ok && !math.IsNaN(v) {
			condition.Threshold = v
		}
		if v, ok := conditionData["metric_name"].(string); ok && v != "" {
			condition.MetricName = v
		}
		if v, ok := conditionData["tag"].([]interface{}); ok {
			tags := make([]*appoptics.Tag, len(v))
			for i, tagData := range v {
				tag := appoptics.Tag{}
				tag.Grouped = tagData.(map[string]interface{})["grouped"].(bool)
				tag.Dynamic = tagData.(map[string]interface{})["dynamic"].(bool)
				tag.Name = tagData.(map[string]interface{})["name"].(string)
				values := tagData.(map[string]interface{})["values"].([]interface{})
				valuesInStrings := make([]string, len(values))
				for i, v := range values {
					valuesInStrings[i] = v.(string)
				}
				tag.Values = valuesInStrings
				tags[i] = &tag
			}

			condition.Tags = tags
		}
		if v, ok := conditionData["duration"].(int); ok {
			condition.Duration = v
		}
		if v, ok := conditionData["summary_function"].(string); ok && v != "" {
			condition.SummaryFunction = v
		}
		conditions[i] = &condition
	}
	alert.Conditions = conditions

	if d.HasChange("attributes") {
		attributeData := d.Get("attributes").([]interface{})
		if attributeData[0] == nil {
			return fmt.Errorf("No attributes found in attributes block")
		}

		alert.Attributes = attributeData[0].(map[string]interface{})
	}

	log.Printf("[INFO] Updating AppOptics alert: %s", alert.Name)
	updErr := client.AlertsService().Update(alert)
	if updErr != nil {
		return fmt.Errorf("Error updating AppOptics alert: %s", updErr)
	}

	log.Printf("[INFO] Updated AppOptics alert %d", id)

	// Wait for propagation since AppOptics updates are eventually consistent
	wait := resource.StateChangeConf{
		Pending:                   []string{fmt.Sprintf("%t", false)},
		Target:                    []string{fmt.Sprintf("%t", true)},
		Timeout:                   5 * time.Minute,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 5,
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Checking if AppOptics Alert %d was updated yet", id)
			changedAlert, getErr := client.AlertsService().Retrieve(int(id))
			if getErr != nil {
				return changedAlert, "", getErr
			}
			return changedAlert, "true", nil
		},
	}

	_, err = wait.WaitForState()
	if err != nil {
		return fmt.Errorf("Failed updating AppOptics Alert %d: %s", id, err)
	}

	return resourceAppOpticsAlertRead(d, meta)
}

func resourceAppOpticsAlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)
	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Alert: %d", id)
	err = client.AlertsService().Delete(int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Alert: %s", err)
	}

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.AlertsService().Retrieve(int(id))
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("alert still exists"))
	})
	if retryErr != nil {
		return fmt.Errorf("Error deleting AppOptics alert: %s", err)
	}

	return nil
}

// used to deal w/ differing structures in API create/read
func alertToAlertRequest(a *appoptics.Alert) *appoptics.AlertRequest {
	aReq := &appoptics.AlertRequest{}
	aReq.ID = a.ID
	aReq.Name = a.Name
	aReq.Description = a.Description
	aReq.Active = a.Active
	aReq.RearmSeconds = a.RearmSeconds
	aReq.Conditions = a.Conditions
	aReq.Attributes = a.Attributes
	aReq.CreatedAt = a.CreatedAt
	aReq.UpdatedAt = a.UpdatedAt

	serviceIDs := make([]int, len(a.Services))
	for i, service := range a.Services {
		serviceIDs[i] = service.ID
	}
	aReq.Services = serviceIDs
	return aReq
}
