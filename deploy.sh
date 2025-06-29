#!/bin/bash

# GCP Deployment Script for Task API
set -e

# Configuration
PROJECT_ID="tasksproject-464417"
REGION="us-central1"
SERVICE_NAME="task-api"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting GCP deployment for Task API${NC}"

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}❌ gcloud CLI is not installed. Please install it first.${NC}"
    exit 1
fi

# Set the project
echo -e "${YELLOW}📋 Setting GCP project to ${PROJECT_ID}${NC}"
gcloud config set project $PROJECT_ID

# Enable required APIs
echo -e "${YELLOW}🔧 Enabling required GCP APIs${NC}"
gcloud services enable run.googleapis.com
gcloud services enable sqladmin.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Build the Docker image
echo -e "${YELLOW}🏗️  Building Docker image${NC}"
docker build -t $IMAGE_NAME .

# Push to Google Container Registry
echo -e "${YELLOW}📤 Pushing image to Google Container Registry${NC}"
docker push $IMAGE_NAME

# Create Cloud SQL instance (if it doesn't exist)
echo -e "${YELLOW}🗄️  Setting up Cloud SQL PostgreSQL instance${NC}"
if ! gcloud sql instances describe task-api-db --region=$REGION &> /dev/null; then
    gcloud sql instances create task-api-db \
        --database-version=POSTGRES_14 \
        --tier=db-f1-micro \
        --region=$REGION \
        --storage-auto-increase \
        --backup-start-time=03:00
    
    # Set postgres user password
    echo -e "${YELLOW}🔐 Setting database password${NC}"
    gcloud sql users set-password postgres \
        --instance=task-api-db \
        --password="TaskAPI2025!"
    
    # Create the database
    gcloud sql databases create taskdb --instance=task-api-db
else
    echo -e "${GREEN}✅ Cloud SQL instance already exists${NC}"
fi

# Get the private IP of the Cloud SQL instance
DB_HOST=$(gcloud sql instances describe task-api-db --format="value(ipAddresses[0].ipAddress)")

# Deploy to Cloud Run
echo -e "${YELLOW}🌐 Deploying to Cloud Run${NC}"
gcloud run deploy $SERVICE_NAME \
    --image=$IMAGE_NAME \
    --platform=managed \
    --region=$REGION \
    --allow-unauthenticated \
    --port=8080 \
    --memory=512Mi \
    --cpu=1 \
    --set-env-vars="DB_HOST=${DB_HOST},DB_PORT=5432,DB_USER=postgres,DB_PASSWORD=TaskAPI2025!,DB_NAME=taskdb,DB_SSL_MODE=require,JWT_SECRET=TaskAPI-JWT-Secret-2025-Production,PORT=8080,GIN_MODE=release"

# Get the service URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")

echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo -e "${GREEN}📱 Your API is available at: ${SERVICE_URL}${NC}"
echo -e "${GREEN}🏥 Health check: ${SERVICE_URL}/health${NC}"
echo -e "${GREEN}📚 API docs: ${SERVICE_URL}/api/v1/auth/register${NC}"

# Test the deployment
echo -e "${YELLOW}🧪 Testing the deployment${NC}"
if curl -s "${SERVICE_URL}/health" > /dev/null; then
    echo -e "${GREEN}✅ Health check passed!${NC}"
else
    echo -e "${RED}❌ Health check failed. Check the logs:${NC}"
    echo "gcloud run logs read --service=$SERVICE_NAME --region=$REGION"
fi