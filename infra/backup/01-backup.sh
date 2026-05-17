#!/usr/bin/env bash
# Backs up all jz.io resources before CDK redeploy.
# Safe to run multiple times — nothing destructive.
set -euo pipefail

BACKUP_DIR="$(cd "$(dirname "$0")" && pwd)"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

echo "=== Backing up Route53 ==="
aws route53 list-resource-record-sets --hosted-zone-id Z3E1BI0B0MOT2U --output json \
  > "$BACKUP_DIR/jz.io-route53-$TIMESTAMP.json"
echo "  Saved to jz.io-route53-$TIMESTAMP.json"

echo ""
echo "=== Backing up S3 content (excluding .git/ and yolo-darknet/) ==="
S3_BACKUP="$BACKUP_DIR/s3-jz.io"
mkdir -p "$S3_BACKUP"
aws s3 sync s3://jz.io/ "$S3_BACKUP/" \
  --exclude ".git/*" \
  --exclude "yolo-darknet/*" \
  --exclude "yolo.html"
echo "  Saved to $S3_BACKUP"

echo ""
echo "=== Backup complete ==="
echo "  Route53: $BACKUP_DIR/jz.io-route53-$TIMESTAMP.json"
echo "  S3:      $S3_BACKUP"
