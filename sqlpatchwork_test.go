package sqlpatchwork

import (
	"testing"
)

func Test_E2E_SimplePatchwork(t *testing.T) {
	_, err := NewSimplePatchwork("./test/nothing")
	if err == nil {
		t.Errorf("Error should be occurrd.")
	}

	spw, err := NewSimplePatchwork("./test/simple_patchwork_test.sql")
	if err != nil {
		t.Errorf("Error should not be occurrd.")
	}
	spw.AddQueryPiecesToBuild("prefix")

	for i := 0; i < 3; i++ {
		if i != 0 {
			spw.AddQueryPiecesToBuild("loopDelim")
		}
		spw.AddQueryPiecesToBuild("loopVal")
	}

	var expected string
	expected = "INSERT INTO hoge_table (col1, col2) VALUES (:col1_0,:col2_0) , (:col1_1,:col2_1) , (:col1_2,:col2_2)"
	if spw.BuildQuery() != expected {
		t.Errorf("E2E of SimplePatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQuery())
	}

	expected = "INSERT /* ./test/simple_patchwork_test.sql [prefix loopVal loopDelim loopVal loopDelim loopVal] */ INTO hoge_table (col1, col2) VALUES (:col1_0,:col2_0) , (:col1_1,:col2_1) , (:col1_2,:col2_2)"
	if spw.BuildQueryWithTraceDesc() != expected {
		t.Errorf("E2E of SimplePatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQueryWithTraceDesc())
	}
}

func Test_E2E_OnOffPatchwork(t *testing.T) {
	_, err := NewOnOffPatchwork("./test/nothing")
	if err == nil {
		t.Errorf("Error should be occurrd.")
	}

	spw, err := NewOnOffPatchwork("./test/onoff_patchwork_test.sql")
	if err != nil {
		t.Errorf("Error should not be occurrd.")
	}
	spw.AddQueryPiecesToBuild("itemTypeNotNil")

	var expected string
	expected = "SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.item_type = :item_type GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"
	if spw.BuildQuery() != expected {
		t.Errorf("E2E of OnOffPatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQuery())
	}

	expected = "SELECT /* ./test/onoff_patchwork_test.sql [__default itemTypeNotNil] */ s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.item_type = :item_type GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"
	if spw.BuildQueryWithTraceDesc() != expected {
		t.Errorf("E2E of OnOffPatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQueryWithTraceDesc())
	}

	spw, err = NewOnOffPatchwork("./test/onoff_patchwork_test.sql")
	if err != nil {
		t.Errorf("Error should not be occurrd.")
	}
	expected = "SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s WHERE 1=1 GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"
	if spw.BuildQuery() != expected {
		t.Errorf("E2E of OnOffPatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQuery())
	}

}
