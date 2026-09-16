#!/usr/bin/env bash
set -e

if [ -z "$DATABASE_URL" ]; then
  echo "DATABASE_URL is not set"
  exit 1
fi

if [ -z "$BUCKET_NAME" ]; then
  echo "BUCKET_NAME is not set"
  exit 1
fi

TIMESTAMP=$(date +%Y-%m-%dT%H-%M-%S)
BACKUP_FILE="/tmp/backup-${TIMESTAMP}.sql"

pg_dump -v "$DATABASE_URL" > "$BACKUP_FILE"

gsutil cp "$BACKUP_FILE" "gs://${BUCKET_NAME}/backup-${TIMESTAMP}.sql"

echo "Backup uploaded to gs://${BUCKET_NAME}/backup-${TIMESTAMP}.sql"
