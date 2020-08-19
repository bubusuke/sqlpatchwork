# [WIP]SQL Patchwork
====

# TODO
  - parseテストの整備
  - error処理を上に
  - goコメント整理 1h

# Overview
SQL file based dynamic sql-query builder.
SQL Patchwork has only build query function and doesn't have db access function.

## Description
### Summary
SQL Patchwork provides a dynamic query in 3 steps.

1. Divide the SQL file into query-pieces.
2. Pick query-pieces on golang.
3. Build query through concatenate picked query-pieces.

You can change the query dynamically by implementing step.2.
All dynamic logic required implemented on Golang, so you don't need to learn the special syntax described in SQL file (as typified by MyBatis, dynamic sql library in Java ).

##### Image
![Summary.JPG](https://github.com/bubusuke/sqlpatchwork/blob/master/doc_materials/Summary.JPG)

##### Point 1
* Start the Query-Piece block with the `@start` keyword.
* The word following the `@start` is the Query-piece ID.
    * The second and subsequent words will be ignored. eg. `/*@start ThisIsID ThisIsIgnored*/`
    * In SimplePatchwork mode (described later), IDs must be unique.
    * In OnOffPatchwork mode (described later), the same ID can be assigned to multiple Query-pieces.
    * In OnOff Patchwork mode (described later), multiple IDs can be assigned to one Query-piece by separating them with `/`.
* End the Query-Piece block with the `@end` keyword.
* Nested structure of Query-Piece block is not supported.

##### Point 2
* All Commentout and Commentblock in SQL are ignored (not output in the final output).

##### Point 3
* It provides two modes, SimplePatchwork mode and OnOffPatchwork mode, which differ in the method of selecting Query-Pieces. The processing of the figure shown in the example can be realized in either mode (details will be described later).

##### Point 4
* The final output of this function is SQL Query text (type string). There is no function such as DB access.
* The final output is a query that has undergone formatting such as no line breaks, no comments, and trim.
* It is also possible to output a query in which the SQL file used and the ID of the selected QueryPiece are described as comments.

### Concept

#### Low learning cost
We decided to reduce the learning cost as much as possible in order to reduce the special syntax written in SQL file.

#### Coexistence with DB access library
There are already some great libraries around DB access and ORM. We aimed for coexistence with such libraries. In other words, this library is just "SQL Query Builder" and does not have any functions related to DB access and ORM.

#### Easy to test
If tools have special syntax on SQL files, you will not benefit from Golang's great testing tools. For this reason, we chose to implement the selection logic in Golang. This will not bother you with the problem of "can you cover all branching logic with special control syntax?"

### Two modes (SimplePatchwork mode and OnOffPatchwork mode)
It provides two functions, SimplePatchwork mode and OnOffPatchwork mode, which differ in the method of picking Query-Pieces.

#### SimplePatchwork
This is a simple and powerful mode that **concatenates the Query-Pieces in the picked order**.

##### Image
![SimplePatchwork.JPG](https://github.com/bubusuke/sqlpatchwork/blob/master/doc_materials/SimplePatchwork.JPG)

##### Point 1
* Query-piece ID must be set uniquely.

##### Point 2, 3
* It is possible to pick the same Query-Piece multiple times.
* If you add the `@@` keyword to the bind variable, the bind variable name will be rewritten as many times as the picked Query-Pieces (This function not available in OnOff Patchwork). This allows you to build queries such as Multiple-Row-insert.

#### OnOffPatchwork
Since SimplePatchwork is a simple concatenated mode, it is highly versatile. But there is concern that the readability will be reduced because the picking logic will become long and query text in the SQL file will become complex.

OnOffPatchwork is the mode that provides a solution to this concern.
In this mode, **Query-pieces are concatenated in the order they are written in the SQL file**.
And then, when it concatenate them, it judge whether each Query-Piece is applicable or not.

##### Image
![OnOffPatchwork.JPG](https://github.com/bubusuke/sqlpatchwork/blob/master/doc_materials/OnOffPatchwork.JPG)

##### Point 1
* The Query-Pieces that are used at any time can omit the specification of `@start` and `@end` (the omitted parts are given an ID of `__default`).

##### Point 2
* The same ID can be given to multiple Query-Pieces.
* You can assign multiple IDs to a query. To do so, you need to write multiple IDs separating by `/`.

##### Point 3
* Picked IDs concatenate the corresponding Query-pieces in the order written in the SQL file.


### Special syntaxes requred for SQL files in this tool
Only three.

* `/*@start someID*/` starts Query-Piece and the ID is "someID".
* `/*@end*/` : ends Query-Piece.
* `@@` : It is replaced by the number of times Query-piece is used. If you use it repeatedly, it will be incremented as `0`,`1`. For details, please check Point 2 in SimplePatchwork.

## Usage
### SimplePatchwork
```go
package main

import (
	"fmt"
	"log"

	"github.com/bubusuke/sqlpatchwork"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type hoge struct {
	Col1 int    `db:"col1"`
	Col2 string `db:"col2"`
}

func main() {
	hoges := []hoge{
		{Col1: 100, Col2: "foo"},
		{Col1: 200, Col2: "bar"},
		{Col1: 300, Col2: "hoge"},
	}
	spw, err := sqlpatchwork.NewSimplePatchwork("./sqls/simplePatchwork.sql")
	if err != nil {
		fmt.Println(err)
	}
	bindMap := make(map[string]interface{})

	if err := spw.AddQueryPiecesToBuild("prefix"); err != nil {
		fmt.Println(err)
	}
	for i, h := range hoges {
		if i != 0 {
			if err := spw.AddQueryPiecesToBuild("loopDelim"); err != nil {
				fmt.Println(err)
			}
		}
		if err := spw.AddQueryPiecesToBuild("loopVal"); err != nil {
			fmt.Println(err)
		}
		bindMap[sqlpatchwork.LoopNoAttach("col1_@@", i)] = h.Col1
		bindMap[sqlpatchwork.LoopNoAttach("col2_@@", i)] = h.Col2
	}

	fmt.Println("==========================")
	fmt.Println("TargetIDs are")
	fmt.Println(spw.TargetIDs())
	fmt.Println("==========================")
	fmt.Println("BuildingQuery is")
	fmt.Println(spw.BuildQuery())
	fmt.Println("==========================")
	fmt.Println("BuildingQueryWithTrace is")
	fmt.Println(spw.BuildQueryWithTraceDesc())
	fmt.Println("==========================")
	fmt.Println("BindMap is")
	fmt.Println(bindMap)
	fmt.Println("==========================")

	insertDemoByUsingSqlx(spw.BuildQuery(), bindMap)
}

//insertDemoByUsingSqlx is using jomoiron/sqlx to execute SQL.
func insertDemoByUsingSqlx(query string, bindMap map[string]interface{}) {
	db, err := sqlx.Connect("postgres", "user=postgres password=pass dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec("DROP TABLE IF EXISTS hoge_table")
	db.MustExec("CREATE TABLE hoge_table (col1 int, col2 varchar(25))")

	tx := db.MustBegin()
	_, err = tx.NamedExec(query, bindMap)
	if err != nil {
		log.Fatalln(err)
	}
	tx.Commit()

	hoges := []hoge{}
	err = db.Select(&hoges, "SELECT col1, col2 FROM hoge_table ORDER BY col1")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("==========================")
	fmt.Println("The result of query execution is")
	fmt.Println(hoges)
	fmt.Println("==========================")

}

// The file content of ./sqls/simplePatchwork.sql
// --------------------------------------------
// /*@start prefix */
// INSERT INTO hoge_table (col1, col2) VALUES 
// /*@end*/
// /*@start loopVal */
// (:col1_@@,:col2_@@)
// /*@end*/
// /*@start loopDelim */
// ,
// /*@end*/

// STDOUT
// --------------------------------------------
// ==========================
// TargetIDs are
// [prefix loopVal loopDelim loopVal loopDelim loopVal]
// ==========================
// BuildingQuery is
// INSERT INTO hoge_table (col1, col2) VALUES (:col1_0,:col2_0) , (:col1_1,:col2_1) , (:col1_2,:col2_2)
// ==========================
// BuildingQueryWithTrace is
// INSERT /* ./sqls/simplePatchwork.sql [prefix loopVal loopDelim loopVal loopDelim loopVal] */ INTO hoge_table (col1, col2) VALUES (:col1_0,:col2_0) , (:col1_1,:col2_1) , (:col1_2,:col2_2)
// ==========================
// BindMap is
// map[col1_0:100 col1_1:200 col1_2:300 col2_0:foo col2_1:bar col2_2:hoge]
// ==========================
// ==========================
// The result of query execution is
// [{100 foo} {200 bar} {300 hoge}]
// ==========================
```
### OnOffPatchwork
```go
package main

import (
	"fmt"
	"log"

	"github.com/bubusuke/sqlpatchwork"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type req struct {
	ItemType  string `db:"item_type"`
	ColorCode string `db:"color_code"`
}

type Res struct {
	ItemCode  int    `db:"item_code"`
	SalesDate string `db:"sales_date"`
	Count     int    `db:"count"`
}

func buildQuery(r *req, isTraceDesc bool) string {
	spw, err := sqlpatchwork.NewOnOffPatchwork("./sqls/onoffPatchwork.sql")
	if err != nil {
		fmt.Println(err)
	}
	if r.ItemType != "" {
		spw.AddQueryPiecesToBuild("itemTypeNotNil")
	}
	if r.ColorCode != "" {
		spw.AddQueryPiecesToBuild("colorCodeNotNil")
	}
	fmt.Println("==========================")
	fmt.Println("request is")
	fmt.Println(r)
	fmt.Println("==========================")
	fmt.Println("TargetIDs are")
	fmt.Println(spw.TargetIDs())
	fmt.Println("==========================")
	fmt.Println("BuildingQuery is")
	fmt.Println(spw.BuildQuery())
	fmt.Println("==========================")
	fmt.Println("BuildingQueryWithTrace is")
	fmt.Println(spw.BuildQueryWithTraceDesc())
	fmt.Println("==========================")

	if isTraceDesc {
		return spw.BuildQueryWithTraceDesc()
	}
	return spw.BuildQuery()
}

func main() {

	resetDB()

	req1 := &req{ItemType: "T-shirts"}
	query1 := buildQuery(req1, false)
	selectDemoByUsingSqlx(query1, req1)

	req2 := &req{ColorCode: "W"}
	query2 := buildQuery(req2, false)
	selectDemoByUsingSqlx(query2, req2)

	req3 := &req{}
	query3 := buildQuery(req3, false)
	selectDemoByUsingSqlx(query3, req3)
}

//resetDB reset DB Data.
func resetDB() {
	db, err := sqlx.Connect("postgres", "user=postgres password=pass dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec("DROP TABLE IF EXISTS sales_tran")
	db.MustExec("DROP TABLE IF EXISTS item_master")
	db.MustExec("CREATE TABLE item_master (item_code varchar(4), item_type varchar(25), color_code varchar(2))")
	db.MustExec("INSERT INTO item_master VALUES " +
		"('1000', 'T-shirts', 'B')" +
		",('2000', 'T-shirts', 'W')" +
		",('3000', 'pants', 'W')")

	db.MustExec("CREATE TABLE sales_tran (item_code varchar(4), sales_date date)")
	db.MustExec("INSERT INTO sales_tran VALUES " +
		"('1000','2020-08-20')" +
		",('1000','2020-08-20')" +
		",('1000','2020-08-19')" +
		",('2000','2020-08-20')" +
		",('2000','2020-08-19')" +
		",('2000','2020-08-19')" +
		",('3000','2020-08-19')" +
		",('3000','2020-08-18')" +
		",('3000','2020-08-18')")
}

//selectDemoByUsingSqlx is using jomoiron/sqlx to execute SQL.
func selectDemoByUsingSqlx(query string, rq *req) {
	db, err := sqlx.Connect("postgres", "user=postgres password=pass dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.NamedQuery(query, rq)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("==========================")
	fmt.Println("The result of query execution is")
	res := Res{}
	for rows.Next() {
		err := rows.StructScan(&res)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(res)
	}
	fmt.Println("==========================")

}

// The file content of ./sqls/onoffPatchwork.sql
// --------------------------------------------
// SELECT
//  s.item_code
// , s.sales_date
// , COUNT(*) AS count
// FROM
//  sales_tran s
// /*@start itemTypeNotNil/colorCodeNotNil */
// INNER JOIN
//  item_master i
// ON
//  i.item_code = s.item_code
// /*@end*/
// WHERE 1=1
// /*@start itemTypeNotNil*/
// AND i.item_type = :item_type
// /*@end*/
// /*@start colorCodeNotNil*/
// AND i.color_code = :color_code
// /*@end*/
// GROUP BY
//  s.item_code
// , s.sales_date
// ORDER BY
//  s.item_code
// , s.sales_date

// STDOUT
// --------------------------------------------
// ==========================
// request is
// &{T-shirts }
// ==========================
// TargetIDs are
// [__default itemTypeNotNil]
// ==========================
// BuildingQuery is
// SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.item_type = :item_type GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date
// ==========================
// BuildingQueryWithTrace is
// SELECT /* ./sqls/onoffPatchwork.sql [__default itemTypeNotNil] */ s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.item_type = :item_type GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date
// ==========================
// ==========================
// The result of query execution is
// {1000 2020-08-19T00:00:00Z 1}
// {1000 2020-08-20T00:00:00Z 2}
// {2000 2020-08-19T00:00:00Z 2}
// {2000 2020-08-20T00:00:00Z 1}
// ==========================
// ==========================
// request is
// &{ W}
// ==========================
// TargetIDs are
// [__default colorCodeNotNil]
// ==========================
// BuildingQuery is
// SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.color_code = :color_code GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date
// ==========================
// BuildingQueryWithTrace is
// SELECT /* ./sqls/onoffPatchwork.sql [__default colorCodeNotNil] */ s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.color_code = :color_code GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date
// ==========================
// ==========================
// The result of query execution is
// {2000 2020-08-19T00:00:00Z 2}
// {2000 2020-08-20T00:00:00Z 1}
// {3000 2020-08-18T00:00:00Z 2}
// {3000 2020-08-19T00:00:00Z 1}
// ==========================
// ==========================
// request is
// &{ }
// ==========================
// TargetIDs are
// [__default]
// ==========================
// BuildingQuery is
// SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s WHERE 1=1 GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date
// ==========================
// BuildingQueryWithTrace is
// SELECT /* ./sqls/onoffPatchwork.sql [__default] */ s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s WHERE 1=1 GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date
// ==========================
// ==========================
// The result of query execution is
// {1000 2020-08-19T00:00:00Z 1}
// {1000 2020-08-20T00:00:00Z 2}
// {2000 2020-08-19T00:00:00Z 2}
// {2000 2020-08-20T00:00:00Z 1}
// {3000 2020-08-18T00:00:00Z 2}
// {3000 2020-08-19T00:00:00Z 1}
// ==========================


```
## Install
```
go get github.com/bubusuke/sqlpatchwork
```

## Contribution
1. Fork
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## Licence
MIT [https://github.com/bubusuke/sqlpatchwork/blob/master/LICENSE](https://github.com/bubusuke/sqlpatchwork/blob/master/LICENSE)

## Author
[bubusuke](https://github.com/bubusuke)
