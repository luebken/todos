# Setting up VMs with AWS

```sh
## Install
```sh
terraform init / plan / apply

terraform refresh
PUBLIC_IP=$(terraform show -json | jq -r .values.root_module.resources[0].values.public_ip)
ssh -i ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP
```