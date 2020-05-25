package harbor

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceTasks() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vulnerability_scan_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTasksCreate,
		Read:   resourceTasksRead,
		Update: resourceTasksUpdate,
		Delete: resourceTasksDelete,
	}
}

func resourceTasksCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	schedule, err := getSchedule(d.Get("vulnerability_scan_policy").(string))
	if err != nil {
		return err
	}

	body := &models.AdminJobSchedule{
		Schedule: schedule,
	}

	resp, err := apiClient.Products.GetSystemScanAllSchedule(products.NewGetSystemScanAllScheduleParams(), nil)

	if err != nil {
		log.Fatalf("Fail to load vulnerability_scan %v", err)
	}

	time := resp.Payload.Schedule.Type
	if time != "" {
		log.Printf("Shedule found performing PUT request")

		params := products.NewPutSystemScanAllScheduleParams().WithSchedule(body)

		_, err = apiClient.Products.PutSystemScanAllSchedule(params, nil)

		if err != nil {
			log.Fatalf("Fail to update vulnerability_scan %v", err)
		}
	} else {
		log.Printf("No shedule found performing POST request")

		params := products.NewPostSystemScanAllScheduleParams().WithSchedule(body)

		_, err = apiClient.Products.PostSystemScanAllSchedule(params, nil)

		if err != nil {
			log.Fatalf("Fail to create new vulnerability_scan %v", err)
		}
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", "vulnerability_scan")))
	return resourceTasksRead(d, m)
}

func resourceTasksRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	resp, err := apiClient.Products.GetSystemScanAllSchedule(products.NewGetSystemScanAllScheduleParams(), nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := d.Set("vulnerability_scan_policy", strings.ToLower(resp.Payload.Schedule.Type)); err != nil {
		return err
	}

	return nil
}

func resourceTasksUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	schedule, err := getSchedule(d.Get("vulnerability_scan_policy").(string))
	if err != nil {
		return err
	}

	body := &models.AdminJobSchedule{
		Schedule: schedule,
	}

	query := products.NewPutSystemScanAllScheduleParams().WithSchedule(body)

	_, err = apiClient.Products.PutSystemScanAllSchedule(query, nil)
	if err != nil {
		return err
	}

	return resourceTasksRead(d, m)
}

func resourceTasksDelete(d *schema.ResourceData, m interface{}) error {
	// https://github.com/goharbor/harbor/issues/11083
	log.Printf("Not inplemented at the moment look gh issue: %s", "https://github.com/goharbor/harbor/issues/11083")
	return nil
}

func getSchedule(schedule string) (*models.AdminJobScheduleObj, error) {
	switch schedule {
	case "hourly":
		return &models.AdminJobScheduleObj{
			Cron: "0 0 * * * *",
			Type: "Hourly",
		}, nil
	case "daily":
		return &models.AdminJobScheduleObj{
			Cron: "0 0 0 * * *",
			Type: "Daily",
		}, nil
	case "weekly":
		return &models.AdminJobScheduleObj{
			Cron: "0 0 0 * * 0",
			Type: "Weekly",
		}, nil
	}
	return &models.AdminJobScheduleObj{}, errors.New("Not a Valid schedule name %s" + schedule)
}
