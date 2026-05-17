#!/usr/bin/env python3
import aws_cdk as cdk

from config import DOMAIN, SITE_BUCKET_NAME
from stacks.stack_dns import DnsStack
from stacks.stack_email import EmailStack
from stacks.stack_site import SiteStack

app = cdk.App()

env = cdk.Environment(region="us-east-1")

dns_stack = DnsStack(app, "DnsStack", domain_name=DOMAIN, env=env)

email_stack = EmailStack(app, "EmailStack", hosted_zone=dns_stack.hosted_zone, env=env)
email_stack.add_dependency(dns_stack)

site_stack = SiteStack(
    app, "SiteStack",
    hosted_zone=dns_stack.hosted_zone,
    site_bucket_name=SITE_BUCKET_NAME or None,
    env=env,
)
site_stack.add_dependency(dns_stack)

app.synth()
