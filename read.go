package sqlb

import (
	. "database/sql"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"html"
	"strconv"
	"strings"
)

func (db *QueryInit) Select(sel string) *QueryInit {
	db.sel = sel
	return db
}

func (db *QueryInit) Join(table, on, join string) *QueryInit {
	join = strings.ToLower(join)
	joinStr := join + " join " + table + " on " + on
	db.join = append(db.join, joinStr)
	return db
}

func (db *QueryInit) Where(field string, value interface{}) *QueryInit {
	db.where = append(db.where, whereDb{and: true, field: field, value: value})
	return db
}

func (db *QueryInit) WhereOr(field string, value interface{}) *QueryInit {
	db.where = append(db.where, whereDb{and: false, field: field, value: value})
	return db
}

func (db *QueryInit) WhereIn(field string, value interface{}) *QueryInit {
	db.where = append(db.where, whereDb{and: true, field: field, value: value, op: "in"})
	return db
}

func (db *QueryInit) WhereInOr(field string, value interface{}) *QueryInit {
	db.where = append(db.where, whereDb{and: false, field: field, value: value, op: "in"})
	return db
}

func (db *QueryInit) WhereNotIn(field string, value interface{}) *QueryInit {
	db.where = append(db.where, whereDb{and: true, field: field, value: value, op: "notIn"})
	return db
}

func (db *QueryInit) WhereNotInOr(field string, value interface{}) *QueryInit {
	db.where = append(db.where, whereDb{and: false, field: field, value: value, op: "notIn"})
	return db
}

func (db *QueryInit) WhereBetween(field, startDate, endDate string) *QueryInit {
	db.where = append(db.where, whereDb{and: true, field: field, starDate: startDate, endDate: endDate, op: "between"})
	return db
}

func (db *QueryInit) WhereBetweenOr(field, startDate, endDate string) *QueryInit {
	db.where = append(db.where, whereDb{and: false, field: field, starDate: startDate, endDate: endDate, op: "between"})
	return db
}

func (db *QueryInit) WhereRaw(sql string) *QueryInit {
	db.where = append(db.where, whereDb{and: true, sqlRaw: sql, op: "raw"})
	return db
}

func (db *QueryInit) WhereRawOr(sql string) *QueryInit {
	db.where = append(db.where, whereDb{and: false, sqlRaw: sql, op: "raw"})
	return db
}

func (db *QueryInit) StartGroup() *QueryInit {
	db.where = append(db.where, whereDb{and: true, op: "startGroup"})
	return db
}

func (db *QueryInit) StartGroupOr() *QueryInit {
	db.where = append(db.where, whereDb{and: false, op: "startGroup"})
	return db
}

func (db *QueryInit) EndGroup() *QueryInit {
	db.where = append(db.where, whereDb{op: "endGroup"})
	return db
}

func (db *QueryInit) OrderBy(orderBy string, dir string) *QueryInit {
	db.orderBy = append(db.orderBy, orderBy+" "+dir)
	return db
}

func (db *QueryInit) GroupBy(groupBy string) *QueryInit {
	db.groupBy = groupBy
	return db
}

func (db *QueryInit) Limit(limit int, start int) *QueryInit {
	db.limit = "limit " + strconv.Itoa(start) + "," + strconv.Itoa(limit)
	return db
}

func (db *QueryInit) Having(str string) *QueryInit {
	db.having = str
	return db
}

func (db *Init) joinBuild() {
	for _, v := range db.queryBuilder.join {
		db.query = append(db.query, strings.TrimSpace(v))
	}
}

func whereInValue(dt interface{}) []string {
	var valStr []string
	switch v := dt.(type) {
	case []string:
		for _, v2 := range v {
			valStr = append(valStr, "'"+v2+"'")
		}
	case []int:
		for _, v2 := range v {
			valStr = append(valStr, strconv.Itoa(v2))
		}
	case []uint:
		for _, v2 := range v {
			valStr = append(valStr, strconv.Itoa(int(v2)))
		}
	default:
		panic("error, type data where in not support")
	}
	return valStr
}
func (db *Init) whereBuild() {
	if len(db.queryBuilder.where) >= 1 {
		db.query = append(db.query, "where")
	}
	var opBefore string
	for i, v := range db.queryBuilder.where {
		query := ""
		if i != 0 && v.op != "endGroup" && opBefore != "startGroup" {
			query = "and "
			if !v.and {
				query = "or "
			}
		}
		opBefore = v.op
		switch v.op {
		case "startGroup":
			query += "("
		case "endGroup":
			query += ")"
		case "between":
			query += v.field + " between '" + v.starDate + "' and '" + v.endDate + "'"
		case "raw":
			query += v.sqlRaw
		case "in":
			whereIn := whereInValue(v.value)
			query += v.field + " in(" + strings.Join(whereIn, ",") + ")"
		case "notIn":
			whereIn := whereInValue(v.value)
			query += v.field + " not in(" + strings.Join(whereIn, ",") + ")"
		default:
			fieldAr := strings.Split(strings.TrimSpace(v.field), " ")
			getValue := true
			if len(fieldAr) == 1 {
				query += fieldAr[0] + " = "
			} else if len(fieldAr) >= 2 {
				if strings.ToLower(fieldAr[1]) == "sql" {
					query += addSlash(v.value.(string))
					getValue = false
				} else {
					query += fieldAr[0] + " " + fieldAr[1]
				}
			} else {
				continue
			}
			if getValue {
				value := valueInterface(v.value)
				if value[1] == "string" {
					query += " '" + addSlash(value[0]) + "'"
				} else {
					query += " " + value[0]
				}
			}
		}
		//append to query
		db.query = append(db.query, query)
	}
}

