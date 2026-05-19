#!/usr/bin/env python3
import aws_cdk as cdk

from config import AWS_ACCOUNT, DOMAIN, GITHUB_CONNECTION_ARN, SITE_BUCKET_NAME
from stacks.stack_cert import CertStack
from stacks.stack_dns import DnsStack
from stacks.stack_email import EmailStack
from stacks.stack_pipeline import PipelineStack
from stacks.stack_site import SiteStack

app = cdk.App()

# DNS and email stay in us-east-1. ACM cert must also be us-east-1 (CloudFront requirement).
# Site bucket + CloudFront + DNS records live in us-west-2.
env_east = cdk.Environment(account=AWS_ACCOUNT, region="us-east-1")
env_west = cdk.Environment(account=AWS_ACCOUNT, region="us-west-2")

dns_stack = DnsStack(app, "DnsStack", domain_name=DOMAIN, env=env_east)

email_stack = EmailStack(app, "EmailStack", hosted_zone=dns_stack.hosted_zone, env=env_east)
email_stack.add_dependency(dns_stack)

cert_stack = CertStack(app, "CertStack", hosted_zone=dns_stack.hosted_zone, env=env_east)
cert_stack.add_dependency(dns_stack)

# cross_region_references=True lets SiteStack (us-west-2) consume the cert from CertStack (us-east-1).
site_stack = SiteStack(
    app, "SiteStack",
    hosted_zone=dns_stack.hosted_zone,
    certificate=cert_stack.certificate,
    site_bucket_name=SITE_BUCKET_NAME or None,
    env=env_west,
    cross_region_references=True,
)
site_stack.add_dependency(cert_stack)

# Self-mutating pipeline: watches master, runs cdk deploy + s3 sync on every push.
# Pipeline lives in us-west-1 to match the CodeConnections connection region.
env_west1 = cdk.Environment(account=AWS_ACCOUNT, region="us-west-1")
PipelineStack(app, "PipelineStack", connection_arn=GITHUB_CONNECTION_ARN, env=env_west1)

app.synth()
