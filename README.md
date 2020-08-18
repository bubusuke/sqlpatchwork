# [WIP]SQL Patchwork
====

# TODO
  - usage
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
* `@@` : It is replaced by the number of times Query-piece is used. If you use it repeatedly, it will be incremented as `1`,`2`. For details, please check Point 2 in SimplePatchwork.

## Usage

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
