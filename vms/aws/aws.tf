locals {
  region              = "eu-central-1"
  availability_zone_a = "eu-central-1a"

  ami = "ami-0084a47cc718c111a" # Ubuntu Server 24.04 user:'ubuntu'
  instance_type = "t2.micro"
  key_name      = "id_rsa_dev"
}
provider "aws" {
  region = local.region
}

# Main VPC
#
resource "aws_vpc" "main" {
  tags                 = { Name = "main_vpc" }
  cidr_block           = "10.0.0.0/16" // from /16 to /28
  enable_dns_support   = true
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
}

# Public Web-Subnet
#
resource "aws_subnet" "web_subnet" {
  tags                    = { Name = "web_subnet" }
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = local.availability_zone_a
}
resource "aws_route_table" "web_subnet_rt" {
  tags   = { Name = "web_subnet_rt" }
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}
resource "aws_route_table_association" "web_subnet_association" {
  subnet_id      = aws_subnet.web_subnet.id
  route_table_id = aws_route_table.web_subnet_rt.id
}

# Public Web Instances
# 
resource "aws_instance" "webapp_vm_1" {
  tags          = { Name = "webapp_vm_1" }
  ami           = local.ami
  instance_type = local.instance_type
  key_name      = local.key_name
  subnet_id     = aws_subnet.web_subnet.id

  vpc_security_group_ids = [aws_security_group.webapp_vm_sg.id]
}
resource "aws_instance" "webapp_vm_2" {
  tags          = { Name = "webapp_vm_2" }
  ami           = local.ami
  instance_type = local.instance_type
  key_name      = local.key_name
  subnet_id     = aws_subnet.web_subnet.id

  vpc_security_group_ids = [aws_security_group.webapp_vm_sg.id]
}
resource "aws_security_group" "webapp_vm_sg" {
  name_prefix = "webapp_vm_1_sg-"
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
    from_port   = -1 # all ICMP types
    to_port     = -1
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

# -----------------

# Private Data-Subnet
# 
resource "aws_subnet" "data_subnet" {
  tags                    = { Name = "data_subnet" }
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = local.availability_zone_a
  map_public_ip_on_launch = false # !!!
}
resource "aws_route_table" "data_subnet_rt" {
  tags   = { Name = "data_subnet_rt" }
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }
}
resource "aws_route_table_association" "data_subnet_association" {
  subnet_id      = aws_subnet.data_subnet.id
  route_table_id = aws_route_table.data_subnet_rt.id
}

# NAT Gateway
# Enable resources in the private subnet to access the internet.
# Create an Elastic IP for the NAT Gateway
resource "aws_eip" "nat_eip" {
  vpc = true
}
resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat_eip.id
  subnet_id     = aws_subnet.web_subnet.id
}

# Private Data Instances
#
resource "aws_instance" "data_vm_1" {
  tags          = { Name = "data_vm_1" }
  ami           = local.ami
  instance_type = local.instance_type
  key_name      = local.key_name
  subnet_id     = aws_subnet.data_subnet.id

  vpc_security_group_ids = [aws_security_group.data_vm_sg.id]
}
resource "aws_instance" "data_vm_2" {
  tags          = { Name = "data_vm_2" }
  ami           = local.ami
  instance_type = local.instance_type
  key_name      = local.key_name
  subnet_id     = aws_subnet.data_subnet.id

  vpc_security_group_ids = [aws_security_group.data_vm_sg.id]
}

#TODO only allow access for SSH and Postgres from the bastion host 
resource "aws_security_group" "data_vm_sg" {
  name_prefix = "data_vm_sg-"
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
    from_port   = -1 # all ICMP types
    to_port     = -1
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
