# patchworksql

```
/*@start id=skelton*/
SELECT
  #{}
FROM
  test
WHERE 1=1
/*@end*/
/*@start id=itemNonNull*/
  item_code = :id
/*@end*/
ORDER BY
  , 
```


- 動的SQLの実装
　- SQLファイルの読み込み
  - SQLファイルのパース
　- バインドの名前寄せ
　- SQLの分岐
  - SQLインジェクション
  - struct とのマッピング
