package harbor

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLabel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deleted": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
		Read: dataSourceLabelRead,
	}
}

func dataSourceLabelRead(d *schema.ResourceData, m interface{}) error {
	if _, ok := d.GetOk("name"); ok {
		if _, ok := d.GetOk("scope"); ok {
			registry, err := findLabelByNameAndScope(d, m)
			if err != nil {
				return err
			}

			if err := setLabelSchema(d, registry); err != nil {
				return err
			}

			return nil
		}
	}

	if labelID, ok := d.GetOk("id"); ok {
		d.SetId(strconv.Itoa(labelID.(int)))

		label, err := findLabelByID(d, m)
		if err != nil {
			return err
		}

		if err := setLabelSchema(d, label); err != nil {
			return err
		}

		return nil
	}

	return errors.New("please specify a combination of name and scope or Id to lookup for a label")
}
