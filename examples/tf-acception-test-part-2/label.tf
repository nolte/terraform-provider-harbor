resource "harbor_label" "project_label" {
  name        = "projectlabel-acc-classic"
  description = "Test Label for Project"
  color       = "#333333"
  scope       = "p"
  project_id  = data.harbor_project.project_1.id
}
