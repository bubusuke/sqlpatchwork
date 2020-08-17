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