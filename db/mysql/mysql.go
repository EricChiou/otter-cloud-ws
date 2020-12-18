package mysql

import (
	"database/sql"
	"strings"

	"otter-cloud-ws/config"
	"otter-cloud-ws/constants/api"
)

// DB mysql connecting
var DB *sql.DB

// Page QueryPage input struct
type Page struct {
	TableName   string
	ColumnNames []string
	UniqueKey   string
	OrderBy     string
	Page        int
	Limit       int
}

// Init connect MySQL
func Init() (err error) {
	for _, c := range specificCharStr {
		specificChar[int(c)] = true
	}

	cfg := config.Get()
	userName := cfg.MySQLUserName
	password := cfg.MySQLPassword
	addr := cfg.MySQLAddr
	port := cfg.MySQLPort
	dbName := cfg.MySQLDBNAME

	DB, err = sql.Open("mysql", userName+":"+password+"@tcp("+addr+":"+port+")/"+dbName)
	return err
}

// Close close mysql connecting
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// ErrMsgHandler error message handler
func ErrMsgHandler(err error) api.RespStatus {
	if strings.Contains(err.Error(), "Duplicate") {
		return api.Duplicate
	}
	return api.DBError
}
