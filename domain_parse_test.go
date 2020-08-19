package sqlpatchwork

import (
	"fmt"
	"testing"
)

func Test_StartComment(t *testing.T) {
	lines := []string{"/* comment */",
		"SELECT /* comment ",
		"SELECT /**/",
	}
	for _, line := range lines {
		dps := &domainParser{isInCommentBlock: false}
		for i, c := range line {
			if dps.isCommentStart(i, c, &line) {
				fmt.Printf(" '/*' was found at %v\n", i)
			}
		}
	}
}
func Test_EndComment(t *testing.T) {
	lines := []string{"/* comment */",
		"comment */ ",
		"SELECT /**/",
	}
	for _, line := range lines {
		dps := &domainParser{isInCommentBlock: true}
		for i, c := range line {
			if dps.isCommentEnd(i, c, &line) {
				fmt.Printf(" '*/' was found at %v\n", i)
			}
		}
	}
}

func Test_CommentOut(t *testing.T) {
	lines := []string{"-- test",
		"--test",
		"SELECT --test",
	}
	for _, line := range lines {
		dps := &domainParser{isInCommentBlock: false}
		for i, c := range line {
			if dps.isCommentOut(i, c, &line) {
				fmt.Printf(" '--' was found at %v\n", i)
			}
		}
	}
}

func Test_StartKey(t *testing.T) {
	lines := []string{"@start testID",
		"@start */",
		"@starthogehoge",
	}
	for _, line := range lines {
		dps := &domainParser{isInCommentBlock: true, isInPatchBlock: false}
		for i, c := range line {
			if dps.existStartKey(i, c, &line) {
				fmt.Printf(" '@start' was found at %v\n", i)
				patchworkID, err := dps.getID(i, c, &line)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(patchworkID)
				}
			}
		}
	}
}

func Test_EndKey(t *testing.T) {
	lines := []string{"@enddummy",
		"@end */",
		"@end*/",
	}
	for _, line := range lines {
		dps := &domainParser{isInCommentBlock: true, isInPatchBlock: true}
		for i, c := range line {
			if dps.existEndKey(i, c, &line) {
				fmt.Printf(" '@end' was found at %v\n", i)
			}
		}
	}
}
