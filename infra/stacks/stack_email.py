import aws_cdk as cdk
import aws_cdk.aws_route53 as route53
from constructs import Construct

from config import (
    FASTMAIL_DKIM_KEYS,
    FASTMAIL_MX,
    FASTMAIL_SPF,
    FASTMAIL_SRV,
    GOOGLE_SITE_VERIFICATION,
    SENDGRID_RECORDS,
    TTL,
    WEBCAM_CNAME,
)


class EmailStack(cdk.Stack):
    """Stack B — Email + third-party DNS records.

    Depends on DnsStack. Edit config.py to enable/disable individual record groups.
    """

    def __init__(
        self,
        scope: Construct,
        id: str,
        hosted_zone: route53.IHostedZone,
        **kwargs,
    ):
        super().__init__(scope, id, **kwargs)

        domain = hosted_zone.zone_name
        ttl = cdk.Duration.seconds(TTL)

        # ── FastMail MX ───────────────────────────────────────────────────────
        mx_values = [
            route53.MxRecordValue(priority=r["priority"], host_name=r["host"])
            for r in FASTMAIL_MX
        ]
        route53.MxRecord(self, "ApexMx", zone=hosted_zone, ttl=ttl, values=mx_values)
        route53.MxRecord(
            self, "WildcardMx", zone=hosted_zone, record_name="*", ttl=ttl, values=mx_values
        )

        # ── SPF ───────────────────────────────────────────────────────────────
        route53.TxtRecord(
            self, "ApexSpf", zone=hosted_zone, ttl=ttl, values=[FASTMAIL_SPF]
        )

        # ── FastMail DKIM CNAMEs ──────────────────────────────────────────────
        for key in FASTMAIL_DKIM_KEYS:
            safe = key.replace("-", "_")
            route53.CnameRecord(
                self, f"Dkim_{safe}",
                zone=hosted_zone,
                record_name=f"{key}._domainkey",
                ttl=ttl,
                domain_name=f"{key}.{domain}.dkim.fmhosted.com",
            )

        # ── FastMail SRV service-discovery ───────────────────────────────────
        if FASTMAIL_SRV:
            srv_entries = [
                ("SrvPop3s",     "_pop3s._tcp",      ["10 1 995 pop.fastmail.com"]),
                ("SrvCaldav",    "_caldav._tcp",     ["0 0 0 ."]),
                ("SrvCaldavs",   "_caldavs._tcp",    ["0 1 443 caldav.fastmail.com"]),
                ("SrvCarddav",   "_carddav._tcp",    ["0 0 0 ."]),
                ("SrvCarddavs",  "_carddavs._tcp",   ["0 1 443 carddav.fastmail.com"]),
                ("SrvImap",      "_imap._tcp",       ["0 0 0 ."]),
                ("SrvImaps",     "_imaps._tcp",      ["0 1 993 imap.fastmail.com"]),
                ("SrvPop3",      "_pop3._tcp",       ["0 0 0 ."]),
                ("SrvSubmission","_submission._tcp", ["0 1 587 smtp.fastmail.com"]),
            ]
            for res_id, name, values in srv_entries:
                route53.CfnRecordSet(
                    self, res_id,
                    hosted_zone_id=hosted_zone.hosted_zone_id,
                    name=f"{name}.{domain}.",
                    type="SRV",
                    ttl=str(TTL),
                    resource_records=values,
                )

        # ── SendGrid ─────────────────────────────────────────────────────────
        for i, rec in enumerate(SENDGRID_RECORDS):
            route53.CnameRecord(
                self, f"Sendgrid{i}",
                zone=hosted_zone,
                record_name=rec["record_name"],
                ttl=ttl,
                domain_name=rec["cname_target"],
            )

        # ── Google site verification ──────────────────────────────────────────
        if GOOGLE_SITE_VERIFICATION:
            route53.CnameRecord(
                self, "GoogleVerify",
                zone=hosted_zone,
                record_name=GOOGLE_SITE_VERIFICATION["record_name"],
                ttl=ttl,
                domain_name=GOOGLE_SITE_VERIFICATION["cname_target"],
            )

        # ── Webcam ────────────────────────────────────────────────────────────
        if WEBCAM_CNAME:
            route53.CnameRecord(
                self, "Webcam",
                zone=hosted_zone,
                record_name="webcam",
                ttl=ttl,
                domain_name=WEBCAM_CNAME,
            )
