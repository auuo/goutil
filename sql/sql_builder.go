package sql

import (
	"strconv"
	"strings"
)

type Builder struct {
	aliasNum int

	selectClause []string
	baseQuery    *Query
	joinQuery    []*Join
	whereClause  [][]string // 一维之间使用 or 连接, 二维之间使用 and 连接
}

type Query struct {
	// sql, table 二选一
	sql   string
	table string
	alias string
}

type Join struct {
	Query
	joinType string
	on       string
	sb       *Builder
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (sb *Builder) NewAlias() string {
	sb.aliasNum++
	return "t" + strconv.Itoa(sb.aliasNum)
}

func (sb *Builder) Select(names ...string) *Builder {
	sb.selectClause = append(sb.selectClause, names...)
	return sb
}

func (sb *Builder) From(sql, alias string) *Builder {
	q := Query{
		sql:   sql,
		alias: alias,
	}
	sb.baseQuery = &q
	return sb
}

func (sb *Builder) FromTable(table, alias string) *Builder {
	q := Query{
		table: table,
		alias: alias,
	}
	sb.baseQuery = &q
	return sb
}

func (sb *Builder) LeftJoin(query, alias string) *Join {
	return sb.join("left", "", query, alias)
}

func (sb *Builder) LeftJoinTable(table, alias string) *Join {
	return sb.join("left", table, "", alias)
}

func (sb *Builder) InnerJoin(query, alias string) *Join {
	return sb.join("inner", "", query, alias)
}

func (sb *Builder) InnerJoinTable(table, alias string) *Join {
	return sb.join("left", table, "", alias)
}

func (sb *Builder) join(joinType, table, query, alias string) *Join {
	join := Join{
		Query: Query{
			sql:   query,
			table: table,
			alias: alias,
		},
		on:       "",
		sb:       sb,
		joinType: joinType,
	}
	sb.joinQuery = append(sb.joinQuery, &join)
	return &join
}

func (sb *Builder) Where(sql string) *Builder {
	if len(sb.whereClause) == 0 {
		sb.whereClause = append(sb.whereClause, []string{})
	}
	last := sb.whereClause[len(sb.whereClause)-1]
	sb.whereClause[len(sb.whereClause)-1] = append(last, sql)
	return sb
}

func (sb *Builder) Or(sql string) *Builder {
	if len(sb.whereClause) == 0 {
		sb.whereClause = append(sb.whereClause, []string{})
	}
	last := sb.whereClause[len(sb.whereClause)-1]
	if len(last) != 0 {
		sb.whereClause = append(sb.whereClause, []string{})
	}
	last = sb.whereClause[len(sb.whereClause)-1]
	sb.whereClause[len(sb.whereClause)-1] = append(last, sql)
	return sb
}

func (j *Join) On(sql string) *Builder {
	j.on = sql
	return j.sb
}

func (sb *Builder) Build() string {
	sql := "select "
	for i, c := range sb.selectClause {
		if i != 0 {
			sql += ", "
		}
		sql += c
	}
	sql += " from " + sb.baseQuery.build()
	for _, q := range sb.joinQuery {
		sql += " " + q.joinType + " join "
		sql += q.build()
		sql += " on " + q.on
	}
	return sql + buildWhere(sb.whereClause)
}

func (q *Query) build() string {
	if q.table != "" {
		return q.table + " as " + q.alias
	}
	return "(" + q.sql + ") as " + q.alias
}

func buildWhere(whereClause [][]string) string {
	if len(whereClause) == 0 {
		return ""
	}
	var andGroup []string
	for _, s := range whereClause {
		andGroup = append(andGroup, buildAndGroup(s))
	}
	return " where " + buildOrGroup(andGroup)
}

func buildOrGroup(group []string) string {
	s := strings.Join(group, ") or (")
	if len(group) > 1 {
		s = "(" + s + ")"
	}
	return s
}

func buildAndGroup(clause []string) string {
	return strings.Join(clause, " and ")
}
