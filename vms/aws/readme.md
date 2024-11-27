# Setting up VMs with AWS


## Create Infra
```sh
# Configure & Test AWS access
aws sts get-caller-identity

tofu init
tofu plan
tofu apply

chmod 600 id_rsa_todos.pem
```

## SSH Access
```sh
PUBLIC_IP_APP_1=$(tofu show --json | jq -r '.values.root_module.resources[] | select(.address=="aws_instance.webapp_vm_1").values.public_ip')
ssh -i id_rsa_todos.pem ubuntu@$PUBLIC_IP_APP_1

PUBLIC_IP_APP_2=$(tofu show --json | jq -r '.values.root_module.resources[] | select(.address=="aws_instance.webapp_vm_2").values.public_ip')
ssh -i id_rsa_todos.pem ubuntu@$PUBLIC_IP_APP_2
```

## Network overview

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
```

## Configure

```sh
# Check if servce came up
# TODO replace with TODOs
systemctl status podman
cat /var/log/cloud-init-output.log
```
See also https://dev.to/gdenn/aws-best-practices-three-tier-vpc-37n0