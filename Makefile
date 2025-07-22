APP_NAME=gobid
REGION=us-east-1
VPC_ID=vpc-096fc31bf0f9ef9e0
SG_NAME=gobid-sg

create-sg:
	@if ! aws ec2 describe-security-groups --filters "Name=group-name,Values=$(SG_NAME)" --region $(REGION) --query "SecurityGroups[*].GroupId" --output text | grep -qE 'sg-'; then \
		echo "Creating security group $(SG_NAME)..."; \
		aws ec2 create-security-group \
			--group-name $(SG_NAME)
			--description "Allow Postgres for GoBid" \
			--vpc_id $(VPC_ID) \
			--region $(REGION); \
	else \
		echo "Security group $(SG_NAME) already existis.";
	fi;
	SG_ID=$$(aws ec2 describe-security-groups --filters "Name=group-name,Values=$(SG_NAME)" --region $(REGION) --query "SecurityGroups[*].GroupId" --output text); \
	echo "Authorizing port 5432 ON SG $$SG_ID..."; \
	aws ec2 authorize-security-group-ingress \
		--group-id $$SG_ID \
		--protocol tcp \
		--port 5432 \
		--cidr 0.0.0.0/0 \
		--region $(REGION); \ || echo "Ingress rule already exists or failed silently."
