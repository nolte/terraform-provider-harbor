resource "harbor_webhook_policy" "main" {
  name         = "test-policy"
  description  = "Testing"
  project_id   = harbor_project.main.id
  endpoint_url = "https://www.googgle.com"
  event_types  = ["SCANNING_COMPLETED"]
}
