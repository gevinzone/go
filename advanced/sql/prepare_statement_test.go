package sqlp

import (
	"context"
	"database/sql"
)

func (s *sqlTestsuite) TestPrepareStatement() {
	stmt, err := s.db.Prepare("SELECT * FROM `test_model` WHERE `id` = ?")
	s.checkErr(err)

	rows, err := stmt.QueryContext(context.Background(), 999)
	s.checkErr(err)
	s.mapRows(rows)

	rows, err = stmt.QueryContext(context.Background(), 2)
	s.checkErr(err)
	s.mapRows(rows)

	err = stmt.Close()
	s.checkErr(err)
}

func (s *sqlTestsuite) mapRows(rows *sql.Rows) {
	for rows.Next() {
		tm := &TestModel{}
		err := rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		s.checkErr(err)
		s.T().Log(tm)
	}
}
