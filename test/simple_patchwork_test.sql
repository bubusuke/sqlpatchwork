/*@start prefix */
INSERT INTO hoge_table (col1, col2) VALUES 
/*@end*/
/*@start loopVal */
(:col1_@@,:col2_@@)
/*@end*/
/*@start loopDelim */
,
/*@end*/