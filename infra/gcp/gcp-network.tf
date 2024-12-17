# Main VPC
# 
resource "google_compute_network" "main_vpc" {
  name                    = "main-vpc"
  auto_create_subnetworks = false
}

# Public Web-Subnet
# 
resource "google_compute_subnetwork" "public_subnet" {
  name          = "public-subnet"
  network       = google_compute_network.main_vpc.id
  ip_cidr_range = "10.0.1.0/24"
  region        = local.region
}

# Private Subnet
resource "google_compute_subnetwork" "private_subnet" {
  name          = "private-subnet"
  network       = google_compute_network.main_vpc.id
  ip_cidr_range = "10.0.2.0/24"
  region        = local.region
}

# NAT Gateway (Cloud NAT)
resource "google_compute_router" "router" {
  name    = "main-router"
  region  = local.region
  network = google_compute_network.main_vpc.id
}

resource "google_compute_router_nat" "nat_gateway" {
  name                               = "main-nat"
  router                             = google_compute_router.router.name
  region                             = local.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.private_subnet.id
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}