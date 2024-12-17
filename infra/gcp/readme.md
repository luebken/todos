# Setting up VMs with GCP

WIP !!!

## Create Infra
```sh
gcloud auth login

tofu init
tofu plan
tofu apply

chmod 600 id_rsa_todos.pem
```

## SSH Access
```sh
PUBLIC_IP_APP_1=$(tofu show --json | jq -r '.values.root_module.resources[] | select(.address=="google_compute_instance.web_vm_1").values.network_interface[0].access_config[0].nat_ip')
ssh -i id_rsa_todos.pem ubuntu@$PUBLIC_IP_APP_1
```