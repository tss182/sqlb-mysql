package main

import (
	"fmt"
	"github.com/tss182/sqlb-mysql"
)

func main() {
	var query sqlb.QueryInit

	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"nama":   "triyana",
		"alamat": "bandung",
	})

	data = append(data, map[string]interface{}{
		"nama":   "fiqri",
		"alamat": "bandung",
	})

	data = append(data, map[string]interface{}{
		"nama":   "dani",
		"alamat": "bandung",
	})

	queryRaw, value, err := query.From("member").InsertBatch(data)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	fmt.Println("sqlRaw", queryRaw)
	fmt.Println("value", value)
}
