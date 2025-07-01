# Service Account for Cloud Run
resource "google_service_account" "cloud_run_sa" {
  account_id   = "${var.service_name}-sa"
  display_name = "Service Account for ${var.service_name}"
  description  = "Service account used by Cloud Run service for ${var.service_name}"
}

# Service Account for API access
resource "google_service_account" "api_caller" {
  account_id   = "api-caller"
  display_name = "API Caller"
  description  = "Service account for authenticated API access"
}

# Grant Cloud SQL Client role to Cloud Run service account
resource "google_project_iam_member" "cloud_run_sql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

# Grant Secret Manager Secret Accessor role to Cloud Run service account
resource "google_project_iam_member" "cloud_run_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

# Create service account key for API caller
resource "google_service_account_key" "api_caller_key" {
  service_account_id = google_service_account.api_caller.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

# Save the service account key to a local file
resource "local_file" "api_caller_key" {
  content  = base64decode(google_service_account_key.api_caller_key.private_key)
  filename = "${path.root}/../api-caller-key.json"
  
  provisioner "local-exec" {
    command = "chmod 600 ${path.root}/../api-caller-key.json"
  }
}