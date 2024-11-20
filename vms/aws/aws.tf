provider "aws" {
  region = "eu-central-1"
}

# aws ec2 describe-vpcs --filters "Name=tag:Name,Values=MainVPC"
resource "aws_vpc" "main" {
  tags                 = { Name = "MainVPC" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "main" {
  tags   = { Name = "MainInternetGateway" }
  vpc_id = aws_vpc.main.id

}

# aws ec2 describe-subnets --filters "Name=tag:Name,Values=PublicSubnet"
# aws ec2 describe-subnets --filters "Name=tag:Name,Values=PublicSubnet" | jq -r .Subnets[0].SubnetId
resource "aws_subnet" "public_subnet" {
  tags                    = { Name = "PublicSubnet" }
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "eu-central-1a"
}

# aws ec2 describe-route-tables --filters "Name=association.subnet-id,Values="
resource "aws_route_table" "main" {
  tags   = { Name = "MainRouteTable" }
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}
resource "aws_route_table_association" "public_subnet_association" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.main.id
}

# aws ec2 describe-instances --filters "Name=tag:Name,Values=AppVM1"
resource "aws_instance" "app_vm_1" {
  tags          = { Name = "AppVM1" }
  ami           = "ami-0084a47cc718c111a" # Ubuntu Server 24.04 user:'ubuntu'
  instance_type = "t2.micro"
  key_name      = "id_rsa_dev"
  subnet_id     = aws_subnet.public_subnet.id

  vpc_security_group_ids = [aws_security_group.dev_vm_sg.id]
}

# aws ec2 describe-security-groups
resource "aws_security_group" "dev_vm_sg" {
  name_prefix = "dev-vm-sg-"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "Allow SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow HTTP"
    from_port   = 8000
    to_port     = 8000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow All ICMP - IPv4"
    from_port   = -1 # -1 means all ICMP types
    to_port     = -1 # -1 means all ICMP codes
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
