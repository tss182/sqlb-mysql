package sqlb

import (
	. "database/sql"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
)

func update(db *Init, value []interface{}) error {
	defer db.Clear()
	var err error
	var stmt *Stmt
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return err
		}
	}
	if db.transaction != nil {
		stmt, err = db.transaction.Prepare(strings.Join(db.query, " "))
	} else {
		stmt, err = db.dbs.Prepare(strings.Join(db.query, " "))
	}
	//defer stmt.Close()
	_, err = stmt.Exec(value...)
	if err != nil {
		return err
	}
	stmt.Close()
	return nil
}

func (db *Init) Update(query map[string]interface{}) error {
	db.query = []string{}
	if db.queryBuilder.from == "" {
		//table not init
		return errors.New("table not found")
	}
	if query == nil || len(query) == 0 {
		return errors.New("query invalid")
	}

	db.query = append(db.query, "update "+db.queryBuilder.from+" set ")

	var set []string
	var value []interface{}
	for i, v := range query {
		set = append(set, i+"=?")
		if db.queryBuilder.removeSpecialChar {
			v = removeSpecialChar(v)
		}
		value = append(value, v)
	}

	db.joinBuild()                                      //join
	db.query = append(db.query, strings.Join(set, ",")) //update
	db.whereBuild()                                     //where

	return update(db, value)
}

func (db *Init) UpdateStruct(update interface{}) error {
	json := jsoniter.Config{EscapeHTML: true, TagKey: "sqlb", OnlyTaggedField: true}.Froze()
	r, err := json.Marshal(update)
	if err != nil {
		return err
	}
	var updateMap map[string]interface{}
	err = json.Unmarshal(r, &updateMap)
	if err != nil {
		return err
	}
	return db.Update(updateMap)
}

func (db *Init) UpdateBatch(query []map[string]interface{}, id string) error {
	db.query = []string{}
	id = strings.TrimSpace(id)
	if db.queryBuilder.from == "" {
		//table not init
		return errors.New("table not found")
	}
	if query == nil || len(query) == 0 {
		return errors.New("query invalid")
	}

	var set map[string][]string
	set = map[string][]string{}
	var value map[string][]interface{}
	value = map[string][]interface{}{}
	values := []interface{}{}
	whereIn := []string{}
	for i, v := range query {
		if v[id] == nil {
			return errors.New("primary key for update, not found")
		}
		var reflectValue = reflect.ValueOf(v[id])
		var valId string
		var idInt int
		switch reflectValue.Kind() {
		case reflect.String:
			valId = strings.TrimSpace(reflectValue.String())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			idInt = int(reflectValue.Uint())
			valId = strconv.Itoa(idInt)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			idInt = int(reflectValue.Int())
			valId = strconv.Itoa(idInt)
		}
		whereIn = append(whereIn, "'"+valId+"'")

		for i2, v2 := range v {
			if i2 == id {
				continue
			}
			if i == 0 {
				set[i2] = append(set[i2], i2+" = (CASE "+id+"\n")
			}
			set[i2] = append(set[i2], "WHEN '"+valId+"' THEN ?\n")
			if db.queryBuilder.removeSpecialChar {
				v2 = removeSpecialChar(v2)
			}
			value[i2] = append(value[i2], v2)
		}
	}

	db.query = append(db.query, "update "+db.queryBuilder.from) //update
	db.joinBuild()                                              //join
	db.query = append(db.query, " SET ")
	var setQuery []string
	for i, v := range set {
		setQuery = append(setQuery, strings.Join(v, "")+" END)")
		for _, v2 := range value[i] {
			values = append(values, v2)
		}
	}

	for i, v := range setQuery {
		if len(setQuery)-1 != i {
			v += ","
		}
		db.query = append(db.query, v)
	}
	db.query = append(db.query, "where "+id+" in("+strings.Join(whereIn, ",")+")")
	return update(db, values)
}

func (db *Init) UpdateBatchStruct(insert interface{}, id string) error {
	json := jsoniter.Config{EscapeHTML: true, TagKey: "sqlb", OnlyTaggedField: true}.Froze()
	r, err := json.Marshal(insert)
	if err != nil {
		return err
	}
	var updateMap []map[string]interface{}
	err = json.Unmarshal(r, &updateMap)
	if err != nil {
		return err
	}
	return db.UpdateBatch(updateMap, id)
}

func (db *QueryInit) Update(query map[string]interface{}) (queryRaw string, values []interface{}, err error) {
	var dbs Init
	dbs.queryBuilder = *db
	dbs.query = []string{}
	if db.from == "" {
		//table not init
		err = errors.New("table not found")
		return
	}
	if query == nil || len(query) == 0 {
		err = errors.New("query invalid")
		return
	}

	dbs.query = append(dbs.query, "update "+db.from+" set ")

	var set []string
	for i, v := range query {
		set = append(set, i+"=?")
		if db.removeSpecialChar {
			v = removeSpecialChar(v)
		}
		values = append(values, v)
	}

	dbs.joinBuild()                                       //join
	dbs.query = append(dbs.query, strings.Join(set, ",")) //update
	dbs.whereBuild()                                      //where

	queryRaw = strings.Join(dbs.query, "\n")
	return
}

func (db *QueryInit) UpdateBatch(query []map[string]interface{}, id string) (queryRaw string, values []interface{}, err error) {
	var dbs Init
	dbs.queryBuilder = *db
	dbs.query = []string{}
	id = strings.TrimSpace(id)
	if db.from == "" {
		//table not init
		err = errors.New("table not found")
		return
	}
	if query == nil || len(query) == 0 {
		err = errors.New("query invalid")
		return
	}

	var set map[string][]string
	set = map[string][]string{}
	var value map[string][]interface{}
	value = map[string][]interface{}{}
	var whereIn []string
	for i, v := range query {
		if v[id] == nil {
			err = errors.New("primary key for update, not found")
			return
		}
		var reflectValue = reflect.ValueOf(v[id])
		var valId string
		var idInt int
		switch reflectValue.Kind() {
		case reflect.String:
			valId = strings.TrimSpace(reflectValue.String())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			idInt = int(reflectValue.Uint())
			valId = strconv.Itoa(idInt)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			idInt = int(reflectValue.Int())
			valId = strconv.Itoa(idInt)
		}
		whereIn = append(whereIn, "'"+valId+"'")

		for i2, v2 := range v {
			if i2 == id {
				continue
			}
			if i == 0 {
				set[i2] = append(set[i2], i2+" = (CASE "+id+"\n")
			}
			set[i2] = append(set[i2], "WHEN '"+valId+"' THEN ?\n")
			if db.removeSpecialChar {
				v2 = removeSpecialChar(v2)
			}
			value[i2] = append(value[i2], v2)
		}
	}

	dbs.query = append(dbs.query, "update "+db.from) //update
	dbs.joinBuild()                                  //join
	dbs.query = append(dbs.query, " SET ")
	var setQuery []string
	for i, v := range set {
		setQuery = append(setQuery, strings.Join(v, "")+" END)")
		for _, v2 := range value[i] {
			values = append(values, v2)
		}
	}

	for i, v := range setQuery {
		if len(setQuery)-1 != i {
			v += ","
		}
		dbs.query = append(dbs.query, v)
	}
	dbs.query = append(dbs.query, "where "+id+" in("+strings.Join(whereIn, ",")+")")

	queryRaw = strings.Join(dbs.query, "\n")
	return
}
