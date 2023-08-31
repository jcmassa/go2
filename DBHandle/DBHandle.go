package dbhandle

import (
	"database/sql"
	"flag"
	"fmt"
	Config "rankapi/ConfigHelper"
)

// var (
// 	server   = flag.String("S", "DESKTOP-CFFJJ9N", "server_name[\\instance_name]")
// 	instance = flag.String("I", "SQLEXPRESS", "server_name[\\instance_name]")
// 	database = flag.String("d", "BDRESERVAS", "db_name")
// 	port     = flag.Int("p", 1433, "db_name")
// )

var (
	server   = flag.String("S", Config.ReadValue("SqlSrv"), "server_name[\\instance_name]")
	instance = flag.String("I", Config.ReadValue("Instance"), "server_name[\\instance_name]")
	database = flag.String("d", Config.ReadValue("Database"), "db_name")
	port     = flag.Int("p", 1433, "db_name")
)

func SetDBConnection() (*sql.DB, error) {
	connstring := fmt.Sprintf("sqlserver://za:za@%s/%s?database=%s", *server, *instance, *database)
	//fmt.Print(connstring)
	return sql.Open("sqlserver", connstring)
}

func RunCommand(comando string) error {
	db, err := SetDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(comando)
	if err != nil {
		return err
	}
	return nil
}

func GetScalarVal(comando string) (int, error) {
	var sclVal int
	db, err := SetDBConnection()
	if err != nil {
		return 0, err //Se debe revisar del otro lado si el error tiene valor
	}
	defer db.Close()
	row := db.QueryRow(comando)
	switch err := row.Scan(&sclVal); err {
	case nil:
		return sclVal, err
	default:
		return 0, err
	}
}
