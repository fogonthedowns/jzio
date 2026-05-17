.PHONY: deploy start help infra infra-dns infra-email infra-site infra-bootstrap infra-setup infra-synth _venv

SITE_STACK  ?= SiteStack
AWS_REGION  ?= us-east-1
VENV        = infra/.venv

export CDK_DEFAULT_ACCOUNT ?= $(shell aws sts get-caller-identity --query Account --output text 2>/dev/null)

help:
	@echo "Targets:"
	@echo "  infra-setup      create Python venv and install dependencies (run once)"
	@echo "  infra-bootstrap  cdk bootstrap for us-east-1 (run once per account)"
	@echo "  infra-synth      synthesise all stacks (dry-run)"
	@echo "  infra-dns        deploy DnsStack — Route53 hosted zone"
	@echo "                   copy NameServers output to your registrar before continuing"
	@echo "  infra-email      deploy EmailStack — FastMail + third-party DNS records"
	@echo "  infra-site       deploy SiteStack  — cert, S3, CloudFront, apex/www records"
	@echo "  infra            deploy all three stacks in dependency order"
	@echo "  deploy           sync ./  to S3 and invalidate CloudFront"
	@echo "  start            start local dev server"
	@echo ""
	@echo "Edit infra/config.py to change domain, email settings, or optional records."

# ── Setup ─────────────────────────────────────────────────────────────────────

_venv:
	@test -f $(VENV)/bin/python || (python3 -m venv $(VENV) && $(VENV)/bin/pip install -q -r infra/requirements.txt)

infra-setup: _venv
	@echo "Venv ready at $(VENV)"

# ── CDK ───────────────────────────────────────────────────────────────────────

infra-bootstrap: _venv
	@test -n "$(CDK_DEFAULT_ACCOUNT)" || (echo "Set CDK_DEFAULT_ACCOUNT or configure AWS CLI." && exit 1)
	cd infra && npx cdk bootstrap aws://$(CDK_DEFAULT_ACCOUNT)/$(AWS_REGION)

infra-synth: _venv
	cd infra && npx cdk synth

infra-dns: _venv
	@test -n "$(CDK_DEFAULT_ACCOUNT)" || (echo "Set CDK_DEFAULT_ACCOUNT or configure AWS CLI." && exit 1)
	cd infra && CDK_DEFAULT_ACCOUNT=$(CDK_DEFAULT_ACCOUNT) npx cdk deploy DnsStack --require-approval never

infra-email: _venv
	@test -n "$(CDK_DEFAULT_ACCOUNT)" || (echo "Set CDK_DEFAULT_ACCOUNT or configure AWS CLI." && exit 1)
	cd infra && CDK_DEFAULT_ACCOUNT=$(CDK_DEFAULT_ACCOUNT) npx cdk deploy EmailStack --require-approval never

infra-site: _venv
	@test -n "$(CDK_DEFAULT_ACCOUNT)" || (echo "Set CDK_DEFAULT_ACCOUNT or configure AWS CLI." && exit 1)
	cd infra && CDK_DEFAULT_ACCOUNT=$(CDK_DEFAULT_ACCOUNT) npx cdk deploy SiteStack --require-approval never

infra: _venv
	@test -n "$(CDK_DEFAULT_ACCOUNT)" || (echo "Set CDK_DEFAULT_ACCOUNT or configure AWS CLI." && exit 1)
	cd infra && CDK_DEFAULT_ACCOUNT=$(CDK_DEFAULT_ACCOUNT) npx cdk deploy DnsStack EmailStack SiteStack --require-approval never

# ── Site deploy ───────────────────────────────────────────────────────────────

deploy:
	@BUCKET=$$(aws cloudformation describe-stacks --stack-name "$(SITE_STACK)" --region "$(AWS_REGION)" --query "Stacks[0].Outputs[?OutputKey=='SiteBucketName'].OutputValue | [0]" --output text 2>/dev/null); \
	DIST_ID=$$(aws cloudformation describe-stacks --stack-name "$(SITE_STACK)" --region "$(AWS_REGION)" --query "Stacks[0].Outputs[?OutputKey=='CloudFrontDistributionId'].OutputValue | [0]" --output text 2>/dev/null); \
	if [ -z "$$BUCKET" ] || [ "$$BUCKET" = "None" ]; then echo "Could not read SiteBucketName from $(SITE_STACK). Run: make infra-site"; exit 1; fi; \
	if [ -z "$$DIST_ID" ] || [ "$$DIST_ID" = "None" ]; then echo "Could not read CloudFrontDistributionId from $(SITE_STACK). Run: make infra-site"; exit 1; fi; \
	echo "Deploying to s3://$$BUCKET (invalidating $$DIST_ID)"; \
	aws s3 sync . "s3://$$BUCKET/" \
		--region "$(AWS_REGION)" \
		--delete \
		--exclude ".git/*" \
		--exclude "infra/*" \
		--exclude ".cursor/*" \
		--exclude ".gitignore" \
		--exclude ".arcconfig"; \
	aws cloudfront create-invalidation --region "$(AWS_REGION)" --distribution-id "$$DIST_ID" --paths "/*"

start:
	npm start
