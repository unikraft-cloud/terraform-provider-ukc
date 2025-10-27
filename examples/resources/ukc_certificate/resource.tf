resource "ukc_certificate" "example" {
  name  = "my-certificate"
  cn    = "example.com"
  chain = file("path/to/certificate.crt")
  pkey  = file("path/to/private.key")
}
