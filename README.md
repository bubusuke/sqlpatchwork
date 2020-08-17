# [TIP]SQL Patchwork
====

# TODO
  - parseテストの整備
  - Readme.md の整備
  - error処理を上に
  - byte string を行ったりきたり
  - goコメント整理 1h


# Overview
SQL file based dynamic sql-query builder.
SQL Patchwork has only build query function and doesn't have db access function.

## Description
SQL Patchwork provides a dynamic query in 3 steps.

1. Divide the SQL file into query-pieces.
2. Pick query-pieces.
3. Build query through concatenate picked query-pieces.

You can change the query dynamically by implementing step.2.
All dynamic logic required implemented on Golang, so you don't need to learn the special syntax described in SQL file (as typified by MyBatis, dynamic sql library in Java ).

SQL Patchwork has two building modes.

One is simple patchwork mode. Multiple Row query, such as multiple row insert, is possible because the same query piece can be used repeatedly.

One is on/off patchwork mode. Concatenate valid query pieces in the order in the original SQL file. We can only choose apply or not about each query-piece.By limiting the functions, on/off mode keeps the readability of the original SQL file itself.

## Demo
### simple patchwork
SQL file
```
/*@start prefix */
INSERT INTO hoge ( foo, bar) VALUES (
/*@end*/
/*@start loop */
( :foo_@@, :bar_@@ )
/*@end*/
/*@start loop_delim */
,
/*@end*/
/*@start surfix */
)
/*@end*/
```
-> output(loop 3)
```
INSERT INTO hoge ( foo, bar) VALUES ( ( :foo_0, :bar_0 ) , ( :foo_1, :bar_1 ) , ( :foo_2, :bar_2 ) )
```

### simple patchwork
SQL file
```
SELECT
  *
FROM
  hoge_table
WHERE 1=1
  /*@start fooIsNotNull */
AND  foo = :foo
  /*@end*/
  /*@start barIsNotNull */
AND  bar = :bar
  /*@end*/
  /*@start fooIsNotNull/barIsNotNull */
AND  foobar = true
  /*@end*/
```
-> output (when fooIsNotNull is choosed)
```
SELECT * FROM hoge_table WHERE 1=1 AND foo = :foo AND foobar = true
```
-> output (when barIsNotNull is choosed)
```
SELECT * FROM hoge_table WHERE 1=1 bar = :bar AND foobar = true
```
-> output (when not choosed)
```
SELECT * FROM hoge_table WHERE 1=1
```


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
MIT

## Author


# sqlpatchwork



```
SELECT
  #{}
FROM
  test
WHERE 1=1
/*@start itemNonNull*/
  item_code = :id
/*@end*/
ORDER BY
  , 
```



