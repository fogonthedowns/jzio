#!/usr/bin/env bash
# Deletes old jz.io infrastructure so CDK can create it fresh.
# Run ONLY after 01-backup.sh has completed successfully.
# EMAIL WILL BE DOWN from the moment Route53 records are deleted until
# `cdk deploy` finishes creating them (~15-20 minutes).
set -euo pipefail

echo "=== Step 1: Disable and delete old CloudFront distribution E9KIROWWL1OJ2 ==="
echo "  Fetching current config..."
ETAG=$(aws cloudfront get-distribution-config --id E9KIROWWL1OJ2 --query 'ETag' --output text)
aws cloudfront get-distribution-config --id E9KIROWWL1OJ2 --query 'DistributionConfig' --output json \
  | python3 -c "
import sys, json
cfg = json.load(sys.stdin)
cfg['Enabled'] = False
cfg['Aliases'] = {'Quantity': 0, 'Items': []}
print(json.dumps(cfg))
" > /tmp/cf-disabled.json

echo "  Disabling distribution and removing jz.io alias..."
NEW_ETAG=$(aws cloudfront update-distribution \
  --id E9KIROWWL1OJ2 \
  --distribution-config file:///tmp/cf-disabled.json \
  --if-match "$ETAG" \
  --query 'ETag' --output text)
echo "  Waiting for distribution to fully disable (can take 5-10 minutes)..."
aws cloudfront wait distribution-deployed --id E9KIROWWL1OJ2
echo "  Deleting distribution..."
aws cloudfront delete-distribution --id E9KIROWWL1OJ2 --if-match "$NEW_ETAG"
echo "  CloudFront E9KIROWWL1OJ2 deleted."

echo ""
echo "=== Step 2: Delete S3 bucket jz.io ==="
echo "  Removing all objects..."
aws s3 rm s3://jz.io/ --recursive
echo "  Deleting bucket..."
aws s3api delete-bucket --bucket jz.io --region us-west-1
echo "  jz.io bucket deleted."

echo ""
echo "=== Step 3: Delete S3 bucket www.jz.io ==="
aws s3 rm s3://www.jz.io/ --recursive 2>/dev/null || true
aws s3api delete-bucket --bucket www.jz.io --region us-west-1
echo "  www.jz.io bucket deleted."

echo ""
echo "=== Step 4: Delete conflicting Route53 records ==="
echo "  (All records except NS, SOA, and old cert validation CNAME)"
aws route53 change-resource-record-sets \
  --hosted-zone-id Z3E1BI0B0MOT2U \
  --change-batch file://$(dirname "$0")/route53-deletes.json
echo "  Route53 records deleted. EMAIL IS NOW DOWN."

echo ""
echo "=== Done. Run: cd infra && cdk deploy ==="
echo "  Email will be restored once cdk deploy completes (~15-20 minutes)."
