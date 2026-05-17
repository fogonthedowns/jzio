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
    """Stack C — ACM certificate, S3 bucket, CloudFront distribution, site DNS records.

    Depends on DnsStack (cert DNS validation requires the hosted zone to exist and
    the registrar NS records to be pointing at Route53).
    """

    def __init__(
        self,
        scope: Construct,
        id: str,
        hosted_zone: route53.IHostedZone,
        site_bucket_name: Optional[str] = None,
        **kwargs,
    ):
        super().__init__(scope, id, **kwargs)

        domain = hosted_zone.zone_name
        www_domain = f"www.{domain}"

        # ── ACM certificate ───────────────────────────────────────────────────
        certificate = acm.Certificate(
            self, "SiteCert",
            domain_name=domain,
            subject_alternative_names=[f"*.{domain}"],
            validation=acm.CertificateValidation.from_dns(hosted_zone),
        )

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
            certificate=certificate,
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
        cdk.CfnOutput(self, "SiteBucketName", value=self.site_bucket.bucket_name,
                      description="S3 bucket for site content")
        cdk.CfnOutput(self, "CloudFrontDistributionId", value=self.distribution.distribution_id,
                      description="CloudFront distribution ID for cache invalidation")
        cdk.CfnOutput(self, "CloudFrontDomainName", value=self.distribution.distribution_domain_name,
                      description="CloudFront domain (*.cloudfront.net)")
        cdk.CfnOutput(self, "SiteUrl", value=f"https://{domain}",
                      description="Primary site URL")
