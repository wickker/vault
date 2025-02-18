package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type String struct {
	Pointer *string
	Like    bool
}

func (s String) ToPgText() pgtype.Text {
	if s.Pointer != nil && *s.Pointer != "" {
		if s.Like {
			return pgtype.Text{String: fmt.Sprintf("%%%s%%", *s.Pointer), Valid: true}
		}
		return pgtype.Text{String: *s.Pointer, Valid: true}
	}
	return pgtype.Text{}
}
