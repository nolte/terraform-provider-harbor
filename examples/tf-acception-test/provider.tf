


# example for harbor v1 api usage
provider "harbor" {
  host     = "harbor.172-17-0-1.sslip.io"
  schema   = "https"
  insecure = true
  basepath = "/api/v2.0"
  username = "admin"
  password = "Harbor12345"
}
