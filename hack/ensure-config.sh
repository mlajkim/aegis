#!/usr/bin/env bash
set -e

# Configuration file for this operator:
CONFIG_DIR="."
CONFIG_FILE="$CONFIG_DIR/config.yaml"
SAMPLE_FILE="./config.sample.yaml"

# Assume we do not need to create one if it already exists:
if [ -f "$CONFIG_FILE" ]; then
  exit 0
fi

echo "----------------------------------------------------------------"
echo "⚠️  Configuration file for this operator is NOT found at $CONFIG_FILE"
echo "🚀  Let's generate one for local development!"
echo "----------------------------------------------------------------"

# Create directory
mkdir -p "$CONFIG_DIR"

# 1. Get user input (with default values)
read -p "ZMS URL? [Hit ENTER for default: https://localhost:4443/zms/v1]: " INPUT_ZMS
ZMS_URL=${INPUT_ZMS:-"https://localhost:4443/zms/v1"}

echo ""
default_cert_path="$(pwd)/certs/athenz_admin.cert.pem"
default_key_path="$(pwd)/keys/athenz_admin.private.pem"

read -p "X.509 Cert File to connect to Athenz Server? [Hit ENTER for default: $default_cert_path]: " INPUT_CERT_PATH
CERT_PATH=${INPUT_CERT_PATH:-"$default_cert_path"}

read -p "X.509 Key File to connect to Athenz Server? [Hit ENTER for default: $default_key_path]: " INPUT_KEY_PATH
KEY_PATH=${INPUT_KEY_PATH:-"$default_key_path"}

if [ -z "$CERT_PATH" ] || [ -z "$KEY_PATH" ]; then
  echo "❌ Error: CertPath and KeyPath are required!"
  exit 1
fi

# Creates a basic config file, based on the sample yaml:

sed -e "s|{{ZMS_URL}}|$ZMS_URL|g" \
    -e "s|{{CERT_PATH}}|$CERT_PATH|g" \
    -e "s|{{KEY_PATH}}|$KEY_PATH|g" \
    "$SAMPLE_FILE" > "$CONFIG_FILE"

echo "----------------------------------------------------------------"
echo "✅  Config generated at: $CONFIG_FILE"
echo "----------------------------------------------------------------"