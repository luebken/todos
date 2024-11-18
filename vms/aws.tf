provider "aws" {
  region = "eu-central-1"
}

resource "aws_instance" "dev_vm_1" {
  ami           = "ami-0084a47cc718c111a" # Ubuntu Server 24.04 'ubuntu'
  instance_type = "t2.micro"
  key_name      = "id_rsa_dev"

  # Security Group
  vpc_security_group_ids = [aws_security_group.dev_vm_sg.id]
}

resource "aws_security_group" "dev_vm_sg" {
  name_prefix = "dev-vm-sg-"

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
