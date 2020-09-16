package main

import (
	"fmt"
	"github.com/tss182/sqlb-mysql"
)

func main() {
	var query sqlb.QueryInit

	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"id":     1,
		"nama":   "triyana",
		"alamat": "bandung",
	})

	data = append(data, map[string]interface{}{
		"id":     2,
		"nama":   "fiqri",
		"alamat": "bandung",
	})

	data = append(data, map[string]interface{}{
		"id":     3,
		"nama":   "dani",
		"alamat": "bandung",
	})

	queryRaw, value, err := query.From("member").UpdateBatch(data, "id")
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	fmt.Println("sqlRaw", queryRaw)
	fmt.Println("value", value)
}
