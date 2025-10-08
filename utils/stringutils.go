package utils

import (
	"database/sql"
	"fmt"
)

type String struct {
	Pointer *string
	Like    bool
}

func (s String) ToPgText() sql.NullString {
	if s.Pointer != nil && *s.Pointer != "" {
		if s.Like {
			return sql.NullString{String: fmt.Sprintf("%%%s%%", *s.Pointer), Valid: true}
		}
		return sql.NullString{String: *s.Pointer, Valid: true}
	}
	return sql.NullString{}
}
