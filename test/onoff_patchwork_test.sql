SELECT
 s.item_code
, s.sales_date
, COUNT(*) AS count
FROM
 sales_tran s
/*@start itemTypeNotNil/colorCodeNotNil */
INNER JOIN
 item_master i
ON
 i.item_code = s.item_code
/*@end*/
WHERE 1=1
/*@start itemTypeNotNil*/
AND i.item_type = :item_type
/*@end*/
/*@start colorCodeNotNil*/
AND i.color_code = :color_code
/*@end*/
GROUP BY
 s.item_code
, s.sales_date
ORDER BY
 s.item_code
, s.sales_date
