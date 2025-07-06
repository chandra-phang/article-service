package db_client

import (
	"article-service/configloader"
	"article-service/db/transaction"
	"article-service/infrastructure/log"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
)

const (
	DbConnTimeout  = 300 // In seconds
	DbMaxIdleConns = 10
	DbMaxOpenConns = 20
)

var connPoolSingleton ISQLClient

var (
	errNoPendingMigrations = errors.New("no change")
)

// InitDatabase initializes the DB connection
func InitDatabase(ctx context.Context, config configloader.DbConfig) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName)

	// NOTE: configloader.DbConfig contains credentials!
	// DO NOT include it in any log, print, or panic!
	log.Infof(ctx, "[DB_Client] Initializing database")

	dbConn, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client] Failed to open DB")
		panic(err.Error())
	}

	err = dbConn.Ping()
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client] Failed to ping DB")
		panic(err.Error())
	}

	dbConn.SetMaxIdleConns(DbMaxIdleConns)
	dbConn.SetMaxOpenConns(DbMaxOpenConns)
	dbConn.SetConnMaxLifetime(DbConnTimeout * time.Second)

	connPoolSingleton = &dbClient{pool: dbConn}
}

// InitDatabaseMock sets up a mock database connection.
func InitDatabaseMock() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Errorf(context.Background(), err, "[DB_Client][InitDatabaseMock] Error: %s", err.Error())
	}

	connPoolSingleton = &dbClient{pool: db}
	return mock
}

// InitDatabaseMockExpectingTxn sets up a mock database connection.
func InitDatabaseMockExpectingTxn(ctrl *gomock.Controller, expectStartSuccess bool, expectCommitSuccess interface{}) sqlmock.Sqlmock {
	ctx := context.Background()
	_, mock, err := sqlmock.New()
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][InitDatabaseMockExpectingTxn] Error: %s", err.Error())
	}

	connPoolSingleton = &txnImpl{ctrl: ctrl, expectStartSuccess: expectStartSuccess, expectCommitSuccess: expectCommitSuccess}
	return mock
}

func GetDB() ISQLClient {
	if connPoolSingleton == nil {
		// This can happen in 2 cases:
		//   1. In a deployment, InitDB() was not called (unlikely)
		//   2. In a unit test, InitDBMock() or InitDBMockExpectingTxn() was not called (more likely)
		//
		// In the unit test case, we have to call InitDBMock or InitDBMockExpectingTxn(),
		// otherwise `GetClient().Begin()` will panic.
		panic("Database connection pool not initialized. Forgot to InitDBMock() or InitDBMockExpectingTxn()?")
	}
	return connPoolSingleton
}

// RunMigrations initializes the DB connection
func RunMigrations(ctx context.Context, config configloader.DbConfig) {
	connInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.DbName)

	// NOTE: configloader.DbConfig contains credentials!
	// DO NOT include it in any log, print, or panic!
	log.Infof(ctx, "[DB_Client] Running migrations")

	m, err := migrate.New(
		"file://db/migrations",
		connInfo)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client] Failed to init migrations")
		panic(err.Error())
	}
	if err := m.Up(); err != nil {
		if err.Error() == errNoPendingMigrations.Error() {
			log.Infof(ctx, "[DB_Client] No pending migrations to run")
		} else {
			log.Errorf(ctx, err, "[DB_Client] Failed to run migrations")
			panic(err.Error())
		}
	} else {
		log.Infof(ctx, "[DB_Client] Running migrations done")
	}
}

func CreateTestDB(ctx context.Context, config configloader.DbConfig) error {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, "postgres")

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][CreateTestDB] failed to connect to postgres")
		return err
	}
	defer db.Close()

	// CREATE DATABASE
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", config.TestDbName))
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][CreateTestDB] create db %s failed", config.TestDbName)
		return err
	}

	testDbConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.TestDbName)
	testDB, err := sql.Open("postgres", testDbConn)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][CreateTestDB] failed to connect to %s", config.TestDbName)
		return err
	}

	log.Infof(ctx, "[DB_Client][CreateTestDB] database %s created successfully", config.TestDbName)

	if err = MigrateTestDB(ctx, config); err != nil {
		return err
	}

	testDB.SetMaxIdleConns(DbMaxIdleConns)
	testDB.SetMaxOpenConns(DbMaxOpenConns)
	testDB.SetConnMaxLifetime(DbConnTimeout * time.Second)

	connPoolSingleton = &dbClient{pool: testDB}

	return nil
}

