# Project Configuration
variable "project_id" {
  description = "The GCP project ID"
  type        = string
  default     = "tasksproject-464417"
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "The GCP zone"
  type        = string
  default     = "us-central1-c"
}

# Application Configuration
variable "service_name" {
  description = "Name of the Cloud Run service"
  type        = string
  default     = "task-api"
}

variable "image_name" {
  description = "Docker image name"
  type        = string
  default     = "gcr.io/tasksproject-464417/task-api:simple"
}

# Database Configuration
variable "db_instance_name" {
  description = "Cloud SQL instance name"
  type        = string
  default     = "task-api-db"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "taskdb"
}

variable "db_user" {
  description = "Database user"
  type        = string
  default     = "postgres"
}

variable "db_tier" {
  description = "Database tier"
  type        = string
  default     = "db-f1-micro"
}

# Cloud Run Configuration
variable "cpu_limit" {
  description = "CPU limit for Cloud Run service"
  type        = string
  default     = "1"
}

variable "memory_limit" {
  description = "Memory limit for Cloud Run service"
  type        = string
  default     = "512Mi"
}

variable "max_instances" {
  description = "Maximum number of instances"
  type        = number
  default     = 20
}

variable "min_instances" {
  description = "Minimum number of instances"
  type        = number
  default     = 0
}

# Security Configuration
variable "allow_unauthenticated" {
  description = "Allow unauthenticated access to Cloud Run service"
  type        = bool
  default     = false
}