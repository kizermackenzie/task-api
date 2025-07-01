# Task API Terraform Infrastructure

This directory contains Terraform configuration files to deploy the Task API infrastructure on Google Cloud Platform.

## 📁 File Structure

```
terraform/
├── main.tf                    # Provider configuration and API enablement
├── variables.tf               # Input variables
├── database.tf               # Cloud SQL PostgreSQL instance
├── cloud_run.tf              # Cloud Run service configuration
├── iam.tf                    # Service accounts and IAM policies
├── secrets.tf                # Secret Manager for sensitive data
├── outputs.tf                # Output values
├── terraform.tfvars.example  # Example variables file
├── deploy-terraform.sh       # Automated deployment script
└── README.md                 # This file
```

## 🚀 Quick Start

### Prerequisites

1. **Install Terraform**: Download from [terraform.io](https://www.terraform.io/downloads.html)
2. **Install Google Cloud SDK**: `gcloud` CLI must be installed and authenticated
3. **Docker image**: Ensure your Docker image is built and pushed to GCR

### Deployment

1. **Run the automated script**:
   ```bash
   cd terraform
   ./deploy-terraform.sh
   ```

   Or manually:

2. **Copy variables file**:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars if needed
   ```

3. **Initialize Terraform**:
   ```bash
   terraform init
   ```

4. **Plan deployment**:
   ```bash
   terraform plan
   ```

5. **Apply configuration**:
   ```bash
   terraform apply
   ```

## 🏗️ Infrastructure Components

### Cloud SQL PostgreSQL
- **Instance**: `task-api-db`
- **Version**: PostgreSQL 14
- **Tier**: db-f1-micro (configurable)
- **Backup**: Daily at 3 AM UTC
- **SSL**: Disabled for Cloud Run connectivity

### Cloud Run Service
- **Name**: `task-api`
- **Memory**: 512Mi (configurable)
- **CPU**: 1 vCPU (configurable)
- **Scaling**: 0-20 instances (configurable)
- **Authentication**: Required by default

### Service Accounts
- **Cloud Run SA**: For database and secret access
- **API Caller SA**: For authenticated API access

### Secret Manager
- **Database Password**: Auto-generated secure password
- **JWT Secret**: Auto-generated secret for token signing

## 🔧 Configuration

### Variables

Key variables in `terraform.tfvars`:

```hcl
project_id            = "your-gcp-project-id"
region               = "us-central1"
service_name         = "task-api"
image_name           = "gcr.io/your-project/task-api:latest"
allow_unauthenticated = false
```

### Security

- **Database credentials**: Stored in Secret Manager
- **JWT secret**: Auto-generated and stored securely
- **Service accounts**: Minimal required permissions
- **Authentication**: Required by default for API access

## 📊 Outputs

After deployment, Terraform provides:

- **Service URL**: Cloud Run service endpoint
- **Database details**: Connection information
- **Service accounts**: Email addresses
- **API endpoints**: Ready-to-use URLs

## 🔐 API Access

### Authenticated Access

1. **Activate service account**:
   ```bash
   gcloud auth activate-service-account --key-file=../api-caller-key.json
   ```

2. **Get identity token**:
   ```bash
   TOKEN=$(gcloud auth print-identity-token --audiences=YOUR_SERVICE_URL)
   ```

3. **Make API calls**:
   ```bash
   curl -H "Authorization: Bearer $TOKEN" \\
        -X POST YOUR_SERVICE_URL/api/v1/auth/register \\
        -H "Content-Type: application/json" \\
        -d '{"email":"test@example.com","password":"Password123","first_name":"Test","last_name":"User"}'
   ```

## 🧹 Cleanup

To destroy the infrastructure:

```bash
terraform destroy
```

⚠️ **Warning**: This will permanently delete all resources including the database and data.

## 📝 Notes

- **State management**: Consider using remote state for production
- **Environment separation**: Use workspaces or separate directories for dev/staging/prod
- **Secrets rotation**: Implement regular rotation of database passwords and JWT secrets
- **Monitoring**: Add Cloud Monitoring and logging for production deployments

## 🔍 Troubleshooting

### Common Issues

1. **Permission errors**: Ensure your GCP account has necessary IAM roles
2. **Docker image not found**: Build and push the image first
3. **API enablement**: Script automatically enables required APIs
4. **Connection timeouts**: Check Cloud SQL and Cloud Run configurations

### Debugging

```bash
# Check Cloud Run logs
gcloud run logs read task-api --region=us-central1

# Check Terraform state
terraform show

# Validate configuration
terraform validate
```