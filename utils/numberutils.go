package utils

import (
	"database/sql"
)

type Number struct {
	Pointer *int32
}

//func (n Number) ToPgInt4() pgtype.Int4 {
//	if n.Pointer != nil && *n.Pointer != 0 {
//		return pgtype.Int4{
//			Int32: *n.Pointer,
//			Valid: true,
//		}
//	}
//	return pgtype.Int4{}
//}

func (n Number) ToPgInt4() sql.NullInt32 {
	if n.Pointer != nil && *n.Pointer != 0 {
		return sql.NullInt32{
			Int32: *n.Pointer,
			Valid: true,
		}
	}
	return sql.NullInt32{}
}
