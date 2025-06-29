#!/bin/bash

# Script to call the secured Task API using service account authentication
set -e

# Configuration
SERVICE_URL="https://task-api-15006307884.us-central1.run.app"
KEY_FILE="/Users/mackenziekizer/task-api-key.json"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔐 Calling secured Task API${NC}"

# Authenticate with service account
echo -e "${YELLOW}📋 Authenticating with service account${NC}"
gcloud auth activate-service-account --key-file="$KEY_FILE" --quiet

# Get identity token
TOKEN=$(gcloud auth print-identity-token --audiences="$SERVICE_URL")

echo -e "${YELLOW}🧪 Testing API endpoints${NC}"

# Test health endpoint
echo "Testing /health:"
curl -s -H "Authorization: Bearer $TOKEN" "$SERVICE_URL/health" | jq '.'

echo -e "\n"

# Test registration
echo "Testing /api/v1/auth/register:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -X POST "$SERVICE_URL/api/v1/auth/register" \
     -d '{
       "email": "automated@example.com",
       "password": "Password123",
       "first_name": "Auto",
       "last_name": "User"
     }' | jq '.'

# Switch back to user account
echo -e "\n${YELLOW}🔄 Switching back to user account${NC}"
gcloud config set account kizermackenzie@gmail.com --quiet

echo -e "${GREEN}✅ Secure API call completed successfully!${NC}"