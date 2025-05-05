package repository_test

import (
	"fmt"
	"log"
	"os"
	"share-basket-server/core/config"
	"share-basket-server/core/util"
	"share-basket-server/infra/rdb/db"
	"share-basket-server/infra/rdb/repository"
	"sync"
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

const migrationsPath = "../../../migrations"

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

	dummyUser = repository.UserDto{
		ID:         "dummy-user-id",
		CognitoUID: "dummy-cognito-uid",
		Email:      "dummy@example.com",
	}

	dummyAccount = repository.AccountDto{
		ID:     "dummy-account-id",
		UserID: dummyUser.ID,
		Name:   "dummy user",
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
	tables := []string{"users", "accounts", "personal_shopping_items"}
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

func createDummyAccount() error {
	if err := testDB.Create(&dummyUser).Error; err != nil {
		return err
	}

	if err := testDB.Create(&dummyAccount).Error; err != nil {
		return err
	}

	return nil
}
