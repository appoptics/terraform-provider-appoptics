package appoptics

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppOpticsSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppOpticsSpaceCreate,
		Read:   resourceAppOpticsSpaceRead,
		Update: resourceAppOpticsSpaceUpdate,
		Delete: resourceAppOpticsSpaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceAppOpticsSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	name := d.Get("name").(string)

	space, err := client.SpacesService().Create(name)
	if err != nil {
		return fmt.Errorf("Error creating AppOptics space %s: %s", name, err)
	}

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.SpacesService().Retrieve(space.ID)
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if retryErr != nil {
		return retryErr
	}

	d.SetId(strconv.Itoa(space.ID))
	return resourceAppOpticsSpaceReadResult(d, space)
}

func resourceAppOpticsSpaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)

	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	spaceResp, err := client.SpacesService().Retrieve(int(id))
	if err != nil {
		if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading AppOptics Space %s: %s", d.Id(), err)
	}

	return resourceAppOpticsSpaceReadResult(d, &spaceResp.Space)
}

func resourceAppOpticsSpaceReadResult(d *schema.ResourceData, space *appoptics.Space) error {
	if err := d.Set("name", space.Name); err != nil {
		return err
	}
	return nil
}

func resourceAppOpticsSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)
	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		log.Printf("[INFO] Modifying name space attribute for %d: %#v", id, newName)
		if err = client.SpacesService().Update(int(id), newName); err != nil {
			return err
		}
	}

	return resourceAppOpticsSpaceRead(d, meta)
}

func resourceAppOpticsSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*appoptics.Client)
	id, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Space: %d", id)
	err = client.SpacesService().Delete(int(id))
	if err != nil {
		if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
			log.Printf("Space %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting space: %s", err)
	}

	retryErr := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := client.SpacesService().Retrieve(int(id))
		if err != nil {
			if errResp, ok := err.(*appoptics.ErrorResponse); ok && errResp.Response.StatusCode == 404 {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("space still exists"))
	})

	if retryErr != nil {
		return retryErr
	}

	d.SetId("")
	return nil
}
