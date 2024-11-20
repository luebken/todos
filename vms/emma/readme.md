
# Setting up VMs with AWS

```sh
# Set client id and secret as Envs
export TF_VAR_emma_client_id=
export TF_VAR_emma_client_secret=

# Issue Token
EMMA_ACCESS_TOKEN=$(curl -s -X POST https://api.emma.ms/external/v1/issue-token -H "Content-Type: application/json" -d '{"clientId": "'"$TF_VAR_emma_client_id"'","clientSecret": "'"$TF_VAR_emma_client_secret"'"}' | jq -r .accessToken)

# Import SSH Key
SSH_PUBLIC_KEY=$(cat ~/.ssh/rsa_pub)
curl -X POST 'https://api.emma.ms/external/v1/ssh-keys' -H 'Content-Type: application/json' -H "Authorization: Bearer $EMMA_ACCESS_TOKEN"  -d '{"name": "mykey","key": "'"$SSH_PUBLIC_KEY"'"}'

# Login
terraform refresh
PUBLIC_IP=$(terraform show -json | jq -r .values.root_module.resources[0].values.networks[1].ip)
ssh -i ~/.ssh/id_rsa_dev ubuntu@$PUBLIC_IP
```