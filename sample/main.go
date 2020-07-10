package main

import "github.com/tss182/sqlb-mysql"

func main() {
	db := sqlb.Init{}
	db.Setup("tcp(localhost:3306)", "root", "", "bandros_v4")
	db.From("p").
		Join("ta t", "t.id=p.id", "left").
		Where("id", 2)
	db.Delete()
}
