resource "ukc_service" "example" {
  services {
    port             = 443
    destination_port = 8080
    handlers         = ["tls", "http"]
  }
}
