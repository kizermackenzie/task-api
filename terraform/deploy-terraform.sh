#!/bin/bash

# Terraform Deployment Script for Task API
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID="tasksproject-464417"
REGION="us-central1"
IMAGE_NAME="gcr.io/${PROJECT_ID}/task-api:simple"

echo -e "${GREEN}🚀 Starting Terraform deployment for Task API${NC}"

# Check if Terraform is installed
if ! command -v terraform &> /dev/null; then
    echo -e "${RED}❌ Terraform is not installed. Please install it first.${NC}"
    echo "Install from: https://learn.hashicorp.com/tutorials/terraform/install-cli"
    exit 1
fi

# Check if gcloud is configured
if ! gcloud config list --filter=project:$PROJECT_ID --format="value(core.project)" | grep -q $PROJECT_ID; then
    echo -e "${YELLOW}⚠️  Setting GCP project to ${PROJECT_ID}${NC}"
    gcloud config set project $PROJECT_ID
fi

# Ensure Docker image exists
echo -e "${YELLOW}🔍 Checking if Docker image exists${NC}"
if ! gcloud container images describe $IMAGE_NAME &>/dev/null; then
    echo -e "${YELLOW}🏗️  Docker image not found. Building and pushing...${NC}"
    
    # Build for AMD64 architecture
    GOOS=linux GOARCH=amd64 go build -o main ../main.go
    
    # Create simple Dockerfile and build
    docker build -f ../Dockerfile.simple --platform linux/amd64 -t $IMAGE_NAME ..
    docker push $IMAGE_NAME
    
    echo -e "${GREEN}✅ Docker image built and pushed${NC}"
else
    echo -e "${GREEN}✅ Docker image already exists${NC}"
fi

# Navigate to terraform directory
cd "$(dirname "$0")"

# Create terraform.tfvars if it doesn't exist
if [ ! -f "terraform.tfvars" ]; then
    echo -e "${YELLOW}📝 Creating terraform.tfvars from example${NC}"
    cp terraform.tfvars.example terraform.tfvars
    echo -e "${BLUE}ℹ️  You can customize terraform.tfvars if needed${NC}"
fi

# Initialize Terraform
echo -e "${YELLOW}🔧 Initializing Terraform${NC}"
terraform init

# Validate Terraform configuration
echo -e "${YELLOW}✅ Validating Terraform configuration${NC}"
terraform validate

# Plan the deployment
echo -e "${YELLOW}📋 Creating Terraform plan${NC}"
terraform plan -out=tfplan

# Ask for confirmation
echo -e "${BLUE}❓ Do you want to apply this plan? (y/N)${NC}"
read -r response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    # Apply the plan
    echo -e "${YELLOW}🚀 Applying Terraform plan${NC}"
    terraform apply tfplan
    
    # Show outputs
    echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
    echo -e "${BLUE}📊 Deployment outputs:${NC}"
    terraform output
    
    # Test the deployment
    SERVICE_URL=$(terraform output -raw service_url)
    echo -e "${YELLOW}🧪 Testing the deployment${NC}"
    
    # Wait a moment for the service to be ready
    sleep 10
    
    if curl -s "$SERVICE_URL/health" > /dev/null; then
        echo -e "${GREEN}✅ Health check passed!${NC}"
        echo -e "${GREEN}🌐 Your API is available at: ${SERVICE_URL}${NC}"
        echo -e "${GREEN}🏥 Health check: ${SERVICE_URL}/health${NC}"
    else
        echo -e "${YELLOW}⚠️  Health check failed. The service might still be starting up.${NC}"
        echo -e "${BLUE}ℹ️  Check the logs with: gcloud run logs read task-api --region=$REGION${NC}"
    fi
    
    # Instructions for API access
    echo -e "${BLUE}📚 To make authenticated API calls:${NC}"
    echo "1. Use the generated service account key: ../api-caller-key.json"
    echo "2. Run: gcloud auth activate-service-account --key-file=../api-caller-key.json"
    echo "3. Get token: gcloud auth print-identity-token --audiences=$SERVICE_URL"
    echo "4. Use token in Authorization header: Bearer <token>"
    
    # Clean up
    rm -f tfplan
    
else
    echo -e "${YELLOW}❌ Deployment cancelled${NC}"
    rm -f tfplan
fi

echo -e "${GREEN}✨ Script completed${NC}"