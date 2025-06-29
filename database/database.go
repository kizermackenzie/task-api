package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"task-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL_MODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	maxOpenConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnvWithDefault("DB_MAX_IDLE_CONNS", "5"))
	connMaxLifetime, _ := strconv.Atoi(getEnvWithDefault("DB_CONN_MAX_LIFETIME", "300"))

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("Database connected successfully")

	if err := db.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration completed successfully")

	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}