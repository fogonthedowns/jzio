#!/usr/bin/env python3
import aws_cdk as cdk

from config import AWS_ACCOUNT, DOMAIN, LOG_BUCKET_NAME, SITE_BUCKET_NAME
from stacks.stack_cert import CertStack
from stacks.stack_dns import DnsStack
from stacks.stack_email import EmailStack
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
    log_bucket_name=LOG_BUCKET_NAME or None,
    env=env_west,
    cross_region_references=True,
)
site_stack.add_dependency(cert_stack)

app.synth()
