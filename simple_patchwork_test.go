package sqlpatchwork

import (
	"testing"
)

func Test_E2E_NewSimplePatchwork(t *testing.T) {
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

func Test_E2E_NewSimplePWSkipPrs(t *testing.T) {
	qps := map[string]string{
		"prefix":    "INSERT INTO hoge_table (col1, col2) VALUES",
		"loopVal":   "(:col1_@@,:col2_@@)",
		"loopDelim": ",",
	}

	spw := NewSimplePWSkipPrs("skipPrsTest", qps)

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

	expected = "INSERT /* skipPrsTest [prefix loopVal loopDelim loopVal loopDelim loopVal] */ INTO hoge_table (col1, col2) VALUES (:col1_0,:col2_0) , (:col1_1,:col2_1) , (:col1_2,:col2_2)"
	if spw.BuildQueryWithTraceDesc() != expected {
		t.Errorf("E2E of SimplePatchwork is failure.\nEXPECTED: %v\nACTUAL  : %v", expected, spw.BuildQueryWithTraceDesc())
	}
}
