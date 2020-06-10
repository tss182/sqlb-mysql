sqlb-mysql is builder mysql.
You can use mysql so easy in go

`Install library
go get github.com/tss182/sqlb-mysql`

### connection
 	db := sqlb.Init{}
    db.Setup("tcp(127.0.0.1:3306)","userSql","PasswordSql","dbName")
 	defer db.Close()`

db.close is close database connection

sampel select

    db.select("id,name,age").From("member").
    Where("name","Dion"). 
    Where("age >=",10)
    
    result,err := db.result()
### Other Command
 
    More command Join,OrderBy, WhereIn, WhereNotIn,WhereBetween, 
    Having, Limit, etc

### Result

    db.result output []map[string]interface{}
    db.row output map[string]inteface{} with limt 1
    
    
## **Insert, Delete and update coming soon**
