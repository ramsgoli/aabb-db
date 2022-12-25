resource "aws_vpc" "buildkite_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "BuildkiteVPC"
  }
}

resource "aws_subnet" "buildkite_vpc_primary_subnet" {
  vpc_id     = aws_vpc.buildkite_vpc.id
  cidr_block = "10.0.0.0/24"

  map_public_ip_on_launch = true

  tags = {
    Name = "BuildkiteVPCPrimarySubnet"
  }
}

resource "aws_security_group" "buildkite_agent_sg" {
  name        = "buildkite-agent-sg"
  description = "Allow traffic to/from Buildkite Agent"

  vpc_id = aws_vpc.buildkite_vpc.id

  ingress {
    to_port     = 22
    from_port   = 22
    protocol    = "tcp"
    cidr_blocks = ["205.220.129.20/32"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "buildkite-agent-sg"
  }
}

# Internet Gateway provides a VPC access to the internet
resource "aws_internet_gateway" "buildkite_internet_gw" {
  vpc_id = aws_vpc.buildkite_vpc.id

  tags = {
    Name = "BuildkiteInternetGateway"
  }
}

resource "aws_route_table" "buildkite_public_subnet_route_table" {
  vpc_id = aws_vpc.buildkite_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.buildkite_internet_gw.id
  }

  tags = {
    Name = "BuildkitePublicSubnetRouteTable"
  }
}

# Attaching a route from a subnet to the IG effectively makes it a public subnet (reachable from the internet)
resource "aws_route_table_association" "buildkite_public_subnet_route_table_association" {
  subnet_id      = aws_subnet.buildkite_vpc_primary_subnet.id
  route_table_id = aws_route_table.buildkite_public_subnet_route_table.id
}
