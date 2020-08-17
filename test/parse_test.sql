SELECT
  *
FROM
  hoge_table
WHERE 1=1
  /*@start cond1 */
AND  foo = 1
  /*@end*/
  /*@start cond2 */
AND  bar = 1
  /*@end*/
  /*@start cond1/cond2 */
AND  foobar = 1
  /*@end*/