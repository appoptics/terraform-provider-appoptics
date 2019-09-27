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
				Default:  math.NaN(),
				Optional: true,
			},
			"max": {
				Type:     schema.TypeFloat,
				Default:  math.NaN(),
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
						"source": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"stream.composite"},
						},
						"group_function": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"stream.composite"},
						},
						"composite": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"stream.metric", "stream.source", "stream.group_function"},
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
							Default:  math.NaN(),
							Optional: true,
						},
						"max": {
							Type:     schema.TypeFloat,
							Default:  math.NaN(),
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
	buf.WriteString(fmt.Sprintf("%s-", m["source"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["composite"].(string)))

	return hashcode.String(buf.String())
}

func resourceAppOpticsSpaceChartCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	spaceID := uint(d.Get("space_id").(int))

	spaceChart := new(appoptics.Chart)
	if v, ok := d.GetOk("name"); ok {
		spaceChart.Name = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		spaceChart.Type = v.(string)
	}
	if v, ok := d.GetOk("min"); ok {
		if math.IsNaN(v.(float64)) {
			spaceChart.Min = nil
		} else {
			spaceChart.Min = v.(float64)
		}
	}
	if v, ok := d.GetOk("max"); ok {
		if math.IsNaN(v.(float64)) {
			spaceChart.Max = nil
		} else {
			spaceChart.Max = v.(float64)
		}
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
			if v, ok := streamData["source"].(string); ok && v != "" {
				stream.Source = v
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
				stream.Min = v
			}
			if v, ok := streamData["max"].(float64); ok && !math.IsNaN(v) {
				stream.Max = v
			}
			streams[i] = stream
		}
		spaceChart.Streams = streams
	}

	spaceChartResult, _, err := client.Spaces.CreateChart(spaceID, spaceChart)
	if err != nil {
		return fmt.Errorf("Error creating AppOptics space chart %s: %s", *spaceChart.Name, err)
	}

	resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, _, err := client.Spaces.GetChart(spaceID, *spaceChartResult.ID)
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

	spaceID := uint(d.Get("space_id").(int))

	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	chart, _, err := client.Spaces.GetChart(spaceID, uint(id))
	if err != nil {
		if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading AppOptics Space chart %s: %s", d.Id(), err)
	}

	return resourceAppOpticsSpaceChartReadResult(d, chart)
}

func resourceAppOpticsSpaceChartReadResult(d *schema.ResourceData, chart *appoptics.SpaceChart) error {
	d.SetId(strconv.FormatUint(uint64(*chart.ID), 10))
	if chart.Name != nil {
		if err := d.Set("name", *chart.Name); err != nil {
			return err
		}
	}
	if chart.Type != nil {
		if err := d.Set("type", *chart.Type); err != nil {
			return err
		}
	}
	if chart.Min != nil {
		if err := d.Set("min", *chart.Min); err != nil {
			return err
		}
	}
	if chart.Max != nil {
		if err := d.Set("max", *chart.Max); err != nil {
			return err
		}
	}
	if chart.Label != nil {
		if err := d.Set("label", *chart.Label); err != nil {
			return err
		}
	}
	if chart.RelatedSpace != nil {
		if err := d.Set("related_space", *chart.RelatedSpace); err != nil {
			return err
		}
	}

	streams := resourceAppOpticsSpaceChartStreamsGather(d, chart.Streams)
	if err := d.Set("stream", streams); err != nil {
		return err
	}

	return nil
}

func resourceAppOpticsSpaceChartStreamsGather(d *schema.ResourceData, streams []appoptics.SpaceChartStream) []map[string]interface{} {
	retStreams := make([]map[string]interface{}, 0, len(streams))
	for _, s := range streams {
		stream := make(map[string]interface{})
		if s.Metric != nil {
			stream["metric"] = *s.Metric
		}
		if s.Source != nil {
			stream["source"] = *s.Source
		}
		if s.Composite != nil {
			stream["composite"] = *s.Composite
		}
		if s.GroupFunction != nil {
			stream["group_function"] = *s.GroupFunction
		}
		if s.SummaryFunction != nil {
			stream["summary_function"] = *s.SummaryFunction
		}
		if s.TransformFunction != nil {
			stream["transform_function"] = *s.TransformFunction
		}
		if s.Color != nil {
			stream["color"] = *s.Color
		}
		if s.UnitsShort != nil {
			stream["units_short"] = *s.UnitsShort
		}
		if s.UnitsLong != nil {
			stream["units_long"] = *s.UnitsLong
		}
		retStreams = append(retStreams, stream)
	}

	return retStreams
}

func resourceAppOpticsSpaceChartUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	spaceID := uint(d.Get("space_id").(int))
	chartID, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	// Just to have whole object for comparison before/after update
	fullChart, _, err := client.Spaces.GetChart(spaceID, uint(chartID))
	if err != nil {
		return err
	}

	spaceChart := new(appoptics.SpaceChart)
	if d.HasChange("name") {
		spaceChart.Name = d.Get("name").(string)
		fullChart.Name = spaceChart.Name
	}
	if d.HasChange("min") {
		if math.IsNaN(d.Get("min").(float64)) {
			spaceChart.Min = nil
		} else {
			spaceChart.Min = d.Get("min").(float64)
		}
		fullChart.Min = spaceChart.Min
	}
	if d.HasChange("max") {
		if math.IsNaN(d.Get("max").(float64)) {
			spaceChart.Max = nil
		} else {
			spaceChart.Max = d.Get("max").(float64)
		}
		fullChart.Max = spaceChart.Max
	}
	if d.HasChange("label") {
		spaceChart.Label = d.Get("label").(string)
		fullChart.Label = spaceChart.Label
	}
	if d.HasChange("related_space") {
		spaceChart.RelatedSpace = d.Get("related_space").(uint)
		fullChart.RelatedSpace = spaceChart.RelatedSpace
	}
	if d.HasChange("stream") {
		vs := d.Get("stream").(*schema.Set)
		streams := make([]appoptics.SpaceChartStream, vs.Len())
		for i, streamDataM := range vs.List() {
			streamData := streamDataM.(map[string]interface{})
			var stream appoptics.SpaceChartStream
			if v, ok := streamData["metric"].(string); ok && v != "" {
				stream.Metric = v
			}
			if v, ok := streamData["source"].(string); ok && v != "" {
				stream.Source = v
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
				stream.Min = v
			}
			if v, ok := streamData["max"].(float64); ok && !math.IsNaN(v) {
				stream.Max = v
			}
			streams[i] = stream
		}
		spaceChart.Streams = streams
		fullChart.Streams = streams
	}

	_, err = client.Spaces.UpdateChart(spaceID, uint(chartID), spaceChart)
	if err != nil {
		return fmt.Errorf("Error updating AppOptics space chart %s: %s", *spaceChart.Name, err)
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
			changedChart, _, getErr := client.Spaces.GetChart(spaceID, uint(chartID))
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

	spaceID := uint(d.Get("space_id").(int))

	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Chart: %d/%d", spaceID, uint(id))
	_, err = client.Spaces.DeleteChart(spaceID, uint(id))
	if err != nil {
		return fmt.Errorf("Error deleting space: %s", err)
	}

	resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, _, err := client.Spaces.GetChart(spaceID, uint(id))
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
