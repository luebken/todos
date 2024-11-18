# VMs

For running the application as VM we ware re-using the Docker containers and starting them with systemd. As a container runtime we are using [podman](https://podman.io/).

```sh
terraform init / plan / apply

PUBLIC_IP_1=$(terraform show -json | jq -r .values.root_module.resources[0].values.public_ip)
ssh -i ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP_1

# Install Podman
sudo apt-get update
sudo apt-get -y install podman

# Export & import the container
docker save luebken/todos -o todos.save
podman import todos.save

# Upload the service definition and image
scp -i ~/.ssh/id_rsa_dev todos.service todos.service ubuntu@$PUBLIC_IP_1:/home/ubuntu

# Setup and start the service
sudo chown root:root todos.service
sudo mv todos.service /etc/systemd/system/

sudo systemctl daemon-reload
sudo systemctl enable todos
sudo systemctl start todos

# Troubleshooting
sudo systemctl status todos
sudo journalctl -u todos.service -e
```
