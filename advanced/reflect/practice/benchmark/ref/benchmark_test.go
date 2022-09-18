package homework

import (
	"database/sql"
	"testing"
)

func BenchmarkInsertStmt(b *testing.B) {
	entity := User{
		BaseEntity: BaseEntity{
			CreateTime: 123,
			UpdateTime: ptrInt64(456),
		},
		Id:       789,
		NickName: sql.NullString{String: "Tom", Valid: true},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = InsertStmt(entity)
	}
}
