package sqlpatchwork

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var getIDError = errors.New("queryPieceID is not found. Please describe ID after '@start'. eg. /*@start ThisIsID*/")

// parseResult represent the result of file parsing.
type parseResult struct {
	onOffQueryPieces  []onOffQP
	simpleQueryPieces simpleQPs
	queryPieceIDs     map[string]bool
}

// domainParser represent the skelton parser.
type domainParser struct {
	isInCommentBlock  bool
	isInPatchBlock    bool
	isInCommentOut    bool
	queryBuf          []byte
	tmpIDs            []string
	onOffQueryPieces  []onOffQP
	simpleQueryPieces simpleQPs
	queryPieceIDs     map[string]bool
	defaultValue      func() ([]string, map[string]bool)
	customParseID     func(string) ([]string, error)
	checkEndedCorrect func() error
	appendQP          func()
}

func newDomainParser() *domainParser {
	return &domainParser{
		isInCommentBlock:  false,
		isInPatchBlock:    false,
		isInCommentOut:    false,
		queryBuf:          nil,
		tmpIDs:            nil,
		onOffQueryPieces:  nil,
		simpleQueryPieces: make(simpleQPs),
		queryPieceIDs:     nil,
	}
}

// parse parses buf-reader.
func (dps *domainParser) parse(reader *bufio.Reader) (*parseResult, error) {
	dps.tmpIDs, dps.queryPieceIDs = dps.defaultValue()
	for {
		lineBuf, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			// All read done.
			break
		} else if err != nil {
			return nil, err
		}

		if isPrefix {
			// Same line to previous loop step.
			if dps.isInCommentOut {
				// When comment out was found at previous loop step, this line is still in comment out.
				continue
			}
		} else {
			dps.isInCommentOut = false
		}
		line := string(lineBuf)

		// Parse rune by rune.
		for i, c := range line {
			if dps.isQuery(i, c, &line) {
				dps.queryBuf = append(dps.queryBuf, byte(c))
			}
			if dps.isCommentOut(i, c, &line) {
				dps.isInCommentOut = true
				break
			}
			if dps.isCommentStart(i, c, &line) {
				dps.isInCommentBlock = true
			}
			if dps.isCommentEnd(i, c, &line) {
				dps.isInCommentBlock = false
			}
			if dps.existStartKey(i, c, &line) {
				dps.isInPatchBlock = true
				if err = dps.appendEachQP(); err != nil {
					return nil, err
				}
				tmpID, err := dps.getID(i, c, &line)
				if err != nil {
					return nil, err
				}
				dps.tmpIDs, err = dps.customParseID(tmpID)
				if err != nil {
					return nil, err
				}
				for _, id := range dps.tmpIDs {
					dps.queryPieceIDs[id] = true
				}
			}
			if dps.existEndKey(i, c, &line) {
				dps.isInPatchBlock = false
				if err = dps.appendEachQP(); err != nil {
					return nil, err
				}
			}
		}
		dps.queryBuf = append(dps.queryBuf, " "...)
	}

	if dps.queryBuf != nil {
		if err := dps.checkEndedCorrect(); err != nil {
			return nil, err
		}
		if err := dps.appendEachQP(); err != nil {
			return nil, err
		}
	}

	return &parseResult{
		onOffQueryPieces:  dps.onOffQueryPieces,
		simpleQueryPieces: dps.simpleQueryPieces,
		queryPieceIDs:     dps.queryPieceIDs,
	}, nil
}

// appendQueryPiece appends a query-piece to domainParser.queryPieces and initializes querybuffer and tmpIDs.
// When queryBuffer is blank (only space or tab), not append.
func (dps *domainParser) appendEachQP() error {
	if dps.spaceTabRemove(string(dps.queryBuf)) != "" {
		if len(dps.tmpIDs) == 0 {
			return errors.New("Id not found.")
		}
		dps.appendQP()
	}
	dps.queryBuf = []byte(" ")
	dps.tmpIDs, _ = dps.defaultValue()
	return nil
}

// spaceTabRemove removes all spaces and tabs from str and return it.
func (dps *domainParser) spaceTabRemove(str string) string {
	return strings.Replace(strings.Replace(str, " ", "", -1), "	", "", -1)
}

// isBufAppend gets whether this position is in sql query or in comment.
// true: sqlquery.
func (dps *domainParser) isQuery(i int, c rune, str *string) bool {
	if dps.isInCommentBlock {
		return false
	}
	if dps.isCommentOut(i, c, str) {
		return false
	}
	if dps.isCommentStart(i, c, str) {
		return false
	}
	if dps.isCommentEnd(i, c, str) {
		return false
	}
	return true
}

// isCommentStart gets whether "/*" found here or not.
// true: "/*" found.
func (dps *domainParser) isCommentStart(i int, c rune, str *string) bool {
	if dps.isInCommentBlock {
		return false
	}
	if isLast := len([]rune(*str)) == i+1; isLast {
		return false
	}
	if c == '/' {
		if []rune(*str)[i+1] == '*' {
			return true
		}
	}
	return false
}

// isCommentEnd gets whether "*/" found here or not.
// true: "*/" found.
func (dps *domainParser) isCommentEnd(i int, c rune, str *string) bool {
	if !dps.isInCommentBlock {
		return false
	}
	if isFirst := i == 0; isFirst {
		return false
	}
	if c == '/' {
		if []rune(*str)[i-1] == '*' {
			return true
		}
	}
	return false
}

// isCommentOut gets whether "--" found here or not.
// true: "--" found.
func (dps *domainParser) isCommentOut(i int, c rune, str *string) bool {
	if dps.isInCommentBlock {
		return false
	}
	if isLast := len([]rune(*str)) == i+1; isLast {
		return false
	}
	if c == '-' {
		if []rune(*str)[i+1] == '-' {
			return true
		}
	}
	return false
}

// isCommentOut gets whether "@start " found here or not.
// true: "@start " found.
// The space in "@start " is important in order to avoid misreading "@startX"
func (dps *domainParser) existStartKey(i int, c rune, str *string) bool {
	if !dps.isInCommentBlock || dps.isInPatchBlock {
		return false
	}
	key := "@start "
	keyLen := len(key)
	if i+keyLen > len([]rune(*str)) {
		return false
	}
	if c == '@' {
		if (*str)[i:i+keyLen] == key {
			return true
		}
	}
	return false
}

// isCommentOut gets whether "@end " or "@end*" found here or not.
// true: they found.
// The space in "@end " is important in order to avoid misreading "@startX"
// The "*" in "@end*" is important in order to allow "@end*/".
func (dps *domainParser) existEndKey(i int, c rune, str *string) bool {
	if !dps.isInCommentBlock || !dps.isInPatchBlock {
		return false
	}
	keys := map[string]bool{
		"@end ": true,
		"@end*": true,
	}
	keyLen := 5
	if i+keyLen > len([]rune(*str)) {
		return false
	}
	if c == '@' {
		if _, isContain := keys[(*str)[i:i+keyLen]]; isContain {
			return true
		}
	}
	return false
}

// getIDs reads key of the query piece.
// Query piece key is described after "@start" keyword.
func (dps *domainParser) getID(i int, c rune, str *string) (string, error) {
	s := strings.Replace((*str)[i:], "	", " ", -1)
	for {
		if !strings.Contains(s, "  ") {
			break
		}
		s = strings.Replace(s, "  ", " ", -1)
	}
	words := strings.Split(s, " ")
	if len(words) < 2 {
		return "", getIDError
	}
	return strings.Split(words[1], "*/")[0], nil
}
