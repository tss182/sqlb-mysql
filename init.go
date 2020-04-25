package sqlb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strings"
)

type Init struct {
	ConnectionType    string //config, env
	Host              string
	User              string
	Password          string
	DbName            string
	sel               string
	from              string
	join              []string
	orderBy           []string
	groupBy           string
	limit             string
	having            string
	where             []whereDb
	query             []string
	removeSpecialChar bool
	call              bool
	dbs               *sql.DB
	transaction       *sql.Tx
}

type whereDb struct {
	and      bool
	field    string
	value    interface{}
	valueIn  []interface{}
	starDate string
	endDate  string
	op       string //'',in,notIn,startGroup,endGroup,between
}

func (db *Init) From(from string) *Init {
	db.from = from
	db.call = false
	return db
}

func (db *Init) mysqlConnect() (*sql.DB, error) {
	var dbs *sql.DB
	var err error
	if db.ConnectionType == "env" {
		db.Host = os.Getenv("sqlHost")
		db.User = os.Getenv("sqlUser")
		db.Password = os.Getenv("sqlPassword")
		db.DbName = os.Getenv("sqlDb")
	}
	dbs, err = sql.Open("mysql", db.User+":"+db.Password+"@"+db.Host+"/"+db.DbName)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (db *Init) Close() {
	if db.dbs != nil {
		_ = db.dbs.Close()
	}
}

func (db *Init) Clear() {
	tx := db.transaction
	dbs := db.dbs
	p := reflect.ValueOf(db).Elem()
	p.Set(reflect.Zero(p.Type()))
	db.transaction = tx
	db.dbs = dbs
}

func (db *Init) QueryView() string {
	return strings.Join(db.query, "/n")
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
