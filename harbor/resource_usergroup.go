package harbor

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

var groupTypeName2Num = map[string]int64{
	"ldap": 1,
	"http": 2,
	"oidc": 3,
}

var groupTypeNum2Name = map[int64]string{
	1: "ldap",
	2: "http",
	3: "oidc",
}

func resourceUsergroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ldap_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceUsergroupCreate,
		Read:   resourceUsergroupRead,
		Update: resourceUsergroupUpdate,
		Delete: resourceUsergroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func getGroupTypeName(name string) (int64, error) {
	if num, ok := groupTypeName2Num[name]; ok {
		return num, nil
	}
	return 0, fmt.Errorf("group type \"%s\" is unknown", name)
}

func resourceUsergroupCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	usergroupTypeNum, err := getGroupTypeName(d.Get("type").(string))
	if err != nil {
		return err
	}

	body := products.NewPostUsergroupsParams().WithUsergroup(&models.UserGroup{
		GroupType:   usergroupTypeNum,
		LdapGroupDn: d.Get("ldap_dn").(string),
		GroupName:   d.Get("name").(string),
	})

	if _, err := apiClient.Products.PostUsergroups(body, nil); err != nil {
		return err
	}

	usergroup, err := findUsergroupByName(d, m)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(usergroup.ID)))

	return resourceUsergroupRead(d, m)
}

func findUsergroupByName(d *schema.ResourceData, m interface{}) (*models.UserGroup, error) {
	apiClient := m.(*client.Harbor)

	if usergroupName, ok := d.GetOk("name"); ok {
		query := products.NewGetUsergroupsParams()

		resp, err := apiClient.Products.GetUsergroups(query, nil)
		if err != nil {
			d.SetId("")
			return &models.UserGroup{}, err
		}

		for _, usergroup := range resp.Payload {
			if usergroup.GroupName == usergroupName {
				return usergroup, nil
			}
		}

		return &models.UserGroup{}, fmt.Errorf("no usergroups with name %v", usergroupName)
	}

	return &models.UserGroup{}, fmt.Errorf("fail to lookup usergroup by Name")
}

func resourceUsergroupRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if usergroupID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		query := products.NewGetUsergroupsGroupIDParams().WithGroupID(usergroupID)

		resp, err := apiClient.Products.GetUsergroupsGroupID(query, nil)
		if err != nil {
			return err
		}

		if err := setUsergroupSchema(d, resp.Payload); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("fail to load the project")
}

func resourceUsergroupUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	usergroupTypeNum, err := getGroupTypeName(d.Get("type").(string))
	if err != nil {
		return err
	}

	if usergroupID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		body := products.NewPutUsergroupsGroupIDParams().WithUsergroup(&models.UserGroup{
			GroupName:   d.Get("name").(string),
			LdapGroupDn: d.Get("ldap_dn").(string),
			GroupType:   usergroupTypeNum,
		}).WithGroupID(usergroupID)

		if _, err := apiClient.Products.PutUsergroupsGroupID(body, nil); err != nil {
			return err
		}

		return resourceUsergroupRead(d, m)
	}

	return fmt.Errorf("Usergroup Id not a Integer")
}

func resourceUsergroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if usergroupID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		delete := products.NewDeleteUsergroupsGroupIDParams().WithGroupID(usergroupID)
		if _, err := apiClient.Products.DeleteUsergroupsGroupID(delete, nil); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Usergroup Id not a Integer")
}

func setUsergroupSchema(data *schema.ResourceData, usergroup *models.UserGroup) error {
	data.SetId(strconv.Itoa(int(usergroup.ID)))

	if err := data.Set("name", usergroup.GroupName); err != nil {
		return err
	}

	if err := data.Set("ldap_dn", usergroup.LdapGroupDn); err != nil {
		return err
	}

	usergroupTypeName := groupTypeNum2Name[usergroup.GroupType]
	if err := data.Set("type", usergroupTypeName); err != nil {
		return err
	}

	return nil
}
