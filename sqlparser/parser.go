// Package sqlparser  用于解析 SQL 构建缓存定位标识符
package sqlparser

import (
	"fmt"
	sp "github.com/xwb1989/sqlparser"
	"strings"
)

type SelectOption bool

const (
	ALL      SelectOption = true
	DISTINCT SelectOption = false
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
	raw          string
	stmt         sp.Statement
	SelectOption SelectOption
	SelectExpr   string
	hashKey      string
	// TODO
}

// HashKey TODO 更好的 Hash 算法
func (i Identifier) init() {
	i.hashKey = i.raw
}

func (i Identifier) HashKey() string {
	return i.hashKey
}
