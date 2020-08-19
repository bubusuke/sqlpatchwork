package sqlpatchwork

import (
	"fmt"
	"testing"
)

func TestOnOffPatchwork(t *testing.T) {
	spw, err := NewOnOffPatchwork("test/parse_test.sql")
	if err != nil {
		fmt.Println(err.Error())
	}
	spw.AddQueryPiecesToBuild("cond1")
	fmt.Println(spw.BuildQuery())
	spw.AddQueryPiecesToBuild("cond2")
	fmt.Println(spw.BuildQuery())
	fmt.Println(spw.BuildQueryWithTraceDesc())
	spw, _ = NewOnOffPatchwork("test/parse_test.sql")
	fmt.Println(spw.BuildQuery())
}

func TestSimplePatchwork(t *testing.T) {
	spw, err := NewSimplePatchwork("test/simple_parse_test.sql")
	if err != nil {
		fmt.Println(err.Error())
	}

	data := []map[string]int{
		{"foo": 1, "bar": 2},
		{"foo": 11, "bar": 12},
		{"foo": 21, "bar": 21},
	}
	bindData := make(map[string]interface{})

	spw.AddQueryPiecesToBuild("prefix")
	for i, v := range data {
		if i != 0 {
			spw.AddQueryPiecesToBuild("loop_delim")
		}
		spw.AddQueryPiecesToBuild("loop")
		bindData[LoopNoAttach("foo_@@", i)] = v["foo"]
		bindData[LoopNoAttach("bar_@@", i)] = v["bar"]
	}
	spw.AddQueryPiecesToBuild("surfix")
	fmt.Println(spw.BuildQuery())
	fmt.Println(bindData)
	fmt.Println(spw.BuildQueryWithTraceDesc())
}
