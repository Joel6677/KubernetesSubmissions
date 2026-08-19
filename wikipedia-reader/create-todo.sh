#!/usr/bin/env bash
set -e

if [ -z "$TODO_BACKEND_URL" ]; then
  echo "TODO_BACKEND_URL is not set"
  exit 1
fi

LOCATION=$(curl -s -o /dev/null -D - "https://en.wikipedia.org/wiki/Special:Random" | grep -i "^location:" | awk '{print $2}' | tr -d '\r')

if [ -z "$LOCATION" ]; then
  echo "Failed to get location"
  exit 1
fi

if [[ "$LOCATION" == //* ]]; then
  LOCATION="https:$LOCATION"
fi

if [[ "$LOCATION" != http* ]]; then
  LOCATION="https://en.wikipedia.org$LOCATION"
fi

TEXT="Read $LOCATION"
echo "Creating todo: $TEXT"

RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$TODO_BACKEND_URL" \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"$TEXT\"}")

if [ "$RESPONSE" != "201" ]; then
  echo "Failed to create todo, got status $RESPONSE"
  exit 1
fi

echo "Todo created successfully"
