package sqlb

import (
	"errors"
	"strings"
)

func (db *Init) Delete() error {
	db.query = []string{}
	defer db.Clear()
	var err error

	if db.queryBuilder.from == "" {
		//table not init
		return errors.New("table not found")
	}
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return err
		}
	}

	db.query = append(db.query, "DELETE FROM "+db.queryBuilder.from)
	db.joinBuild()
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

func (db *QueryInit) Delete() (queryRaw string, err error) {
	var dbs Init
	dbs.queryBuilder = *db
	dbs.query = []string{}
	if db.from == "" {
		//table not init
		err = errors.New("table not found")
		return
	}

	dbs.query = append(dbs.query, "DELETE FROM "+db.from)
	dbs.joinBuild()
	dbs.whereBuild()
	queryRaw = strings.Join(dbs.query, "\n")
	return
}
