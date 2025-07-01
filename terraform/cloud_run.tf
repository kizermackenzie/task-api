# Cloud Run Service
resource "google_cloud_run_service" "task_api" {
  name     = var.service_name
  location = var.region

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"         = var.max_instances
        "autoscaling.knative.dev/minScale"         = var.min_instances
        "run.googleapis.com/cloudsql-instances"    = google_sql_database_instance.postgres.connection_name
        "run.googleapis.com/cpu-throttling"        = "false"
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }

    spec {
      service_account_name = google_service_account.cloud_run_sa.email
      
      containers {
        image = var.image_name

        ports {
          container_port = 8080
        }

        resources {
          limits = {
            cpu    = var.cpu_limit
            memory = var.memory_limit
          }
        }

        env {
          name  = "DB_HOST"
          value = "/cloudsql/${google_sql_database_instance.postgres.connection_name}"
        }

        env {
          name  = "DB_PORT"
          value = "5432"
        }

        env {
          name  = "DB_USER"
          value = var.db_user
        }

        env {
          name = "DB_PASSWORD"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret_version.db_password.secret
              key  = "latest"
            }
          }
        }

        env {
          name  = "DB_NAME"
          value = var.db_name
        }

        env {
          name  = "DB_SSL_MODE"
          value = "disable"
        }

        env {
          name = "JWT_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret_version.jwt_secret.secret
              key  = "latest"
            }
          }
        }

        env {
          name  = "GIN_MODE"
          value = "release"
        }

        startup_probe {
          initial_delay_seconds = 0
          timeout_seconds       = 240
          period_seconds        = 240
          failure_threshold     = 1
          tcp_socket {
            port = 8080
          }
        }

        liveness_probe {
          http_get {
            path = "/health"
            port = 8080
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.cloud_run,
    google_sql_database_instance.postgres,
    google_secret_manager_secret_version.db_password,
    google_secret_manager_secret_version.jwt_secret
  ]
}

# IAM policy for Cloud Run service
resource "google_cloud_run_service_iam_member" "public_access" {
  count = var.allow_unauthenticated ? 1 : 0

  service  = google_cloud_run_service.task_api.name
  location = google_cloud_run_service.task_api.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# IAM policy for service account access
resource "google_cloud_run_service_iam_member" "service_account_access" {
  service  = google_cloud_run_service.task_api.name
  location = google_cloud_run_service.task_api.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.api_caller.email}"
}