package database_test

import (
	"fmt"
	"log"
	"os"
	"share-basket-server/core/config"
	"share-basket-server/core/db"
	"share-basket-server/core/util"
	"sync"
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

const migrationsPath = "../../migrations"

var (
	once   sync.Once
	testDB *gorm.DB

	cfg = config.DB{
		Port:     util.GetEnv("DB_PORT", "5432"),
		Host:     util.GetEnv("TEST_DB_HOST", "test_db"),
		Name:     util.GetEnv("POSTGRES_DB", "share-basket"),
		User:     util.GetEnv("POSTGRES_USER", "postgres"),
		Password: util.GetEnv("POSTGRES_PASSWORD", "postgres"),
	}
)

func TestMain(m *testing.M) {
	setupTestDB()

	code := m.Run()

	teardownTestDB()

	os.Exit(code)
}

func setupTestDB() {
	once.Do(func() {
		var err error

		testDB, err = db.New(cfg)
		if err != nil {
			log.Fatalf("initialize test database: %v", err)
		}

		sqlDB, err := testDB.DB()
		if err != nil {
			log.Fatalf("failed to get test database: %v", err)
		}

		if err := db.Migrate(sqlDB, migrationsPath); err != nil {
			log.Fatalf("migrate test database: %v", err)
		}

	})
}

func clearTestData() {
	tables := []string{"users", "accounts"}
	for _, table := range tables {
		err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)).Error
		if err != nil {
			log.Fatalf("failed to clear test data from %s: %v", table, err)
		}
	}

	fmt.Println("Test data cleared successfully")
}

func teardownTestDB() {
	sqlDB, err := testDB.DB()
	if err != nil {
		log.Fatalf("get test database: %v", err)
	}

	defer sqlDB.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	n, err := migrate.ExecMax(sqlDB, "postgres", migrations, migrate.Down, -1)
	if err != nil {
		log.Fatalf("Failed to apply down migration: %v", err)
	}

	fmt.Printf("Rolled back %d migration(s)\n", n)
}
