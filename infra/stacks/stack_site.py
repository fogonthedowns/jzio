from typing import Optional

import aws_cdk as cdk
import aws_cdk.aws_athena as athena
import aws_cdk.aws_certificatemanager as acm
import aws_cdk.aws_cloudfront as cloudfront
import aws_cdk.aws_cloudfront_origins as origins
import aws_cdk.aws_glue as glue
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

        # ── CloudFront access logs bucket ─────────────────────────────────────
        log_bucket = s3.Bucket(
            self, "LogBucket",
            bucket_name="jzio-logs",
            block_public_access=s3.BlockPublicAccess.BLOCK_ALL,
            encryption=s3.BucketEncryption.S3_MANAGED,
            enforce_ssl=True,
            removal_policy=cdk.RemovalPolicy.RETAIN,
            auto_delete_objects=False,
            object_ownership=s3.ObjectOwnership.BUCKET_OWNER_PREFERRED,
            access_control=s3.BucketAccessControl.LOG_DELIVERY_WRITE,
            lifecycle_rules=[
                s3.LifecycleRule(expiration=cdk.Duration.days(90)),
            ],
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
            enable_logging=True,
            log_bucket=log_bucket,
            log_file_prefix="cf/",
        )

        # ── DNS records (apex + www) ──────────────────────────────────────────
        cf_target = route53.RecordTarget.from_alias(
            targets.CloudFrontTarget(self.distribution)
        )
        route53.ARecord(self, "ApexA", zone=hosted_zone, target=cf_target)
        route53.AaaaRecord(self, "ApexAAAA", zone=hosted_zone, target=cf_target)
        route53.ARecord(self, "WwwA", zone=hosted_zone, record_name="www", target=cf_target)
        route53.AaaaRecord(self, "WwwAAAA", zone=hosted_zone, record_name="www", target=cf_target)

        # ── Observability (Glue table + Athena named queries) ────────────────
        glue.CfnTable(
            self, "CloudFrontLogsTable",
            catalog_id=cdk.Aws.ACCOUNT_ID,
            database_name="default",
            table_input=glue.CfnTable.TableInputProperty(
                name="cloudfront_logs",
                table_type="EXTERNAL_TABLE",
                parameters={
                    "EXTERNAL": "TRUE",
                    "skip.header.line.count": "2",
                    "classification": "cloudfront",
                },
                storage_descriptor=glue.CfnTable.StorageDescriptorProperty(
                    location=f"s3://{log_bucket.bucket_name}/cf",
                    input_format="org.apache.hadoop.mapred.TextInputFormat",
                    output_format="org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat",
                    compressed=False,
                    stored_as_sub_directories=False,
                    serde_info=glue.CfnTable.SerdeInfoProperty(
                        serialization_library="org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe",
                        parameters={
                            "field.delim": "\t",
                            "line.delim": "\n",
                            "serialization.format": "\t",
                        },
                    ),
                    columns=[
                        glue.CfnTable.ColumnProperty(name="date",                       type="string"),
                        glue.CfnTable.ColumnProperty(name="time",                       type="string"),
                        glue.CfnTable.ColumnProperty(name="x_edge_location",            type="string"),
                        glue.CfnTable.ColumnProperty(name="bytes",                      type="bigint"),
                        glue.CfnTable.ColumnProperty(name="ip",                         type="string"),
                        glue.CfnTable.ColumnProperty(name="method",                     type="string"),
                        glue.CfnTable.ColumnProperty(name="host",                       type="string"),
                        glue.CfnTable.ColumnProperty(name="uri",                        type="string"),
                        glue.CfnTable.ColumnProperty(name="status",                     type="int"),
                        glue.CfnTable.ColumnProperty(name="referer",                    type="string"),
                        glue.CfnTable.ColumnProperty(name="user_agent",                 type="string"),
                        glue.CfnTable.ColumnProperty(name="query_string",               type="string"),
                        glue.CfnTable.ColumnProperty(name="cookie",                     type="string"),
                        glue.CfnTable.ColumnProperty(name="x_edge_result_type",         type="string"),
                        glue.CfnTable.ColumnProperty(name="x_edge_request_id",          type="string"),
                        glue.CfnTable.ColumnProperty(name="x_host_header",              type="string"),
                        glue.CfnTable.ColumnProperty(name="protocol",                   type="string"),
                        glue.CfnTable.ColumnProperty(name="bytes_out",                  type="bigint"),
                        glue.CfnTable.ColumnProperty(name="time_taken",                 type="float"),
                        glue.CfnTable.ColumnProperty(name="x_forwarded_for",            type="string"),
                        glue.CfnTable.ColumnProperty(name="ssl_protocol",               type="string"),
                        glue.CfnTable.ColumnProperty(name="ssl_cipher",                 type="string"),
                        glue.CfnTable.ColumnProperty(name="x_edge_response_result_type", type="string"),
                        glue.CfnTable.ColumnProperty(name="protocol_version",           type="string"),
                    ],
                ),
            ),
        )

        # Workgroup with a 1 GB per-query scan limit — at $5/TB this caps a single
        # runaway query at $0.005 and keeps the monthly bill well under $1.
        workgroup = athena.CfnWorkGroup(
            self, "AthenaWorkGroup",
            name="jzio",
            state="ENABLED",
            work_group_configuration=athena.CfnWorkGroup.WorkGroupConfigurationProperty(
                enforce_work_group_configuration=True,
                bytes_scanned_cutoff_per_query=1_073_741_824,  # 1 GB
                result_configuration=athena.CfnWorkGroup.ResultConfigurationProperty(
                    output_location=f"s3://{log_bucket.bucket_name}/athena/",
                ),
                publish_cloud_watch_metrics_enabled=False,
            ),
        )

        _named_queries = [
            ("AthenaQueryAccessLog", "Access log", "default",
             "This looks like a server access log",
             "SELECT\n"
             "    CONCAT(date, ' ', time) as timestamp,\n"
             "    ip,\n"
             "    method,\n"
             "    uri,\n"
             "    status,\n"
             "    bytes,\n"
             "    x_edge_result_type,\n"
             "    x_edge_location\n"
             "  FROM cloudfront_logs\n"
             "  ORDER BY date DESC, time DESC\n"
             "  LIMIT 100;"),
            ("AthenaQueryStatusCodeDist", "Status Code Dist", "default",
             None,
             "SELECT status, COUNT(*) as count\n"
             "FROM cloudfront_logs\n"
             "GROUP BY status\n"
             "ORDER BY count DESC;"),
            ("AthenaQueryTopIPs", "Top IPs", "default",
             None,
             "SELECT ip, COUNT(*) as requests\n"
             "  FROM cloudfront_logs\n"
             "  GROUP BY ip\n"
             "  ORDER BY requests DESC\n"
             "  LIMIT 10;"),
            ("AthenaQueryTopPaths", "Top paths", "default",
             None,
             "SELECT uri, COUNT(*) as hits, AVG(CAST(bytes AS DOUBLE)) as avg_bytes\n"
             "  FROM cloudfront_logs\n"
             "  GROUP BY uri\n"
             "  ORDER BY hits DESC\n"
             "  LIMIT 20;"),
        ]
        for logical_id, name, database, description, query in _named_queries:
            q = athena.CfnNamedQuery(
                self, logical_id,
                name=name,
                database=database,
                query_string=query,
                work_group=workgroup.name,
                **({"description": description} if description else {}),
            )
            q.add_dependency(workgroup)

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
        cdk.CfnOutput(self, "LogBucketName", value=log_bucket.bucket_name,
                      description="S3 bucket for CloudFront access logs")
