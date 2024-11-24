# Setting up VMs with AWS

```sh
# Configure & Test AWS access
aws sts get-caller-identity

# Create
terraform plan
terraform apply
```

# Network overview

* A VPC (in eu-central-1)

* "web_subnet" for public access (10.0.1.0/24)
  * VMs receive static IP addresses (via Internet Gateway)
  * RouteTable "web_subnet_rt" to connect to internet
  * Two VMs "webapp_vm_1" and "webapp_vm_2" with a Security group "webapp_vm_sg"
  
* "data_subnet" for private access (10.0.2.0/24)
  * Natgateway to allow only access from the inside 
  * RouteTable "data_subnet_rt" 
  * Two VMs "data_vm_1" and "data_vm_2" with a Security group "data_vm_sg"

```shell
aws ec2 describe-vpcs --filters "Name=tag:Name,Values=main_vpc"

aws ec2 describe-subnets --filters "Name=tag:Name,Values=web_subnet,data_subnet" | jq -r .Subnets[].SubnetId

aws ec2 describe-instances --filters "Name=tag:Name,Values=webapp_vm_1,webapp_vm_2,data_vm_1,data_vm_2" | jq '.Reservations[].Instances[0].InstanceId'

# Access to web VMs
aws ec2 describe-instances --filters "Name=tag:Name,Values=webapp_vm_1,webapp_vm_2" | jq -r '.Reservations[].Instances[].PublicIpAddress'
PUBLIC_IP1=
ssh -i ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP1

# Check if VM came up
cat /var/log/cloud-init-output.log
# TODO replace with TODOs
systemctl status podman

# Access to data VMs
# We use the first WebVM as a bastion host
DATA_PRIVATE_IP1=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=data_vm_1" | jq -r '.Reservations[0].Instances[0].PrivateIpAddress')
# memorize the private IP
echo $DATA_PRIVATE_IP1
# copy scp key to bastion host
scp -i ~/.ssh/id_rsa_dev ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP1:/home/ubuntu
ssh -i ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP1
# on the bastion host
ssh -i id_rsa_dev ubuntu@$DATA_PRIVATE_IP1
# TODO check  Port Forwarding via SSH Tunnel

```
See also https://dev.to/gdenn/aws-best-practices-three-tier-vpc-37n0