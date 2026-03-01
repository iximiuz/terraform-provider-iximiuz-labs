resource "iximiuz-labs_play_port" "web" {
  play_id = iximiuz-labs_play.example.id
  machine = "docker-01"
  number  = 8080
  access  = "public"
  tls     = true
}

output "port_url" {
  value = iximiuz-labs_play_port.web.url
}
