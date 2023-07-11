package csv_importer

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

var mysqlKeyword = map[string]struct{}{
	"key":      {},
	"int":      {},
	"float":    {},
	"interval": {},
	"group":    {},
	"groups":   {},
	"add":      {},
	"desc":     {},
	"repeat":   {},
	"order":    {},
	"rank":     {},
	"default":  {},
	"show":     {},
}
var invalidChar = []string{" ", "-", ">=", "<=", "=", ">", "<", "(", ")"}

type Column struct {
	Name string
	Type string
}

type Header []Column

type Table struct {
	Name   string
	Header Header
	Data   [][]any
}

func (d Header) Columns() []string {
	h := make([]string, 0, len(d))
	for _, v := range d {
		h = append(h, v.Name)
	}
	return h
}
func (d Header) Types() []string {
	t := make([]string, 0, len(d))
	for _, v := range d {
		t = append(t, v.Type)
	}
	return t
}

const separator rune = ','

func Parse(file string) (*Table, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, nil
	}

	r := csv.NewReader(bytes.NewReader(b))
	r.Comma = separator

	filename := filepath.Base(file)
	extension := filepath.Ext(file)
	// 移除扩展名以获取不带后缀的文件名
	tableName := strings.TrimSuffix(filename, extension)

	var (
		tb = &Table{
			Name: strings.Replace(tableName, "-", "_", -1),
		}
	)

	record, err := r.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("header not enough")
		}
		return nil, err
	}
	tb.Header = make([]Column, len(record))

	for i, s := range record {
		tb.Header[i].Name = strings.TrimSpace(s)
		if len(s) == 0 {
			tb.Header[i].Name = fmt.Sprintf("Filed_%d", i)
		}
		tb.Header[i].Name = validateName(tb.Header[i].Name)
		if isKeyword(tb.Header[i].Name) {
			tb.Header[i].Name += "_"
		}
		if len(tb.Header[i].Name) == 0 {
			tb.Header[i].Name = fmt.Sprintf("Filed_%d", i)
		}
	}

	var (
		firstLine []string
		idx       int
	)
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return tb, err
		}
		if idx == 0 {
			firstLine = record
		}
		idx++
		if len(record) == 0 {
			continue
		}
		var d []interface{}
		for _, r := range record {
			d = append(d, r)
		}
		tb.Data = append(tb.Data, d)
	}
	// 根据第一行分析类型
	for i, s := range firstLine {
		_, err := strconv.ParseFloat(s, 64)
		if err == nil {
			tb.Header[i].Type = "float"
		} else {
			tb.Header[i].Type = "text"
		}
	}

	return tb, nil
}

func validateName(name string) string {
	for _, s := range invalidChar {
		name = strings.ReplaceAll(name, s, "_")
	}
	return strings.Trim(name, "_")
}

func isKeyword(name string) bool {
	_, ok := mysqlKeyword[strings.ToLower(name)]
	return ok
}
