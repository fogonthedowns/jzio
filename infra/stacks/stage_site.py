import aws_cdk as cdk
from constructs import Construct

from config import AWS_ACCOUNT, SITE_BUCKET_NAME
from stacks.stack_site import SiteStack


class SiteStage(cdk.Stage):
    """Pipeline stage that manages only SiteStack.

    DnsStack, EmailStack, and CertStack are never touched by the pipeline.
    SiteStack imports the hosted zone and certificate by reference from config.py.
    """

    def __init__(self, scope: Construct, id: str, **kwargs):
        super().__init__(scope, id, **kwargs)

        self.site_stack = SiteStack(
            self, "SiteStack",
            stack_name="SiteStack",
            site_bucket_name=SITE_BUCKET_NAME or cdk.PhysicalName.GENERATE_IF_NEEDED,
            env=cdk.Environment(account=AWS_ACCOUNT, region="us-west-2"),
            cross_region_references=True,
        )
