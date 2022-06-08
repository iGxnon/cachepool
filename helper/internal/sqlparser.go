package internal

import (
	"fmt"
	sp "github.com/xwb1989/sqlparser"
	"strings"
)

type SelectOption bool

const (
	ALL      SelectOption = true
	DISTINCT SelectOption = false
	Format                = "sqlkey::{table}:{column}:{}"
)

// Parse TODO 解析更多语句信息
func Parse(query string, args ...any) (Identifier, error) {
	i := 0
	for strings.Contains(query, "?") {
		query = strings.Replace(query, "?", fmt.Sprint(args[i]), 1)
		i++
	}
	stmt, err := sp.Parse(query)
	if err != nil {
		return Identifier{}, err
	}
	id := Identifier{raw: query, stmt: stmt}
	id.init()
	return id, nil
}

type Identifier struct {
	raw  string
	stmt sp.Statement

	key string
}

// Key table
func (i *Identifier) init() {
	i.key = strings.ReplaceAll(i.raw, " ", "-")
}

func (i *Identifier) Key() string {
	return i.key
}
