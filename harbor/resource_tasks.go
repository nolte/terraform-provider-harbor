package harbor

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

var TypeStr string
var CronStr string

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

	vulnSchedule := d.Get("vulnerability_scan_policy").(string)
	getSchedule(vulnSchedule)

	body := &models.AdminJobSchedule{
		Schedule: &models.AdminJobScheduleObj{
			Cron: CronStr,
			Type: TypeStr,
		}}

	resp, err := apiClient.Products.GetSystemScanAllSchedule(products.NewGetSystemScanAllScheduleParams(), nil)
	if err != nil {
		log.Fatalf("Fail to load vulnerability_scan %v", err)
	}

	time := resp.Payload.Schedule.Type
	if time != "" {
		log.Printf("Shedule found performing PUT request")
		_, err = apiClient.Products.PutSystemScanAllSchedule(products.NewPutSystemScanAllScheduleParams().WithSchedule(body), nil)
		if err != nil {
			log.Fatalf("Fail to update vulnerability_scan %v", err)
		}
	} else {
		log.Printf("No shedule found performing POST request")
		_, err = apiClient.Products.PostSystemScanAllSchedule(products.NewPostSystemScanAllScheduleParams().WithSchedule(body), nil)
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

	vulnSchedule := d.Get("vulnerability_scan_policy").(string)
	getSchedule(vulnSchedule)

	body := &models.AdminJobSchedule{
		Schedule: &models.AdminJobScheduleObj{
			Cron: CronStr,
			Type: TypeStr,
		}}

	_, err := apiClient.Products.PutSystemScanAllSchedule(products.NewPutSystemScanAllScheduleParams().WithSchedule(body), nil)
	if err != nil {
		log.Fatal(err)
	}

	return resourceTasksRead(d, m)
}

func resourceTasksDelete(d *schema.ResourceData, m interface{}) error {
	// https://github.com/goharbor/harbor/issues/11083
	//	apiClient := m.(*client.Harbor)
	//
	//	body := &models.AdminJobSchedule{
	//		Schedule: &models.AdminJobScheduleObj{
	//			Cron: "",
	//		}}
	//
	//	_, err := apiClient.Products.PutSystemScanAllSchedule(products.NewPutSystemScanAllScheduleParams().WithSchedule(body), nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	log.Printf("Not inplemented at the moment look gh issue: %s", "https://github.com/goharbor/harbor/issues/11083")
	return nil
}

func getSchedule(schedule string) {
	switch schedule {
	case "hourly":
		TypeStr = "Hourly"
		CronStr = "0 0 * * * *"
	case "daily":
		TypeStr = "Daily"
		CronStr = "0 0 0 * * *"
	case "weekly":
		TypeStr = "Weekly"
		CronStr = "0 0 0 * * 0"
	}
}
