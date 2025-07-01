# Cloud Run Service URL
output "service_url" {
  description = "URL of the Cloud Run service"
  value       = google_cloud_run_service.task_api.status[0].url
}

# Database Connection Details
output "database_connection_name" {
  description = "Connection name for Cloud SQL instance"
  value       = google_sql_database_instance.postgres.connection_name
}

output "database_private_ip" {
  description = "Private IP address of the Cloud SQL instance"
  value       = google_sql_database_instance.postgres.private_ip_address
}

output "database_public_ip" {
  description = "Public IP address of the Cloud SQL instance"
  value       = google_sql_database_instance.postgres.public_ip_address
}

# Service Account Information
output "cloud_run_service_account" {
  description = "Email of the Cloud Run service account"
  value       = google_service_account.cloud_run_sa.email
}

output "api_caller_service_account" {
  description = "Email of the API caller service account"
  value       = google_service_account.api_caller.email
}

# Secret Information
output "database_password_secret" {
  description = "Secret Manager secret name for database password"
  value       = google_secret_manager_secret.db_password.secret_id
  sensitive   = true
}

output "jwt_secret_name" {
  description = "Secret Manager secret name for JWT secret"
  value       = google_secret_manager_secret.jwt_secret.secret_id
  sensitive   = true
}

# Project Information
output "project_id" {
  description = "The GCP project ID"
  value       = var.project_id
}

output "region" {
  description = "The GCP region"
  value       = var.region
}

# Health Check URL
output "health_check_url" {
  description = "Health check endpoint URL"
  value       = "${google_cloud_run_service.task_api.status[0].url}/health"
}

# API Endpoints
output "api_endpoints" {
  description = "Main API endpoints"
  value = {
    register = "${google_cloud_run_service.task_api.status[0].url}/api/v1/auth/register"
    login    = "${google_cloud_run_service.task_api.status[0].url}/api/v1/auth/login"
    tasks    = "${google_cloud_run_service.task_api.status[0].url}/api/v1/tasks"
    health   = "${google_cloud_run_service.task_api.status[0].url}/health"
  }
}