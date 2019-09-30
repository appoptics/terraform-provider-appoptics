package appoptics

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppOpticsSpaceChart() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppOpticsSpaceChartCreate,
		Read:   resourceAppOpticsSpaceChartRead,
		Update: resourceAppOpticsSpaceChartUpdate,
		Delete: resourceAppOpticsSpaceChartDelete,

		Schema: map[string]*schema.Schema{
			"space_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"min": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"max": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"related_space": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"stream": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"stream.composite"},
						},
						"tag": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"stream.composite"},
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
						"group_function": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"stream.composite"},
						},
						"composite": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"stream.metric", "stream.group_function"},
						},
						"summary_function": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"color": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"units_short": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"units_long": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"max": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"transform_function": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				Set: resourceAppOpticsSpaceChartHash,
			},
		},
	}
}

func resourceAppOpticsSpaceChartHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["metric"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["composite"].(string)))
	tags, present := m["tag"].([]interface{})
	if present && len(tags) > 0 {
		buf.WriteString(fmt.Sprintf("%d-", chartStreamTagsHash(tags)))
	}

	return hashcode.String(buf.String())
}

func chartStreamTagsHash(tags []interface{}) int {
	var buf bytes.Buffer
	for _, v := range tags {
		m := v.(map[string]interface{})
		buf.WriteString(fmt.Sprintf("%s-", m["name"]))
		buf.WriteString(fmt.Sprintf("%s-", m["grouped"]))
		buf.WriteString(fmt.Sprintf("%d-", chartStreamTagsValuesHash(m["values"].([]interface{}))))
	}

	return hashcode.String(buf.String())
}

func chartStreamTagsValuesHash(s []interface{}) int {
	var buf bytes.Buffer
	for _, v := range s {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	return hashcode.String(buf.String())
}

func resourceAppOpticsSpaceChartCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	spaceID := d.Get("space_id").(int)

	spaceChart := new(appoptics.Chart)
	if v, ok := d.GetOk("name"); ok {
		spaceChart.Name = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		spaceChart.Type = v.(string)
	}
	if v, ok := d.GetOk("min"); ok {
		if math.IsNaN(v.(float64)) {
			return fmt.Errorf("Error creating AppOptics space chart. 'min' cannot be converted to a float64. %s", d.Get("min"))
		}
		spaceChart.Min = v.(float64)
	}
	if v, ok := d.GetOk("max"); ok {
		if math.IsNaN(v.(float64)) {
			return fmt.Errorf("Error creating AppOptics space chart. 'max' cannot be converted to a float64. %s", d.Get("max"))
		}
		spaceChart.Max = v.(float64)
	}
	if v, ok := d.GetOk("label"); ok {
		spaceChart.Label = v.(string)
	}
	if v, ok := d.GetOk("related_space"); ok {
		spaceChart.RelatedSpace = v.(int)
	}
	if v, ok := d.GetOk("stream"); ok {
		vs := v.(*schema.Set)
		streams := make([]appoptics.Stream, vs.Len())
		for i, streamDataM := range vs.List() {
			streamData := streamDataM.(map[string]interface{})
			var stream appoptics.Stream
			if v, ok := streamData["metric"].(string); ok && v != "" {
				stream.Metric = v
			}
			if v, ok := streamData["tags"].([]appoptics.Tag); ok {
				stream.Tags = v
			}
			if v, ok := streamData["composite"].(string); ok && v != "" {
				stream.Composite = v
			}
			if v, ok := streamData["group_function"].(string); ok && v != "" {
				stream.GroupFunction = v
			}
			if v, ok := streamData["summary_function"].(string); ok && v != "" {
				stream.SummaryFunction = v
			}
			if v, ok := streamData["transform_function"].(string); ok && v != "" {
				stream.TransformFunction = v
			}
			if v, ok := streamData["color"].(string); ok && v != "" {
				stream.Color = v
			}
			if v, ok := streamData["units_short"].(string); ok && v != "" {
				stream.UnitsShort = v
			}
			if v, ok := streamData["units_longs"].(string); ok && v != "" {
				stream.UnitsLong = v
			}
			if v, ok := streamData["min"].(float64); ok && !math.IsNaN(v) {
				stream.Min = int(v)
			}
			if v, ok := streamData["max"].(float64); ok && !math.IsNaN(v) {
				stream.Max = int(v)
			}
			streams[i] = stream
		}
		spaceChart.Streams = streams
	}

	spaceChartResult, err := client.ChartsService().Create(spaceChart, spaceID)
	if err != nil {
		return fmt.Errorf("Error creating AppOptics space chart %s: %s", spaceChart.Name, err)
	}

	resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.ChartsService().Retrieve(spaceChartResult.ID, spaceID)
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return resourceAppOpticsSpaceChartReadResult(d, spaceChartResult)
}

func resourceAppOpticsSpaceChartRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	spaceID := d.Get("space_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	chart, err := client.ChartsService().Retrieve(id, spaceID)
	if err != nil {
		if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading AppOptics Space chart %s: %s", d.Id(), err)
	}

	return resourceAppOpticsSpaceChartReadResult(d, chart)
}

func resourceAppOpticsSpaceChartReadResult(d *schema.ResourceData, chart *appoptics.Chart) error {
	d.SetId(strconv.FormatUint(uint64(chart.ID), 10))
	if err := d.Set("name", chart.Name); err != nil {
		return err
	}
	if err := d.Set("type", chart.Type); err != nil {
		return err
	}
	if err := d.Set("min", chart.Min); err != nil {
		return err
	}
	if err := d.Set("max", chart.Max); err != nil {
		return err
	}
	if err := d.Set("label", chart.Label); err != nil {
		return err
	}
	if err := d.Set("related_space", chart.RelatedSpace); err != nil {
		return err
	}

	streams := resourceAppOpticsSpaceChartStreamsGather(d, chart.Streams)
	if err := d.Set("stream", streams); err != nil {
		return err
	}

	return nil
}

