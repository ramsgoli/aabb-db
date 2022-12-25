resource "aws_key_pair" "buildkite_agent_keypair" {
  key_name   = "buildkite_agent_key"
  public_key = file("~/.ssh/id_rsa.pub")
}
