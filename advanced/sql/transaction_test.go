package sqlp

import (
	"context"
	"database/sql"
	"time"
)

func (s *sqlTestsuite) TestTransaction() {
	t := s.T()
	ctx, cancal := context.WithTimeout(context.Background(), time.Second)
	defer cancal()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	s.checkErr(err)

	err = s.multiInsert(ctx, tx)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	rows, err := s.db.QueryContext(ctx, "SELECT `id`, `first_name`,`age`, `last_name` FROM `test_model` WHERE id=?", 100)
	s.checkErr(err)
	if !rows.Next() {
		return
	}
	tm := &TestModel{}
	err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	s.checkErr(err)
	t.Log(tm)
}

func (s *sqlTestsuite) multiInsert(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (100, 'Tom', 20, 'Jerry')")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (999, 'Tom', 20, 'Jerry')")
	if err != nil {
		return err
	}
	return nil
}
