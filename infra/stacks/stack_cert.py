import aws_cdk as cdk
import aws_cdk.aws_certificatemanager as acm
import aws_cdk.aws_route53 as route53
from constructs import Construct


class CertStack(cdk.Stack):
    """ACM certificate for CloudFront — must live in us-east-1 regardless of site region."""

    def __init__(self, scope: Construct, id: str, hosted_zone: route53.IHostedZone, **kwargs):
        super().__init__(scope, id, **kwargs)

        self.certificate = acm.Certificate(
            self, "SiteCert",
            domain_name=hosted_zone.zone_name,
            subject_alternative_names=[f"*.{hosted_zone.zone_name}"],
            validation=acm.CertificateValidation.from_dns(hosted_zone),
        )

        cdk.CfnOutput(self, "CertArn", value=self.certificate.certificate_arn,
                      description="ACM certificate ARN (us-east-1, required by CloudFront)")
