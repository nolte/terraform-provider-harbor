resource "harbor_config_email" "conf_email" {
  email_host     = "main2"
  email_port     = 25
  email_username = "main2"
  email_password = "main2"
  email_from     = "main2"
  email_ssl      = false
}

resource "harbor_config_auth" "oidc" {
  auth_mode          = "oidc_auth"
  oidc_name          = "azure"
  oidc_endpoint      = "https://login.microsoftonline.com/v2.0"
  oidc_client_id     = "OIDC Client ID goes here"
  oidc_client_secret = "ODDC Client Secret goes here"
  oidc_scope         = "openid,email"
  oidc_verify_cert   = true
}

resource "harbor_config_system" "main" {
  project_creation_restriction = "everyone"
  robot_token_expiration       = 5259492
}

