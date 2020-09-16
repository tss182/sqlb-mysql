package main

import (
	"fmt"
	"github.com/tss182/sqlb-mysql"
)

func main() {
	var query sqlb.QueryInit
	query.Select("p.id,nama").From("produk_main p").
		Join("produk_variant pv", "p.id=pv.id_produk", "left").
		Where("p.id", 2).
		Limit(1, 0).
		OrderBy("nama", "desc").
		OrderBy("alama", "desc")
	queryRaw, err := query.Result()
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	fmt.Println("sqlRaw", queryRaw)
}
