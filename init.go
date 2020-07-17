package sqlb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strings"
)

type Init struct {
	ConnectionType string //config, env
	connection     struct {
		host     string
		user     string
		password string
		dbName   string
	}
	queryBuilder QueryInit
	query        []string
	dbs          *sql.DB
	transaction  *sql.Tx
}

type QueryInit struct {
	sel               string
	from              string
	join              []string
	orderBy           []string
	groupBy           string
	limit             string
	having            string
	where             []whereDb
	removeSpecialChar bool
	call              bool
}

type whereDb struct {
	and      bool
	field    string
	value    interface{}
	sqlRaw   string
	starDate string
	endDate  string
	op       string //'',in,notIn,startGroup,endGroup,between
}

func (db *QueryInit) RemoveSpecialChar() {
	db.removeSpecialChar = true
}

func (db *Init) Setup(host, user, password, dbname string) {
	db.connection.host = host
	db.connection.user = user
	db.connection.password = password
	db.connection.dbName = dbname
}

func (db *QueryInit) From(from string) *QueryInit {
	db.from = from
	db.call = false
	return db
}

func (db *Init) mysqlConnect() (*sql.DB, error) {
	var dbs *sql.DB
	var err error
	if db.ConnectionType == "env" {
		db.connection.host = os.Getenv("sqlHost")
		db.connection.user = os.Getenv("sqlUser")
		db.connection.password = os.Getenv("sqlPassword")
		db.connection.dbName = os.Getenv("sqlDb")
	}
	dns := fmt.Sprintf("%s:%s@%s/%s", db.connection.user, db.connection.password, db.connection.host, db.connection.dbName)
	dbs, err = sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (db *Init) Close() {
	if db.dbs != nil && db.transaction == nil {
		_ = db.dbs.Close()
		db.dbs = nil
	}
}

func (db *Init) Clear() {
	tx := db.transaction
	dbs := db.dbs
	query := db.query
	connection := db.connection
	p := reflect.ValueOf(db).Elem()
	p.Set(reflect.Zero(p.Type()))
	db.transaction = tx
	db.dbs = dbs
	db.query = query
	db.connection = connection
}

func (db *QueryInit) Clear() {
	p := reflect.ValueOf(db).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (db *Init) Query(query QueryInit) *Init {
	db.queryBuilder = query
	return db
}

//func (db *QueryInit) QueryBackup() QueryInit {
//	var backup QueryInit
//	backup.sel = db.sel
//	backup.from = db.from
//	backup.join = db.join
//	backup.orderBy = db.orderBy
//	backup.groupBy = db.groupBy
//	backup.limit = db.limit
//	backup.having = db.having
//	backup.where = db.where
//	return backup
//}
//
//func (db *QueryInit) QueryRestore(restore QueryInit) {
//	db.sel = restore.sel
//	db.from = restore.from
//	db.join = restore.join
//	db.orderBy = restore.orderBy
//	db.groupBy = restore.groupBy
//	db.limit = restore.limit
//	db.having = restore.having
//	db.where = restore.where
//}

func (db *Init) QueryView() string {
	return strings.Join(db.query, "\n")
}

func (db *Init) Transaction() error {
	var err error
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return err
		}
	}
	tx, err := db.dbs.Begin()
	db.transaction = tx
	return err
}

func (db *Init) Rollback() error {
	var err error
	if db.transaction != nil {
		err = db.transaction.Rollback()
	}
	db.transaction = nil
	return err
}

func (db *Init) Commit() error {
	var err error
	if db.transaction != nil {
		err = db.transaction.Commit()
	}
	db.transaction = nil
	return err
}
