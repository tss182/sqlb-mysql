package sqlb

import (
	"errors"
	"strings"
)

func (db *Init) Delete() error {
	var err error

	if db.from == "" {
		//table not init
		return errors.New("table not found")
	}
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return err
		}
	}

	db.query = append(db.query, "DELETE FROM "+db.from)
	db.whereBuild()
	if db.transaction != nil {
		_, err = db.transaction.Exec(strings.Join(db.query, " "))
	} else {
		_, err = db.dbs.Exec(strings.Join(db.query, " "))
	}

	if err != nil {
		return err
	}
	return nil
}
