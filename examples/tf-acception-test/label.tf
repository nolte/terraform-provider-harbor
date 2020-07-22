resource "harbor_label" "project_label" {
  name        = "projectlabel-acc-classic"
  description = "Test Label for Project"
  color       = "#333333"
  scope       = "p"
  project_id  = harbor_project.main.id
}

resource "harbor_label" "main" {
  name        = "testlabel-acc-classic"
  description = "Test Label"
  color       = "#61717D"
  scope       = "g"
}