func MigrateTestDB(ctx context.Context, config configloader.DbConfig) error {
	connMigrate := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.TestDbName)

	m, err := migrate.New(
		"file://../db/migrations",
		connMigrate)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][MigrateTestDB] Failed to init migrations")
		return err
	}
	if err := m.Up(); err != nil {
		if err.Error() == errNoPendingMigrations.Error() {
			log.Infof(ctx, "[DB_Client][MigrateTestDB] No pending migrations to run")
		} else {
			log.Errorf(ctx, err, "[DB_Client][MigrateTestDB] Failed to run migrations")
			return err
		}
	} else {
		log.Infof(ctx, "[DB_Client][MigrateTestDB] Running migrations done")
	}

	return nil
}

func DropTestDB(ctx context.Context, config configloader.DbConfig) error {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, "postgres")

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][DropTestDB] failed to connect to postgres")
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1 AND pid <> pg_backend_pid()
	`, config.TestDbName)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][DropTestDB] failed to terminate connections, db: %s", config.TestDbName)
	}

	// DROP DATABASE
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", config.TestDbName))
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][DropTestDB] drop db %s failed", config.TestDbName)
		return err
	}

	log.Infof(ctx, "[DB_Client][DropTestDB] database %s dropped successfully", config.TestDbName)

	return nil
}

func RunIntegrationTestSeed(ctx context.Context, config configloader.DbConfig, dir string) error {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.TestDbName)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][CreateTestDB] failed to connect to postgres")
		return err
	}
	defer db.Close()

	absPath, err := filepath.Abs(dir)
	if err != nil {
		log.Errorf(ctx, err, "[RunSeedFiles] unable to resolve seed dir")
		return err
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		log.Errorf(ctx, err, "[RunSeedFiles] unable to read seed files")
		return err
	}

	// Sort files alphabetically (e.g., 001.sql, 002.sql)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		path := filepath.Join(absPath, file.Name())
		sqlContent, err := os.ReadFile(path)
		if err != nil {
			log.Errorf(ctx, err, "[RunSeedFiles] error reading %s", file.Name())
			return err
		}

		if _, err := db.Exec(string(sqlContent)); err != nil {
			log.Errorf(ctx, err, "[RunSeedFiles] error executing %s", file.Name())
			return err
		}

		log.Infof(ctx, "Executed seed: %s\n", file.Name())
	}

	return nil
}

func TruncateTestDB(ctx context.Context, config configloader.DbConfig) error {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.TestDbName)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Errorf(ctx, err, "[DB_Client][CreateTestDB] failed to connect to postgres")
		return err
	}
	defer db.Close()

	db.Exec("TRUNCATE articles, authors RESTART IDENTITY CASCADE")

	return nil
}

func StartTransactionCtx(ctx context.Context) (context.Context, transaction.ITransaction, error) {
	db := GetDB()
	return transaction.Start(ctx, db)
}

// ISQLClient describes a SQL client.
type ISQLClient interface {
	transaction.Manager
	ISQLOperations
}

// ITransaction describes an interface of a database transaction.
type ITransaction interface {
	transaction.ITransaction
	ISQLOperations
}

// ISQLOperations describes methods for querying a database.
type ISQLOperations interface {
	// Select queries the database and writes the output into the destination type.
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// Update executes INSERT/UPDATE/DELETE statement on the database.
	Exec(ctx context.Context, query string, args ...interface{}) (IResult, error)
}

// IResult summarizes an executed SQL command.
type IResult interface {
	RowsAffected() (int64, error)
}
