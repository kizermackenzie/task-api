# Enable Secret Manager API
resource "google_project_service" "secret_manager" {
  service = "secretmanager.googleapis.com"
}

# Secret for database password
resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password"
  
  replication {
    auto {}
  }

  depends_on = [google_project_service.secret_manager]
}

resource "google_secret_manager_secret_version" "db_password" {
  secret         = google_secret_manager_secret.db_password.id
  secret_data_wo = random_password.db_password.result
}

# Secret for JWT secret
resource "google_secret_manager_secret" "jwt_secret" {
  secret_id = "jwt-secret"
  
  replication {
    auto {}
  }

  depends_on = [google_project_service.secret_manager]
}

resource "google_secret_manager_secret_version" "jwt_secret" {
  secret         = google_secret_manager_secret.jwt_secret.id
  secret_data_wo = random_password.jwt_secret.result
}