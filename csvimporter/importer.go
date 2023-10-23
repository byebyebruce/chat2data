package csvimporter

import (
	"database/sql"
	"fmt"
	"strings"
)

func isNumber(t string) bool {
	switch t {
	case "int", "int32", "int64":
		return true
	case "float", "float32", "float64":
		return true
	default:
		return false
	}
}
func type2SQLType(t string) string {
	switch t {
	case "int", "int32", "int64":
		return "INT"
	case "float", "float32", "float64":
		return "FLOAT"
	//case "string":
	//return "VARCHAR(256)"
	default:
		return "TEXT"
		//return "VARCHAR(256)"
	}
}

func ImportCSV2DB(db *sql.DB, file string) (string, error) {
	t, err := Parse(file)
	if err != nil {
		return "", err
	}
	err = ImportTable(db, t)
	return t.Name, err
}

func ImportTable(db *sql.DB, tb *Table) error {
	table := tb.Name
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	if err != nil {
		return err
	}

	cs, ts := tb.Header.Columns(), tb.Header.Types()
	var as []string
	for i, c := range cs {
		as = append(as, fmt.Sprintf("%s %s", c, type2SQLType(ts[i])))
	}

	createStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s)", table, strings.Join(as, ",\n\t"))
	_, err = db.Exec(createStmt)
	if err != nil {
		return err
	}

	if len(tb.Data) == 0 {
		return nil
	}
	qs := strings.Repeat("?,", len(tb.Header))
	qs = qs[:len(qs)-1]

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (%s) values (%s)", table, strings.Join(cs, ","), qs))
	if err != nil {
		return err
	}

	for _, d := range tb.Data {
		for i, a := range d {
			if v, ok := a.(string); ok && len(v) == 0 {
				if isNumber(tb.Header[i].Type) {
					d[i] = 0
				}
			}
		}
		_, err = stmt.Exec(d...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type Schema struct {
	Cid       interface{} `json:"cid"`
	Name      string      `json:"name"`
	ColType   string      `json:"colType"`
	Notnull   interface{} `json:"notnull"`
	DfltValue interface{} `json:"dflt_Value"`
	Pk        interface{} `json:"pk"`
}

func SQLiteTableInfo(db *sql.DB, table string) ([]Schema, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, err
	}

	var ss []Schema
	for rows.Next() {
		var s Schema
		err = rows.Scan(&s.Cid, &s.Name, &s.ColType, &s.Notnull, &s.DfltValue, &s.Pk)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, nil
}

type QueryResult []string

func (q QueryResult) Convert2AnyArray() []any {
	var v = make([]any, len(q))
	for i := 0; i < len(q); i++ {
		v[i] = &q[i]
	}
	return v
}

func ShowTables(db *sql.DB) ([]string, error) {
	r, err := Query(db, "SHOW TABLES;")
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, result := range r {
		tbName := ""
		if len(result) > 0 {
			tbName = result[0]
		}
		if len(tbName) > 0 {
			ret = append(ret, tbName)
		}
	}
	return ret, nil
}

func DropAllTables(db *sql.DB, tbs []string) error {
	for _, table := range tbs {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			return err
		}
	}
	return nil
}

func Query(db *sql.DB, query string) ([]QueryResult, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var a []QueryResult
	for rows.Next() {
		var v QueryResult = make([]string, len(cols))
		rows.ColumnTypes()
		if err := rows.Scan(v.Convert2AnyArray()...); err != nil {
			panic(err)
		}
		a = append(a, v)
	}

	return a, nil
}
