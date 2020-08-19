package sqlpatchwork

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// hasOnOffParseDone have filepathes which have already parsed.
var hasOnOffParseDone = map[string]*parseResult{}

// onOffParseFile parses file and save the result to hasParseDone.
func onOffParseFile(path string) (*parseResult, error) {
	if pr, ok := hasOnOffParseDone[path]; ok {
		return pr, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	dps := newDomainParser()
	dps.defaultValue = dps.onOffDefaultValue
	dps.customParseID = dps.onOffCustomParseID
	dps.checkEndedCorrect = dps.onOffCheckEndedCorrect
	pr, err := dps.parse(reader)
	hasOnOffParseDone[path] = pr
	return pr, err
}

func (dps *domainParser) onOffDefaultValue() (tmpIDs []string, queryPieceIDs map[string]bool) {
	tmpIDs = []string{onoff_default_id}
	queryPieceIDs = map[string]bool{onoff_default_id: true}
	return
}

// onOffCustomParseID splits ID by "/".
// Eg: "key1/key2" => [key1, key2]
func (dps *domainParser) onOffCustomParseID(tmpID string) ([]string, error) {
	return strings.Split(tmpID, "/"), nil
}

//checkEndedCorrect gets whether all "@start" closed by "@end" or not.
func (dps *domainParser) onOffCheckEndedCorrect() error {
	if len(dps.tmpIDs) == 1 && dps.tmpIDs[0] == onoff_default_id {
		return nil
	}
	return errors.New("'@end' is missing.")
}
