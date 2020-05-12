package sqlb

import (
	. "database/sql"
	"errors"
	"strings"
)

func insert(querySql string, value []interface{}, db *Init) (interface{}, error) {
	defer db.Clear()
	var err error
	var stmt *Stmt
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return nil, err
		}
	}
	db.query = append(db.query, querySql)
	if db.transaction != nil {
		stmt, err = db.transaction.Prepare(querySql)
	} else {
		stmt, err = db.dbs.Prepare(querySql)
	}
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(value...)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	err = stmt.Close()
	if err != nil {
		return nil, err
	}
	return id, nil
}
func (db *Init) Insert(query map[string]interface{}) (interface{}, error) {
	db.query = []string{}
	if db.from == "" {
		//table not init
		return nil, errors.New("table not found")
	}
	querySql := "INSERT INTO " + db.from
	tag := ""
	field := ""
	value := []interface{}{}
	for i, v := range query {
		tag += "?,"
		field += i + ","
		if db.removeSpecialChar {
			v = removeSpecialChar(v)
		}
		value = append(value, v)

	}
	tag = tag[0 : len(tag)-1]
	field = field[0 : len(field)-1]
	querySql += "(" + field + ") values " + "(" + tag + ")"

	return insert(querySql, value, db)

}

func (db *Init) InsertBatch(query []map[string]interface{}) (interface{}, error) {
	db.query = []string{}
	if db.from == "" {
		//table not init
		return nil, errors.New("table not found")
	}
	querySql := "INSERT INTO " + db.from
	var value []interface{}
	field := joinMapKey(query[0], ",")
	fieldArray := strings.Split(field, ",")
	tag := strings.Repeat("?,", len(fieldArray))
	tag = tag[0 : len(tag)-1]
	tag = "(" + tag + ")"
	tags := strings.Repeat(tag+",", len(query))
	tags = tags[0 : len(tags)-1]
	for _, v := range query {
		for _, v2 := range fieldArray {
			if db.removeSpecialChar {
				v[v2] = removeSpecialChar(v[v2])
			}
			value = append(value, v[v2])
		}
	}

	querySql += "(" + field + ") values " + tags

	return insert(querySql, value, db)
}
