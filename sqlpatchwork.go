package sqlpatchwork

import (
	"strconv"
	"strings"
)

// Sqlpatchwork represent behaver of this tool.
type Sqlpatchwork interface {

	//AddQueryPieceToBuild adds query-pieces to BuildQuery target.
	AddQueryPiecesToBuild(...string) error

	//targetIDs gets BuildQuery targets.
	targetIDs() []string

	//BuildQuery builds query to join query-pieces and return query.
	BuildQuery() string

	//BuildQueryWithTraceDesc builds query to join query-pieces and add sqlfilename and applied query piese IDs to query as comment.
	BuildQueryWithTraceDesc() string
}

// NewOnOffpatchwork parses sql file and return.
func NewOnOffpatchwork(sqlFilePath string) (*onOffPatchwork, error) {
	pr, err := onOffParseFile(sqlFilePath)
	if err != nil {
		return nil, err
	}
	return &onOffPatchwork{
		sqlName:       sqlFilePath,
		queryPieceIDs: pr.queryPieceIDs,
		queryPieces:   pr.queryPieces,
		applyIDs:      map[string]bool{onoff_default_id: true},
	}, nil
}

// NewSimplePatchwork parses sql file and return.
func NewSimplePatchwork(sqlFilePath string) (*simplePatchwork, error) {
	pr, err := simpleParseFile(sqlFilePath)
	if err != nil {
		return nil, err
	}
	return &simplePatchwork{
		sqlName:       sqlFilePath,
		queryPieceIDs: pr.queryPieceIDs,
		queryPieces:   pr.queryPieces,
		applyIDOrder:  nil,
	}, nil
}

const loopNoIndicater = "@@"

func LoopNoAttach(name string, loopNo int) string {
	return strings.Replace(name, loopNoIndicater, strconv.Itoa(loopNo), -1)
}
func addLoopNoTobindName(query []byte, loopNo int) []byte {
	return []byte(strings.Replace(string(query), loopNoIndicater, strconv.Itoa(loopNo), -1))
}