package sqlb

import (
	. "database/sql"
	"errors"
	"html"
	"strconv"
	"strings"
)

func (db *Init) Select(sel string) *Init {
	db.sel = sel
	return db
}

func (db *Init) Join(table, on, join string) *Init {
	join = strings.ToLower(join)
	joinStr := join + " join " + table + " on " + on
	db.join = append(db.join, joinStr)
	return db
}

func (db *Init) Where(field string, value interface{}) *Init {
	db.where = append(db.where, whereDb{and: true, field: field, value: value})
	return db
}

func (db *Init) WhereOr(field string, value interface{}) *Init {
	db.where = append(db.where, whereDb{and: false, field: field, value: value})
	return db
}

func (db *Init) WhereIn(field string, value []interface{}) *Init {
	db.where = append(db.where, whereDb{and: true, field: field, valueIn: value, op: "in"})
	return db
}

func (db *Init) WhereInOr(field string, value []interface{}) *Init {
	db.where = append(db.where, whereDb{and: false, field: field, valueIn: value, op: "in"})
	return db
}

func (db *Init) WhereNotIn(field string, value []interface{}) *Init {
	db.where = append(db.where, whereDb{and: true, field: field, valueIn: value, op: "notIn"})
	return db
}

func (db *Init) WhereNotInOr(field string, value []interface{}) *Init {
	db.where = append(db.where, whereDb{and: false, field: field, valueIn: value, op: "notIn"})
	return db
}

func (db *Init) WhereBetween(field, startDate, endDate string) *Init {
	db.where = append(db.where, whereDb{and: true, field: field, starDate: startDate, endDate: endDate, op: "between"})
	return db
}

func (db *Init) WhereBetweenOr(field, startDate, endDate string) *Init {
	db.where = append(db.where, whereDb{and: false, field: field, starDate: startDate, endDate: endDate, op: "between"})
	return db
}

func (db *Init) StartGroup() *Init {
	db.where = append(db.where, whereDb{and: true, op: "startGroup"})
	return db
}

func (db *Init) StartGroupOr() *Init {
	db.where = append(db.where, whereDb{and: false, op: "startGroup"})
	return db
}

func (db *Init) EndGroup() *Init {
	db.where = append(db.where, whereDb{op: "endGroup"})
	return db
}

func (db *Init) OrderBy(orderBy string, dir string) *Init {
	db.orderBy = append(db.orderBy, orderBy+" "+dir)
	return db
}

func (db *Init) GroupBy(groupBy string) *Init {
	db.groupBy = groupBy
	return db
}

func (db *Init) Limit(limit int, start int) *Init {
	db.limit = "limit " + strconv.Itoa(start) + "," + strconv.Itoa(limit)
	return db
}

func (db *Init) Having(str string) *Init {
	db.having = str
	return db
}

func (db *Init) joinBuild() {
	for _, v := range db.join {
		db.query = append(db.query, strings.TrimSpace(v))
	}
}

func whereInValue(dt []interface{}) []string {
	var valStr []string
	for _, v := range dt {
		data := valueInterface(v)
		switch data[1] {
		case "string":
			valStr = append(valStr, "'"+data[0]+"'")
		default:
			valStr = append(valStr, data[0])
		}
	}
	return valStr
}
func (db *Init) whereBuild() {
	if len(db.where) >= 1 {
		db.query = append(db.query, "where")
	}
	for i, v := range db.where {
		query := ""
		if i != 0 && v.op != "endGroup" {
			query = "and "
			if !v.and {
				query = "or "
			}
		}
		switch v.op {
		case "startGroup":
			query += "("
		case "endGroup":
			query += ")"
		case "between":
			query += "between '" + v.starDate + "' and '" + v.endDate
		case "in":
			whereIn := whereInValue(v.valueIn)
			query += v.field + " in(" + strings.Join(whereIn, ",") + ")"
		case "notIn":
			whereIn := whereInValue(v.valueIn)
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
	db.sel = strings.TrimSpace(db.sel)
	db.sel = strings.TrimRight(db.sel, ",")
	if db.sel == "" {
		//select not init
		db.sel = "*"
	}

	db.query = append(db.query, "select "+db.sel)

	//add from
	if db.from == "" && db.call == false {
		//table not init
		return errors.New("table not found")
	}
	db.query = append(db.query, "from "+db.from)

	//add join
	db.joinBuild()

	//add where
	db.whereBuild()

	//add groupBy
	if db.groupBy != "" {
		db.query = append(db.query, "group by "+db.groupBy)
	}

	//add having
	if db.having != "" {
		db.query = append(db.query, "having "+db.having)
	}

	//add order by
	for _, v := range db.orderBy {
		db.query = append(db.query, "order by "+v)
	}

	//add limit
	if db.limit != "" {
		db.query = append(db.query, db.limit)
	}

	return nil
}

func (db *Init) Result() ([]map[string]interface{}, error) {
	err := db.buildQuery()
	if err != nil {
		return nil, err
	}
	sqlQuery := strings.Join(db.query, "\n")
	if db.dbs == nil {
		db.dbs, err = db.mysqlConnect()
		if err != nil {
			return nil, err
		}
	}

	var stmt *Stmt
	if db.transaction != nil {
		stmt, err = db.transaction.Prepare(sqlQuery)
	} else {
		stmt, err = db.dbs.Prepare(sqlQuery)
	}
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
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
	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Init) Row() (map[string]interface{}, error) {
	if db.call == false {
		db.Limit(1, 0)
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

func (db *Init) Call(procedure string, value []interface{}) *Init {
	values := ""
	if len(value) != 0 {
		var valAr []string
		for _, v := range value {
			var str string
			if db.removeSpecialChar {
				str = removeSpecialChar(v)
			} else {
				str = valueInterface(v)[0]
			}
			valAr = append(valAr, str)
		}
		values = "('" + strings.Join(valAr, "','") + "')"
	}
	db.query = []string{"call " + procedure + " " + values}
	db.call = true
	return db
}
