resource "harbor_project_member" "developers_main" {
    project_id = harbor_project.main.id
    role       = "guest"
    group_type = "http"
    group_name = harbor_usergroup.developers.name
}
