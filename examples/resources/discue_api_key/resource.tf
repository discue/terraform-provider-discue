resource "discue_api_key" "test_status" {
  alias  = "my-first-api-key"
  status = "enabled"
  scopes = [{
    resource = "messages"
    access   = "write"
    targets  = ["*"]
  }]
}
