resource "harbor_project" "main" {
  name                    = "main"
  public                  = false # (Optional) Default value is false
  vulnerability_scanning  = true  # (Optional) Default value is true. Automatically scan images on push
  reuse_sys_cve_whitelist = false # (Optional) Default value is true.
  cve_whitelist           = ["CVE-2020-12345", "CVE-2020-54321"]
}
