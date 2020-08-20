package sqlpatchwork

import (
	"bufio"
	"errors"
	"os"
)

// hasSimpleParseDone have filepathes which have already parsed.
var hasSimpleParseDone = map[string]*parseResult{}

// simpleParseFile parses file and save the result to hasParseDone.
func simpleParseFile(path string) (*parseResult, error) {
	if pr, ok := hasSimpleParseDone[path]; ok {
		return pr, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	dps := newDomainParser()
	dps.di2Simple()
	pr, err := dps.parse(reader)
	hasSimpleParseDone[path] = pr
	return pr, err
}

// di2Simple injects dependencey as a simplePatchwork Parser.
func (dps *domainParser) di2Simple() {
	dps.defaultValue = dps.simpleDefaultValue
	dps.customParseID = dps.simpleCustomParseID
	dps.checkEndedCorrect = dps.simpleCheckEndedCorrect
	dps.appendQP = dps.simpleAppendQp
}

// simpleDefaultValue sets default value to tmpIDs and queryPieceIDs.
// In this function, not set to parser's field.
func (dps *domainParser) simpleDefaultValue() (tmpIDs []string, queryPieceIDs map[string]bool) {
	tmpIDs = nil
	queryPieceIDs = make(map[string]bool)
	return
}

// simpleCustomParseID checks whether ID is duplicated or not.
// And return on element tmpID.
func (dps *domainParser) simpleCustomParseID(tmpID string) ([]string, error) {
	if _, isDuplicated := dps.queryPieceIDs[tmpID]; isDuplicated {
		return nil, errors.New("duplicated")
	}
	return []string{tmpID}, nil
}

//checkEndedCorrect gets whether all "@start" closed by "@end" or not.
func (dps *domainParser) simpleCheckEndedCorrect() error {
	if len(dps.tmpIDs) == 0 {
		return nil
	}
	return errors.New("'@end' is missing.")
}

//simpleAppendQp appends query piece to simple query pieace field.
//In simplepatchwork mode, tmpIDs have only 1 element. So we use tmpIDs[0].
//It is guaranteed in simpleCustomParseID that number of tmpIDs is one.
func (dps *domainParser) simpleAppendQp() {
	dps.simpleQueryPieces[dps.tmpIDs[0]] = dps.queryBuf
	return
}
