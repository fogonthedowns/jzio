# ── Domain ───────────────────────────────────────────────────────────────────
DOMAIN = "jz.io"

# ── AWS account ───────────────────────────────────────────────────────────────
AWS_ACCOUNT = "919607760751"

# Leave empty for a CDK-generated bucket name, or set a fixed name e.g. "jz.io".
SITE_BUCKET_NAME = ""

LOG_BUCKET_NAME = "jzio-logs"

# ── DNS TTL (seconds) ─────────────────────────────────────────────────────────
TTL = 300

# ── FastMail MX ───────────────────────────────────────────────────────────────
FASTMAIL_MX = [
    {"priority": 10, "host": "in1-smtp.messagingengine.com"},
    {"priority": 20, "host": "in2-smtp.messagingengine.com"},
]

FASTMAIL_SPF = "v=spf1 include:spf.messagingengine.com ?all"

# DKIM key prefixes — each becomes <key>._domainkey CNAME -> <key>.<domain>.dkim.fmhosted.com
FASTMAIL_DKIM_KEYS = ["fm1", "fm2", "fm3", "mesmtp"]

# FastMail CalDAV/CardDAV/IMAP/SMTP service-discovery SRV records.
FASTMAIL_SRV = True

# ── Optional third-party DNS records ─────────────────────────────────────────
# Set to None to skip.

WEBCAM_CNAME = "ec2-18-144-70-98.us-west-1.compute.amazonaws.com"

GOOGLE_SITE_VERIFICATION = {
    "record_name": "wzexeppdb46o",
    "cname_target": "gv-ouol5yihh3pdq6.dv.googlehosted.com",
}

# SendGrid — record_name is relative to the domain.
# Set to [] to skip all SendGrid records.
SENDGRID_RECORDS = [
    {"record_name": "s1._domainkey", "cname_target": "s1.domainkey.u16015049.wl050.sendgrid.net"},
    {"record_name": "s2._domainkey", "cname_target": "s2.domainkey.u16015049.wl050.sendgrid.net"},
    {"record_name": "em5811", "cname_target": "u16015049.wl050.sendgrid.net"},
]