func resourceAppOpticsSpaceChartStreamsGather(d *schema.ResourceData, streams []appoptics.Stream) []map[string]interface{} {
	retStreams := make([]map[string]interface{}, 0, len(streams))
	for _, s := range streams {
		stream := make(map[string]interface{})
		// TODO: support all options in appoptics.Chart
		stream["metric"] = s.Metric
		stream["tags"] = s.Tags
		stream["composite"] = s.Composite
		stream["group_function"] = s.GroupFunction
		stream["summary_function"] = s.SummaryFunction
		stream["transform_function"] = s.TransformFunction
		stream["color"] = s.Color
		stream["units_short"] = s.UnitsShort
		stream["units_long"] = s.UnitsLong
		retStreams = append(retStreams, stream)
	}

	return retStreams
}

func resourceAppOpticsSpaceChartUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	spaceID := d.Get("space_id").(int)
	chartID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	// Just to have whole object for comparison before/after update
	fullChart, err := client.ChartsService().Retrieve(chartID, spaceID)
	if err != nil {
		return err
	}

	spaceChart := new(appoptics.Chart)
	if d.HasChange("name") {
		spaceChart.Name = d.Get("name").(string)
		fullChart.Name = spaceChart.Name
	}
	if d.HasChange("min") {
		if math.IsNaN(d.Get("min").(float64)) {
			return fmt.Errorf("Error updating AppOptics space chart. 'min' cannot be converted to a float64. %s: %s", d.Get("min"), err)
		}
		spaceChart.Min = d.Get("min").(float64)
		fullChart.Min = spaceChart.Min
	}
	if d.HasChange("max") {
		if math.IsNaN(d.Get("max").(float64)) {
			return fmt.Errorf("Error updating AppOptics space chart. 'max' cannot be converted to a float64. %s: %s", d.Get("max"), err)
		}
		spaceChart.Max = d.Get("max").(float64)
		fullChart.Max = spaceChart.Max
	}
	if d.HasChange("label") {
		spaceChart.Label = d.Get("label").(string)
		fullChart.Label = spaceChart.Label
	}
	if d.HasChange("related_space") {
		spaceChart.RelatedSpace = d.Get("related_space").(int)
		fullChart.RelatedSpace = spaceChart.RelatedSpace
	}
	if d.HasChange("stream") {
		vs := d.Get("stream").(*schema.Set)
		streams := make([]appoptics.Stream, vs.Len())
		for i, streamDataM := range vs.List() {
			streamData := streamDataM.(map[string]interface{})
			var stream appoptics.Stream
			if v, ok := streamData["metric"].(string); ok && v != "" {
				stream.Metric = v
			}
			if v, ok := streamData["tags"].([]appoptics.Tag); ok {
				stream.Tags = v
			}
			if v, ok := streamData["composite"].(string); ok && v != "" {
				stream.Composite = v
			}
			if v, ok := streamData["group_function"].(string); ok && v != "" {
				stream.GroupFunction = v
			}
			if v, ok := streamData["summary_function"].(string); ok && v != "" {
				stream.SummaryFunction = v
			}
			if v, ok := streamData["transform_function"].(string); ok && v != "" {
				stream.TransformFunction = v
			}
			if v, ok := streamData["color"].(string); ok && v != "" {
				stream.Color = v
			}
			if v, ok := streamData["units_short"].(string); ok && v != "" {
				stream.UnitsShort = v
			}
			if v, ok := streamData["units_longs"].(string); ok && v != "" {
				stream.UnitsLong = v
			}
			if v, ok := streamData["min"].(int); ok && !math.IsNaN(float64(v)) {
				stream.Min = v
			}
			if v, ok := streamData["max"].(int); ok && !math.IsNaN(float64(v)) {
				stream.Max = v
			}
			streams[i] = stream
		}
		spaceChart.Streams = streams
		fullChart.Streams = streams
	}

	_, err = client.ChartsService().Update(spaceChart, spaceID)
	if err != nil {
		return fmt.Errorf("Error updating AppOptics space chart %s: %s", spaceChart.Name, err)
	}

	// Wait for propagation since AppOptics updates are eventually consistent
	wait := resource.StateChangeConf{
		Pending:                   []string{fmt.Sprintf("%t", false)},
		Target:                    []string{fmt.Sprintf("%t", true)},
		Timeout:                   5 * time.Minute,
		MinTimeout:                2 * time.Second,
		ContinuousTargetOccurence: 5,
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Checking if AppOptics Space Chart %d was updated yet", chartID)
			changedChart, getErr := client.ChartsService().Retrieve(chartID, spaceID)
			if getErr != nil {
				return changedChart, "", getErr
			}
			isEqual := reflect.DeepEqual(*fullChart, *changedChart)
			log.Printf("[DEBUG] Updated AppOptics Space Chart %d match: %t", chartID, isEqual)
			return changedChart, fmt.Sprintf("%t", isEqual), nil
		},
	}

	_, err = wait.WaitForState()
	if err != nil {
		return fmt.Errorf("Failed updating AppOptics Space Chart %d: %s", chartID, err)
	}

	return resourceAppOpticsSpaceChartRead(d, meta)
}

func resourceAppOpticsSpaceChartDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	spaceID := d.Get("space_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Chart: %d/%d", spaceID, uint(id))
	err = client.ChartsService().Delete(id, spaceID)
	if err != nil {
		return fmt.Errorf("Error deleting space: %s", err)
	}

	resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.ChartsService().Retrieve(id, spaceID)
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("space chart still exists"))
	})

	d.SetId("")
	return nil
}
