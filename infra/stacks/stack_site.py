from typing import Optional

import aws_cdk as cdk
import aws_cdk.aws_certificatemanager as acm
import aws_cdk.aws_cloudfront as cloudfront
import aws_cdk.aws_cloudfront_origins as origins
import aws_cdk.aws_route53 as route53
import aws_cdk.aws_route53_targets as targets
import aws_cdk.aws_s3 as s3
from constructs import Construct


class SiteStack(cdk.Stack):
    """S3 bucket, CloudFront distribution, and site DNS records — deployed to us-west-2.

    Certificate lives in CertStack (us-east-1) because CloudFront requires it there.

    When called from the pipeline, hosted_zone and certificate are omitted and
    imported by reference using HOSTED_ZONE_ID / CERT_ARN from config.py.
    When called from app.py for manual deploys, pass the live CDK objects.
    """

    def __init__(
        self,
        scope: Construct,
        id: str,
        hosted_zone: Optional[route53.IHostedZone] = None,
        certificate: Optional[acm.ICertificate] = None,
        site_bucket_name: Optional[str] = None,
        **kwargs,
    ):
        super().__init__(scope, id, **kwargs)

        if hosted_zone is None:
            from config import DOMAIN, HOSTED_ZONE_ID
            hosted_zone = route53.HostedZone.from_hosted_zone_attributes(
                self, "ImportedZone",
                hosted_zone_id=HOSTED_ZONE_ID,
                zone_name=DOMAIN,
            )
        if certificate is None:
            from config import CERT_ARN
            certificate = acm.Certificate.from_certificate_arn(self, "ImportedCert", CERT_ARN)

        domain = hosted_zone.zone_name
        www_domain = f"www.{domain}"

        # ── S3 bucket ─────────────────────────────────────────────────────────
        self.site_bucket = s3.Bucket(
            self, "SiteBucket",
            bucket_name=site_bucket_name or None,
            block_public_access=s3.BlockPublicAccess.BLOCK_ALL,
            encryption=s3.BucketEncryption.S3_MANAGED,
            enforce_ssl=True,
            removal_policy=cdk.RemovalPolicy.RETAIN,
            auto_delete_objects=False,
        )

        # ── CloudFront distribution ───────────────────────────────────────────
        self.distribution = cloudfront.Distribution(
            self, "SiteDist",
            comment=f"Static site for {domain}",
            default_behavior=cloudfront.BehaviorOptions(
                origin=origins.S3BucketOrigin.with_origin_access_control(self.site_bucket),
                viewer_protocol_policy=cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
                allowed_methods=cloudfront.AllowedMethods.ALLOW_GET_HEAD_OPTIONS,
                cached_methods=cloudfront.CachedMethods.CACHE_GET_HEAD_OPTIONS,
                compress=True,
            ),
            domain_names=[domain, www_domain],
            certificate=certificate,  # from CertStack (us-east-1)
            default_root_object="index.html",
            error_responses=[
                cloudfront.ErrorResponse(
                    http_status=403,
                    response_http_status=200,
                    response_page_path="/index.html",
                    ttl=cdk.Duration.minutes(5),
                ),
                cloudfront.ErrorResponse(
                    http_status=404,
                    response_http_status=200,
                    response_page_path="/index.html",
                    ttl=cdk.Duration.minutes(5),
                ),
            ],
            http_version=cloudfront.HttpVersion.HTTP2_AND_3,
            price_class=cloudfront.PriceClass.PRICE_CLASS_100,
            minimum_protocol_version=cloudfront.SecurityPolicyProtocol.TLS_V1_2_2021,
        )

        # ── DNS records (apex + www) ──────────────────────────────────────────
        cf_target = route53.RecordTarget.from_alias(
            targets.CloudFrontTarget(self.distribution)
        )
        route53.ARecord(self, "ApexA", zone=hosted_zone, target=cf_target)
        route53.AaaaRecord(self, "ApexAAAA", zone=hosted_zone, target=cf_target)
        route53.ARecord(self, "WwwA", zone=hosted_zone, record_name="www", target=cf_target)
        route53.AaaaRecord(self, "WwwAAAA", zone=hosted_zone, record_name="www", target=cf_target)

        # ── Outputs ───────────────────────────────────────────────────────────
        self.bucket_name_cfn_output = cdk.CfnOutput(
            self, "SiteBucketName", value=self.site_bucket.bucket_name,
            description="S3 bucket for site content",
        )
        self.distribution_id_cfn_output = cdk.CfnOutput(
            self, "CloudFrontDistributionId", value=self.distribution.distribution_id,
            description="CloudFront distribution ID for cache invalidation",
        )
        cdk.CfnOutput(self, "CloudFrontDomainName", value=self.distribution.distribution_domain_name,
                      description="CloudFront domain (*.cloudfront.net)")
        cdk.CfnOutput(self, "SiteUrl", value=f"https://{domain}",
                      description="Primary site URL")
