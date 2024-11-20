variable "emma_client_id" { type = string }
variable "emma_client_secret" { type = string }

terraform {
  required_providers {
    emma = {
      source  = "emma-community/emma"
      version = "0.0.1"
    }
  }
}

provider "emma" {
  client_id     = var.emma_client_id
  client_secret = var.emma_client_secret
}

resource "emma_vm" "vm" {
  name               = "vm-test1"
  data_center_id     = "aws-eu-central-1"
  os_id              = 5
  cloud_network_type = "multi-cloud"
  vcpu_type          = "shared"
  ssh_key_id         = 918
  volume_type        = "ssd"
  ram_gb             = 1
  volume_gb          = 16
  vcpu               = 2
}
