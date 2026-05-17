import aws_cdk as cdk
import aws_cdk.aws_route53 as route53
from constructs import Construct


class DnsStack(cdk.Stack):
    """Stack A — Route53 hosted zone.

    Deploy this first. Copy the output NameServers to your registrar before
    deploying EmailStack or SiteStack (cert DNS validation will hang otherwise).
    """

    def __init__(self, scope: Construct, id: str, domain_name: str, **kwargs):
        kwargs.setdefault("termination_protection", True)
        super().__init__(scope, id, **kwargs)

        self.hosted_zone = route53.PublicHostedZone(
            self, "HostedZone", zone_name=domain_name
        )
        # Keep the hosted zone even if the stack is somehow removed.
        self.hosted_zone.apply_removal_policy(cdk.RemovalPolicy.RETAIN)

        cdk.CfnOutput(
            self, "HostedZoneId",
            value=self.hosted_zone.hosted_zone_id,
            description="Route53 hosted zone ID",
        )
        cdk.CfnOutput(
            self, "NameServers",
            value=cdk.Fn.join(", ", self.hosted_zone.hosted_zone_name_servers or []),
            description="Set these NS records at your registrar before deploying EmailStack/SiteStack",
        )
