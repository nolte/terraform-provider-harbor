
provider "harbor" {
  url      = "harbor.192-168-178-51.sslip.io"
  insecure = true
  #url      = "demo.goharbor.io"
  #basepath = "/api/v2.0"
  username = "admin"
  password = "Harbor12345"
}

#resource "harbor_project" "main" {
#  name                   = "main"
#  public                 = false # (Optional) Default value is false
#  vulnerability_scanning = true  # (Optional) Default vale is true. Automatically scan images on push 
#}

resource "harbor_label" "main" {
  name  = "testlabel"
  scope = "g"
}