func (db *Init) buildQuery() error {

	//add select
	db.queryBuilder.sel = strings.TrimSpace(db.queryBuilder.sel)
	db.queryBuilder.sel = strings.TrimRight(db.queryBuilder.sel, ",")
	if db.queryBuilder.sel == "" {
		//select not init
		db.queryBuilder.sel = "*"
	}

	db.query = append(db.query, "select "+db.queryBuilder.sel)

	//add from
	if db.queryBuilder.from == "" && db.queryBuilder.call == false {
		//table not init
		return errors.New("table not found")
	}
	db.query = append(db.query, "from "+db.queryBuilder.from)

	//add join
	db.joinBuild()

	//add where
	db.whereBuild()

	//add groupBy
	if db.queryBuilder.groupBy != "" {
		db.query = append(db.query, "group by "+db.queryBuilder.groupBy)
	}

	//add having
	if db.queryBuilder.having != "" {
		db.query = append(db.query, "having "+db.queryBuilder.having)
	}

	//add order by
	if len(db.queryBuilder.orderBy) > 0 {
		db.query = append(db.query, "order by ", strings.Join(db.queryBuilder.orderBy, ","))
	}

	//add limit
	if db.queryBuilder.limit != "" {
		db.query = append(db.query, db.queryBuilder.limit)
	}

	return nil
}

func (db *QueryInit) Result() (queryRaw string, err error) {
	var dbs Init
	dbs.queryBuilder = *db
	dbs.query = []string{}
	err = dbs.buildQuery()
	if err != nil {
		return
	}
	queryRaw = strings.Join(dbs.query, "\n")
	return
}

func (db *QueryInit) Row() (queryRaw string, err error) {
	db.Limit(1, 0)
	queryRaw, err = db.Result()
	return
}

func (db *Init) Result() ([]map[string]interface{}, error) {
	var err error
	defer db.Clear()
	if db.queryBuilder.call == false {
		db.query = []string{}
		err = db.buildQuery()
		if err != nil {
			return nil, err
		}
	}
	sqlQuery := strings.Join(db.query, "\n")
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return nil, err
		}
	}

	var rows *Rows
	var stmt *Stmt

	if db.queryBuilder.call == false {
		if db.transaction != nil {
			stmt, err = db.transaction.Prepare(sqlQuery)
		} else {
			stmt, err = db.dbs.Prepare(sqlQuery)
		}
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()

	} else {
		if db.transaction != nil {
			rows, err = db.transaction.Query(sqlQuery)
		} else {
			rows, err = db.dbs.Query(sqlQuery)
		}
		if err != nil {
			return nil, err
		}
	}

	//set call to false
	db.queryBuilder.call = false

	var columns, _ = rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result := []map[string]interface{}{}
	if count == 0 {
		return nil, nil
	}
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		data := map[string]interface{}{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
				v = html.UnescapeString(v.(string))
			} else {
				v = val
			}
			if v == nil {
				data[col] = ""
			} else {
				data[col] = v
			}
		}
		result = append(result, data)
	}
	if stmt != nil {
		err = stmt.Close()
		if err != nil {
			return nil, err
		}
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Init) Row() (map[string]interface{}, error) {
	if db.queryBuilder.call == false {
		db.queryBuilder.Limit(1, 0)
	}
	result, err := db.Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

func (db *Init) ResultStruct(result interface{}) error {
	r, err := db.Result()
	if err != nil {
		return err
	}
	json := jsoniter.Config{EscapeHTML: true, TagKey: "sqlb", OnlyTaggedField: true}.Froze()
	rb, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rb, &result)
	if err != nil {
		return err
	}
	return nil
}

func (db *Init) RowStruct(result interface{}) error {
	r, err := db.Row()
	if err != nil {
		return err
	}
	json := jsoniter.Config{EscapeHTML: true, TagKey: "sqlb", OnlyTaggedField: true}.Froze()
	rb, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rb, &result)
	if err != nil {
		return err
	}
	return nil
}

func (db *Init) Call(procedure string, value []interface{}) *Init {
	values := ""
	if len(value) != 0 {
		var valAr []string
		for _, v := range value {
			var str string
			if db.queryBuilder.removeSpecialChar {
				str = removeSpecialChar(v)
			} else {
				str = valueInterface(v)[0]
			}
			valAr = append(valAr, str)
		}
		values = "('" + strings.Join(valAr, "','") + "')"
	}
	db.query = []string{"call " + procedure + " " + values}
	db.queryBuilder.call = true
	return db
}
