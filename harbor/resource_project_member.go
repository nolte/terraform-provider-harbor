package harbor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

var roleName2Id = map[string]int64{
	"project_admin": 1,
	"master":        2,
	"developer":     3,
	"guest":         4,
	"limited_guest": 5,
}

var roleId2Name = map[int64]string{
	1: "project_admin",
	2: "master",
	3: "developer",
	4: "guest",
	5: "limited_guest",
}

func resourceProjectMember() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"group_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
		Create: resourceProjectMemberCreate,
		Read:   resourceProjectMemberRead,
		Update: resourceProjectMemberUpdate,
		Delete: resourceProjectMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func getRoleName(name string) (int64, error) {
	if num, ok := roleName2Id[name]; ok {
		return num, nil
	}
	return 0, fmt.Errorf("role id \"%s\" is unknown", name)
}

func resourceProjectMemberCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	usergroupTypeNum, err := getGroupTypeName(d.Get("group_type").(string))
	if err != nil {
		return err
	}
	roleId, err := getRoleName(d.Get("role").(string))
	if err != nil {
		return err
	}

	groupName := d.Get("group_name").(string)

	projectMember := &models.ProjectMember{
		RoleID: roleId,
		MemberGroup: &models.UserGroup{
			GroupType: usergroupTypeNum,
			GroupName: groupName,
		},
	}

	projectID := int64(d.Get("project_id").(int))

	_, err = apiClient.Products.PostProjectsProjectIDMembers(
		products.NewPostProjectsProjectIDMembersParams().
			WithProjectID(projectID).
			WithProjectMember(projectMember),
		nil,
	)

	if err != nil {
		return fmt.Errorf("project member creation failed")
	}

	listResp, err := apiClient.Products.GetProjectsProjectIDMembers(
		products.NewGetProjectsProjectIDMembersParams().
			WithProjectID(projectID),
		nil,
	)

	if err != nil {
		return fmt.Errorf("project member loading failed")
	}

	var foundMember *models.ProjectMemberEntity
	for _, member := range listResp.Payload {
		if member.EntityType == "g" && member.EntityName == groupName {
			foundMember = member
			break
		}
	}

	if foundMember == nil {
		return fmt.Errorf("could not find member %s", groupName)
	}

	groupRef, err := apiClient.Products.GetUsergroupsGroupID(
		products.NewGetUsergroupsGroupIDParams().
			WithGroupID(foundMember.EntityID),
		nil,
	)

	if err != nil {
		return fmt.Errorf("could not find member %s reference, error %v", groupName, err)
	}

	if err := resourceProjectMemberRefresh(d, foundMember, groupRef.Payload.GroupType); err != nil {
		return err
	}

	return nil
}

func resourceProjectMemberRefresh(d *schema.ResourceData, r *models.ProjectMemberEntity, groupType int64) error {
	d.SetId(fmt.Sprintf("%d/%d", r.ProjectID, r.ID))

	roleName := roleId2Name[r.RoleID]
	if err := d.Set("role", roleName); err != nil {
		return err
	}
	usergroupTypeName := groupTypeNum2Name[groupType]
	if err := d.Set("group_type", usergroupTypeName); err != nil {
		return err
	}
	if err := d.Set("group_name", r.EntityName); err != nil {
		return err
	}
	if err := d.Set("project_id", r.ProjectID); err != nil {
		return err
	}

	return nil
}

func resourceProjectMemberParseID(id string) (int64, int64, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid id %s", id)
	}

	projectID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	memberID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return int64(projectID), int64(memberID), nil
}

func resourceProjectMemberRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	projectID, memberID, err := resourceProjectMemberParseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := apiClient.Products.GetProjectsProjectIDMembersMid(
		products.NewGetProjectsProjectIDMembersMidParams().
			WithProjectID(projectID).
			WithMid(memberID),
		nil,
	)

	if err != nil {
		return err
	}

	groupRef, err := apiClient.Products.GetUsergroupsGroupID(
		products.NewGetUsergroupsGroupIDParams().
			WithGroupID(resp.Payload.EntityID),
		nil,
	)

	if err != nil {
		return fmt.Errorf("could not find member reference for ID %d, error %v", memberID, err)
	}

	if err := resourceProjectMemberRefresh(d, resp.Payload, groupRef.Payload.GroupType); err != nil {
		return err
	}

	return nil
}

func resourceProjectMemberUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	projectID, memberID, err := resourceProjectMemberParseID(d.Id())
	if err != nil {
		return err
	}

	roleId, err := getRoleName(d.Get("role").(string))
	if err != nil {
		return err
	}
	roleRequest := &models.RoleRequest{
		RoleID: roleId,
	}

	body := products.NewPutProjectsProjectIDMembersMidParams().
		WithRole(roleRequest).
		WithProjectID(projectID).
		WithMid(memberID)

	_, err = apiClient.Products.PutProjectsProjectIDMembersMid(body, nil)

	if err != nil {
		return fmt.Errorf("Project member update failed")
	}

	return nil
}

func resourceProjectMemberDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	projectID, memberID, err := resourceProjectMemberParseID(d.Id())
	if err != nil {
		return err
	}

	_, err = apiClient.Products.DeleteProjectsProjectIDMembersMid(
		products.NewDeleteProjectsProjectIDMembersMidParams().
			WithProjectID(projectID).
			WithMid(memberID),
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}
