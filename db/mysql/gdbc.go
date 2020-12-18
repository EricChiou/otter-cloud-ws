package mysql

import (
	"database/sql"
	"errors"
	"strconv"
)

// Gdbc sql orm
type Gdbc struct{}

var specificCharStr string = `"':.,;(){}[]&|=+-*%/\<>^`
var specificChar [128]bool

// Insert insert data
func (g *Gdbc) Insert(table string, columnValues map[string]interface{}) (sql.Result, error) {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, errors.New("db not initialized")
	}

	columnSQL := ""
	valueSQL := ""
	var args []interface{}
	for k, v := range columnValues {
		columnSQL += ", " + k
		valueSQL += ", ?"
		args = append(args, v)
	}
	if len(columnSQL) > 2 {
		columnSQL = columnSQL[2:]
	}
	if len(valueSQL) > 2 {
		valueSQL = valueSQL[2:]
	}

	return tx.Exec("INSERT INTO "+table+"( "+columnSQL+" ) VALUES( "+valueSQL+" )", args...)
}

// Exec execute sql
func (g *Gdbc) Exec(sql string, params SQLParams, args ...interface{}) (sql.Result, error) {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, errors.New("db not initialized")
	}

	execSQL := convertSQL(sql, params.kv)
	return tx.Exec(execSQL, args...)
}

// Query query rows
func (g *Gdbc) Query(sql string, params SQLParams, rowMapper func(rows *sql.Rows) error, args ...interface{}) error {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return errors.New("db not initialized")
	}

	execSQL := convertSQL(sql, params.kv)

	rows, err := tx.Query(execSQL, args...)
	defer rows.Close()
	if err != nil {
		return err
	}

	return rowMapper(rows)
}

// QueryRow query one row
func (g *Gdbc) QueryRow(sql string, params SQLParams, rowMapper func(row *sql.Row) error, args ...interface{}) error {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return errors.New("db not initialized")
	}

	execSQL := convertSQL(sql, params.kv)
	row := tx.QueryRow(execSQL, args...)

	return rowMapper(row)
}

// QueryPage query page format
func (g *Gdbc) QueryPage(page Page, whereSQL string, rowMapper func(rows *sql.Rows, total int) error, args ...interface{}) error {
	tx, err := DB.Begin()
	defer tx.Commit()
	if err != nil {
		return errors.New("db not initialized")
	}

	var columnSQL string
	for _, columnName := range page.ColumnNames {
		columnSQL += ", " + columnName
	}
	if len(columnSQL) > 0 {
		columnSQL = columnSQL[2:]
	} else {
		columnSQL = "*"
	}

	sql := "SELECT " + columnSQL + " "
	sql += "FROM " + page.TableName + " "
	sql += "    JOIN ( "
	sql += "    SELECT " + page.UniqueKey + " FROM " + page.TableName + " "
	sql += "    ORDER BY " + page.OrderBy + " "
	sql += "    LIMIT " + strconv.Itoa((page.Page-1)*page.Limit) + ", " + strconv.Itoa(page.Limit) + " "
	sql += ") t "
	sql += "USING ( " + page.UniqueKey + " ) "
	sql += whereSQL

	rows, err := tx.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return err
	}

	var total int
	sql = "SELECT COUNT(*) FROM " + page.TableName + " " + whereSQL
	row := tx.QueryRow(sql, args...)
	row.Scan(&total)

	rowMapper(rows, total)
	return nil
}

func convertSQL(originalSQL string, params map[string]string) string {
	convertSQL := ""

	preIndex := 0
	for i := 0; i < len(originalSQL)-1; i++ {
		if originalSQL[i:i+1] == "#" {
			key := getKey(originalSQL, i+1)
			value := params[key]
			if len(value) > 0 {
				convertSQL += originalSQL[preIndex:i] + value
				i += len(key)
				preIndex = i + 1
			}
		}
	}
	convertSQL += originalSQL[preIndex:]

	return convertSQL
}

func getKey(original string, startIndex int) string {
	for j := startIndex; j < len(original); j++ {
		if isSpecificChar([]rune(original[j : j+1])[0]) {
			key := original[startIndex:j]
			return key
		}
	}

	return original[startIndex:]
}

func isSpecificChar(c rune) bool {
	return (c < 128 && specificChar[c]) || c == ' '
}
