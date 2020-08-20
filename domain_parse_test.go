package sqlpatchwork

import (
	"testing"
)

func TestDomainParser_isCommentStart(t *testing.T) {
	dps := &domainParser{}
	var line string

	dps.isInCommentBlock = false
	line = "/* comment */"
	for i, c := range line {
		if i == 0 {
			if !dps.isCommentStart(i, c, &line) {
				t.Errorf("isCommentStart should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.isCommentStart(i, c, &line) {
				t.Errorf("isCommentStart should be FALSE at index %v.\n", i)
			}
		}
	}

	dps.isInCommentBlock = true
	for i, c := range line {
		if dps.isCommentStart(i, c, &line) {
			t.Errorf("isCommentStart should be FALSE in case isInCommentBlock status is true.\n")
		}
	}

	dps.isInCommentBlock = false
	line = "SELECT /* comment "
	for i, c := range line {
		if i == 7 {
			if !dps.isCommentStart(i, c, &line) {
				t.Errorf("isCommentStart should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.isCommentStart(i, c, &line) {
				t.Errorf("isCommentStart should be FALSE at index %v.\n", i)
			}
		}
	}

	dps.isInCommentBlock = true
	for i, c := range line {
		if dps.isCommentStart(i, c, &line) {
			t.Errorf("isCommentStart should be FALSE in case isInCommentBlock status is true.\n")
		}
	}
}

func TestDomainParser_isCommentEnd(t *testing.T) {
	dps := &domainParser{}
	var line string

	dps.isInCommentBlock = true
	line = "comment */  "
	for i, c := range line {
		if i == 9 {
			if !dps.isCommentEnd(i, c, &line) {
				t.Errorf("isCommentEnd should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.isCommentEnd(i, c, &line) {
				t.Errorf("isCommentEnd should be FALSE at index %v.\n", i)
			}
		}
	}

	dps.isInCommentBlock = false
	for i, c := range line {
		if dps.isCommentEnd(i, c, &line) {
			t.Errorf("isCommentEnd should be FALSE in case isInCommentBlock status is false.\n")
		}
	}
}

func TestDomainParser_isCommentOut(t *testing.T) {
	dps := &domainParser{}
	var line string

	dps.isInCommentBlock = false
	line = "SELECT --test"
	for i, c := range line {
		if i == 7 {
			if !dps.isCommentOut(i, c, &line) {
				t.Errorf("isCommentOut should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.isCommentOut(i, c, &line) {
				t.Errorf("isCommentOut should be FALSE at index %v.\n", i)
			}
		}
	}

	dps.isInCommentBlock = true
	for i, c := range line {
		if dps.isCommentOut(i, c, &line) {
			t.Errorf("isCommentOut should be FALSE in case isInCommentBlock status is true.\n")
		}
	}
}

func TestDomainParser_existStartKey(t *testing.T) {
	dps := &domainParser{}
	var line string

	dps.isInCommentBlock = true
	dps.isInPatchBlock = false
	line = " @start testID"
	for i, c := range line {
		if i == 1 {
			if !dps.existStartKey(i, c, &line) {
				t.Errorf("existStartKey should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.existStartKey(i, c, &line) {
				t.Errorf("existStartKey should be FALSE at index %v.\n", i)
			}
		}
	}
	dps.isInCommentBlock = false
	dps.isInPatchBlock = false
	for i, c := range line {
		if dps.existStartKey(i, c, &line) {
			t.Errorf("existStartKey should be FALSE when isInCommentBlock is false.\n")
		}
	}

	dps.isInCommentBlock = true
	dps.isInPatchBlock = true
	for i, c := range line {
		if dps.existStartKey(i, c, &line) {
			t.Errorf("existStartKey should be FALSE when isInPatchBlock is true.\n")
		}
	}
}

func TestDomainParser_existEndKey(t *testing.T) {
	dps := &domainParser{}
	var line string

	dps.isInCommentBlock = true
	dps.isInPatchBlock = true
	line = " @end */"
	for i, c := range line {
		if i == 1 {
			if !dps.existEndKey(i, c, &line) {
				t.Errorf("existEndKey should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.existEndKey(i, c, &line) {
				t.Errorf("existEndKey should be FALSE at index %v.\n", i)
			}
		}
	}
	dps.isInCommentBlock = false
	dps.isInPatchBlock = true
	line = " @end */"
	for i, c := range line {
		if dps.existEndKey(i, c, &line) {
			t.Errorf("existEndKey should be FALSE when isInCommentBlock is false.\n")
		}
	}
	dps.isInCommentBlock = true
	dps.isInPatchBlock = false
	line = " @end */"
	for i, c := range line {
		if dps.existEndKey(i, c, &line) {
			t.Errorf("existEndKey should be FALSE when isInPatchBlock is false.\n")
		}
	}

	dps.isInCommentBlock = true
	dps.isInPatchBlock = true
	line = " @end*/"
	for i, c := range line {
		if i == 1 {
			if !dps.existEndKey(i, c, &line) {
				t.Errorf("existEndKey should be TRUE at index %v.\n", i)
			}
		} else {
			if dps.existEndKey(i, c, &line) {
				t.Errorf("existEndKey should be FALSE at index %v.\n", i)
			}
		}
	}

	dps.isInCommentBlock = true
	dps.isInPatchBlock = true
	line = " @endd*/"
	for i, c := range line {
		if dps.existEndKey(i, c, &line) {
			t.Errorf("existEndKey should be FALSE.\n")
		}
	}

}

func TestDomainParser_getID(t *testing.T) {
	dps := &domainParser{}
	var line string

	line = "@start foo bar"
	for i, c := range line {
		if i <= 6 {
			expected := "foo"
			ID, err := dps.getID(i, c, &line)
			if err != nil {
				t.Errorf("getID should not be ERROR. ERROR:%v\n", err)
			}
			if ID != expected {
				t.Errorf("getID should be '%v'. Actual '%v'.\n", expected, ID)
			}
		} else if i <= 10 {
			expected := "bar"
			ID, err := dps.getID(i, c, &line)
			if err != nil {
				t.Errorf("getID should not be ERROR. ERROR:%v\n", err)
			}
			if ID != expected {
				t.Errorf("getID should be '%v'. Actual '%v'.\n", expected, ID)
			}
		} else {
			_, err := dps.getID(i, c, &line)
			if err != getIDError {
				t.Errorf("getID should be ERROR. Actual error is %v.\n", err)
			}
		}
	}

	line = "@start  foo		bar"
	for i, c := range line {
		if i <= 7 {
			expected := "foo"
			ID, err := dps.getID(i, c, &line)
			if err != nil {
				t.Errorf("getID should not be ERROR. ERROR:%v\n", err)
			}
			if ID != expected {
				t.Errorf("getID should be '%v'. Actual '%v'.\n", expected, ID)
			}
		} else if i <= 12 {
			expected := "bar"
			ID, err := dps.getID(i, c, &line)
			if err != nil {
				t.Errorf("getID should not be ERROR. ERROR:%v\n", err)
			}
			if ID != expected {
				t.Errorf("getID should be '%v'. Actual '%v'.\n", expected, ID)
			}
		} else {
			_, err := dps.getID(i, c, &line)
			if err != getIDError {
				t.Errorf("getID should be ERROR. Actual error is %v.\n", err)
			}
		}
	}
}

func TestDomainParser_appendEachQP(t *testing.T) {
	dps := &domainParser{
		queryBuf: []byte("hogehoge"),
		tmpIDs:   []string{"foo", "bar"},
		onOffQueryPieces: onOffQPs{
			{
				IDs:   []string{"a"},
				query: []byte("b"),
			},
		},
		defaultValue: func() ([]string, map[string]bool) { return []string{"ini"}, map[string]bool{"test": true} },
	}
	dps.appendQP = dps.onOffAppendQp

	dps.appendEachQP()
	if string(dps.queryBuf) != " " {
		t.Errorf("appendEachQP should initialize queryBuf.\n")
	}
	for _, qp := range dps.onOffQueryPieces {
		if qp.IDs[0] == "foo" && qp.IDs[1] == "bar" {
			if string(qp.query) != "hogehoge" {
				t.Errorf("appendEachQP is something wrong.\n")
			}
		} else if qp.IDs[0] == "a" {
			if string(qp.query) != "b" {
				t.Errorf("appendEachQP is something wrong.\n")
			}
		} else {
			t.Errorf("appendEachQP is something wrong.\n")
		}
	}

	dps = &domainParser{
		queryBuf: []byte("   		   	 	 	 	 	 	"),
		tmpIDs: []string{"foo", "bar"},
		onOffQueryPieces: onOffQPs{
			{
				IDs:   []string{"a"},
				query: []byte("b"),
			},
		},
		defaultValue: func() ([]string, map[string]bool) { return []string{"ini"}, map[string]bool{"test": true} },
	}
	dps.appendQP = dps.onOffAppendQp

	dps.appendEachQP()
	if string(dps.queryBuf) != " " {
		t.Errorf("appendEachQP should initialize queryBuf.\n")
	}
	for _, qp := range dps.onOffQueryPieces {
		if qp.IDs[0] == "a" {
			if string(qp.query) != "b" {
				t.Errorf("appendEachQP is something wrong.\n")
			}
		} else {
			t.Errorf("appendEachQP is something wrong.\n")
		}
	}
}
