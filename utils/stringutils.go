package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

type String struct {
	String string
}

func (s String) ToPgText() pgtype.Text {
	str := strings.ReplaceAll(s.String, "%", "")
	str = strings.TrimSpace(str)
	if str != "" {
		return pgtype.Text{String: s.String, Valid: true}
	}
	return pgtype.Text{}
}
