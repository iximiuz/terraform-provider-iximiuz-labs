resource "iximiuz-labs_play_shell" "terminal" {
  play_id = iximiuz-labs_play.example.id
  machine = "docker-01"
  user    = "root"
  access  = "public"
}

output "shell_url" {
  value = iximiuz-labs_play_shell.terminal.url
}
