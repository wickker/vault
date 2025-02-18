package utils

import "github.com/jackc/pgx/v5/pgtype"

type String struct {
	StringPointer *string
}

func (s String) PointerToPgText() pgtype.Text {
	if s.StringPointer != nil && *s.StringPointer != "" {
		return pgtype.Text{String: *s.StringPointer, Valid: true}
	}
	return pgtype.Text{}
}
