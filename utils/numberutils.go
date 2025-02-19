package utils

import "github.com/jackc/pgx/v5/pgtype"

type Number struct {
	Pointer *int32
}

func (n Number) ToPgInt4() pgtype.Int4 {
	if n.Pointer != nil && *n.Pointer != 0 {
		return pgtype.Int4{
			Int32: *n.Pointer,
			Valid: true,
		}
	}
	return pgtype.Int4{}
}
