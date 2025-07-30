APP_NAME=gobid
REGION=us-east-1
VPC_ID=vpc-096fc31bf0f9ef9e0

ACCOUNT_ID=$(shell aws sts get-caller-identity --query Account --output text)
export AWS_PAGER=
SG_NAME=$(APP_NAME)-sg
ECR_URL=$(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com
REPO_URL=$(ECR_URL)/$(APP_NAME)

create-sg:
	@if ! aws ec2 describe-security-groups --filters "Name=group-name,Values=$(SG_NAME)" --region $(REGION) --query "SecurityGroups[*].GroupId" --output text | grep -qE 'sg-'; then \
		echo "Creating security group $(SG_NAME)..."; \
		aws ec2 create-security-group \
			--group-name $(SG_NAME) \
			--description "Allow Postgres for GoBid" \
			--vpc-id $(VPC_ID) \
			--region $(REGION); \
	else \
		echo "Security group $(SG_NAME) already exists."; \
	fi
	@SG_ID=$$(aws ec2 describe-security-groups --filters "Name=group-name,Values=$(SG_NAME)" --region $(REGION) --query "SecurityGroups[*].GroupId" --output text); \
	echo "Authorizing port 5432 ON SG $$SG_ID..."; \
	AWS_PAGER= aws ec2 authorize-security-group-ingress \
		--group-id $$SG_ID \
		--protocol tcp \
		--port 5432 \
		--cidr 0.0.0.0/0 \
		--region $(REGION) || echo "Ingress rule already exists or failed silently."

create-ecr:
	aws ecr describe-repositories --repository-name $(APP_NAME) --region $(REGION) || \
	aws ecr create-repository --repository-name $(APP_NAME) --region $(REGION)

build:
	docker build -t $(APP_NAME) -f Dockerfile.prod .

tag: build
	docker tag $(APP_NAME):latest $(REPO_URL):latest

push: tag
	aws ecr get-login-password --region $(REGION) | \
	docker login --username AWS --password-stdin $(ECR_URL)
	docker push $(REPO_URL):latest
