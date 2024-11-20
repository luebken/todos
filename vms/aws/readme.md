# Setting up VMs with AWS

```sh
# Configure & Test AWS access
aws sts get-caller-identity

## Install
```sh
terraform init / plan / apply


PUBLIC_IP=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=DevVM1" | jq -r '.Reservations[0].Instances[0].PublicIpAddress')
ssh -i ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP
```