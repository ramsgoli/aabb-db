resource "aws_instance" "buildkite_agent" {
  # Amazon Linux 2 (64-bit x86)
  ami           = "ami-00d8a762cb0c50254"
  instance_type = "t2.micro"

  subnet_id                   = aws_subnet.buildkite_vpc_primary_subnet.id
  vpc_security_group_ids      = [aws_security_group.buildkite_agent_sg.id]
  associate_public_ip_address = true
  key_name                    = aws_key_pair.buildkite_agent_keypair.key_name

  depends_on = [
    aws_internet_gateway.buildkite_internet_gw
  ]

  tags = {
    Name = "BuildkiteAgent"
  }
}
