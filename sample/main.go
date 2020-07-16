package main

import (
	"fmt"
	"github.com/tss182/sqlb-mysql"
)

func main() {
	db := sqlb.Init{}
	db.Setup("tcp(localhost:3306)", "root", "", "bandros_v4")
	var query sqlb.QueryInit
	query.Select("p.id,nama").From("produk_main p").
		Join("produk_variant pv", "p.id=pv.id_produk", "left").
		Where("p.id", 2)
	db.Query(query)
	data, err := db.Row()
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	fmt.Println("data", data)
}
