destroy:
	terraform destroy

create:
	terraform plan
	terraform apply -auto-approve

ssh:
	PUBLIC_IP=`aws ec2 describe-instances --filters "Name=tag:Name,Values=webapp_vm_1" | jq -r '.Reservations[0].Instances[0].PublicIpAddress'`; \
	echo $$PUBLIC_IP; \
	ssh -i ~/.ssh/id_rsa_dev ubuntu@$$PUBLIC_IP