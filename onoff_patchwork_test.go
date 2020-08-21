package sqlpatchwork

import (
	"testing"
)

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

func Test_E2E_NewOnOffPWSkipPrs(t *testing.T) {
	qps := OnOffQPs(
		OnOffQP("SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s"),
		OnOffQP("INNER JOIN item_master i ON i.item_code = s.item_code", "itemTypeNotNil", "colorCodeNotNil"),
		OnOffQP("WHERE 1=1"),
		OnOffQP("AND i.item_type = :item_type", "itemTypeNotNil"),
		OnOffQP("AND i.color_code = :color_code", "colorCodeNotNil"),
		OnOffQP("GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"))

	spw := NewOnOffPWSkipPrs("skipPrsTest", qps)

	spw.AddQueryPiecesToBuild("itemTypeNotNil")

	var expected string
	expected = "SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.item_type = :item_type GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"
	if spw.BuildQuery() != expected {
		t.Errorf("E2E of OnOffPatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQuery())
	}

	expected = "SELECT /* skipPrsTest [__default itemTypeNotNil] */ s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s INNER JOIN item_master i ON i.item_code = s.item_code WHERE 1=1 AND i.item_type = :item_type GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"
	if spw.BuildQueryWithTraceDesc() != expected {
		t.Errorf("E2E of OnOffPatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQueryWithTraceDesc())
	}

	spw = NewOnOffPWSkipPrs("skipPrsTest", qps)
	expected = "SELECT s.item_code , s.sales_date , COUNT(*) AS count FROM sales_tran s WHERE 1=1 GROUP BY s.item_code , s.sales_date ORDER BY s.item_code , s.sales_date"
	if spw.BuildQuery() != expected {
		t.Errorf("E2E of OnOffPatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQuery())
	}
}
