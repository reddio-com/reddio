package database

import (
	"database/sql"
	"testing"
	"time"

	"gorm.io/gorm"
)

var MockConfig = &Config{
	DSN:        "testuser:123456@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
	DriverName: "mysql",
	MaxOpenNum: 10,
	MaxIdleNum: 5,
}

func MockPing(db *gorm.DB) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	return sqlDB, sqlDB.Ping()
}

func TestInitDB(t *testing.T) {
	db, err := InitDB(MockConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := CloseDB(db); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}
	time.Sleep(2 * time.Second)

	if sqlDB.Stats().MaxOpenConnections != MockConfig.MaxOpenNum {
		t.Errorf("Expected MaxOpenConnections to be %d, got %d", MockConfig.MaxOpenNum, sqlDB.Stats().MaxOpenConnections)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

type User_test struct {
	gorm.Model
	Name  string
	Email string
}

func TestCreateAndGetUser(t *testing.T) {
	db, err := InitDB(MockConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := CloseDB(db); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}()

	if err := db.AutoMigrate(&User_test{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	user := User_test{Name: "John Doe", Email: "john.doe@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	var retrievedUser User_test
	if err := db.First(&retrievedUser, user.ID).Error; err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrievedUser.Name != user.Name {
		t.Errorf("Expected user name to be %s, got %s", user.Name, retrievedUser.Name)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected user email to be %s, got %s", user.Email, retrievedUser.Email)
	}
}
