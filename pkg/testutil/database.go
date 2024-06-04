package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func InitTxdbDatabase(t *testing.T) (*sql.DB, error) {
	t.Helper()

	err := godotenv.Load("../../.env")
	if err != nil {
		panic("error loading .env file")
	}

	connString := fmt.Sprintf("%s?parseTime=true", os.Getenv("MYSQL_CONNECTION_STRING"))

	txdb.Register("txdb", "mysql", connString)
	db, err := sql.Open("txdb", uuid.New().String())
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
