package sqlp

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sqlTestsuite struct {
	suite.Suite

	driver string
	dsn    string

	db *sql.DB
}

func (s *sqlTestsuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	s.checkErr(err)

	s.db = db
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	s.checkErr(err)
	_, err = db.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (999, 'Tom', 18, 'Jerry')")
	s.checkErr(err)
}

func (s *sqlTestsuite) TearDownTest() {
	_, err := s.db.Exec("DELETE FROM test_model;")
	s.checkErr(err)
}

func (s *sqlTestsuite) TestCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	s.testInsertRow(ctx)
	s.testQueryRows(ctx)
	s.testUpdateRow(ctx)
	s.testQueryRow(ctx)

}

func (s *sqlTestsuite) testInsertRow(ctx context.Context) {
	db := s.db

	res, err := db.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (1, 'Tom', 18, 'Jerry')")
	s.checkErr(err)

	affectedRows, err := res.RowsAffected()
	s.checkErr(err)
	s.Assertions.Equal(int64(1), affectedRows)
}

func (s *sqlTestsuite) testQueryRows(ctx context.Context) {
	db := s.db
	rows, err := db.QueryContext(ctx,
		"SELECT `id`, `first_name`,`age`, `last_name` FROM `test_model` LIMIT ?", 1)
	s.checkErr(err)
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		tm := &TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		s.checkErr(err)
		s.Equal("Tom", tm.FirstName)
	}
}

func (s *sqlTestsuite) testUpdateRow(ctx context.Context) {
	db := s.db
	res, err := db.ExecContext(ctx,
		"UPDATE `test_model` SET `first_name` = 'changed' WHERE `id` = ?", 1)
	s.checkErr(err)
	affected, err := res.RowsAffected()
	s.checkErr(err)
	s.Equal(int64(1), affected)
}

func (s *sqlTestsuite) testQueryRow(ctx context.Context) {
	db := s.db
	row := db.QueryRowContext(ctx, "SELECT `id`, `first_name`,`age`, `last_name` FROM `test_model` WHERE id=?", 1)
	s.checkErr(row.Err())
	tm := &TestModel{}
	err := row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	s.checkErr(err)
	s.Assertions.Equal("changed", tm.FirstName)
}

func (s *sqlTestsuite) checkErr(err error) {
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &sqlTestsuite{
		driver: "sqlite3",
		dsn:    "file:test.db?cache=shared&mode=memory",
	})
}

type TestModel struct {
	Id        int64 `eorm:"auto_increment,primary_key"`
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
