package sqlpatchwork

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//simplePatchwork imprements Sqlpatchwork.
type simplePatchwork struct {
	sqlName       string
	queryPieceIDs map[string]bool
	queryPieces   simpleQPs
	applyIDOrder  []string
}

// NewSimplePatchwork parses sql file and return Sqlpatchwork.
func NewSimplePatchwork(sqlFilePath string) (Sqlpatchwork, error) {
	pr, err := simpleParseFile(sqlFilePath)
	if err != nil {
		return nil, err
	}
	return &simplePatchwork{
		sqlName:       sqlFilePath,
		queryPieceIDs: pr.queryPieceIDs,
		queryPieces:   pr.simpleQueryPieces,
		applyIDOrder:  nil,
	}, nil
}

// NewSimplePWSkipPrs requires query pieces and return Sqlpatchwork.
// This is the function which skipped sql file parsing process.
// You can start SQL patchwork without preparing sql files by using this method.
func NewSimplePWSkipPrs(sqlName string, simpleQueryPieces map[string]string) Sqlpatchwork {
	sQPs := make(map[string][]byte)
	queryPieceIDs := make(map[string]bool)
	for key, val := range simpleQueryPieces {
		sQPs[key] = []byte(val)
		queryPieceIDs[key] = true
	}
	return &simplePatchwork{
		sqlName:       sqlName,
		queryPieceIDs: queryPieceIDs,
		queryPieces:   sQPs,
		applyIDOrder:  nil,
	}
}

//AddQueryPieceToBuild adds query-pieces to BuildQuery target.
//When ID is not found, return error.
func (spw *simplePatchwork) AddQueryPiecesToBuild(IDs ...string) error {
	//check
	for _, ID := range IDs {
		if _, ok := spw.queryPieceIDs[ID]; !ok {
			return errors.New(fmt.Sprintf("Failure to add. The queryPieceID is not exists. queryPieceID: '%v'\n", ID))
		}
	}
	for _, ID := range IDs {
		spw.applyIDOrder = append(spw.applyIDOrder, ID)
	}
	return nil
}

//BuildQuery builds query to join query-pieces and return query.
func (spw *simplePatchwork) BuildQuery() (query string) {
	queryBuf := make([]byte, 0, 4096)
	loopCount := make(map[string]int)
	// build
	for _, ID := range spw.applyIDOrder {
		loopNo := loopCount[ID]
		loopCount[ID]++
		queryBuf = append(queryBuf, []byte(" ")...)
		queryBuf = append(queryBuf, convertLoopNo(spw.queryPieces[ID], loopNo)...)
	}
	// trim and decrease spaces.
	query = strings.Trim(string(queryBuf), " ")
	for {
		if !strings.Contains(query, "  ") {
			break
		}
		query = strings.Replace(query, "  ", " ", -1)
	}
	return
}

//BuildQueryWithTraceDesc builds query to join query-pieces and add sqlfilename and applied query piese IDs to query as comment and return query.
func (spw *simplePatchwork) BuildQueryWithTraceDesc() (query string) {
	query = spw.BuildQuery()
	// Describe apply condition to trace.
	traceDesc := fmt.Sprintf(" /* %v %v */ ", spw.sqlName, spw.TargetIDs())
	query = strings.Replace(query, " ", traceDesc, 1)
	return
}

//targetIDs gets BuildQuery targets.
func (spw *simplePatchwork) TargetIDs() []string {
	return spw.applyIDOrder
}

// convertLoopNo converts "@@" to loopNo
func convertLoopNo(query []byte, loopNo int) []byte {
	return []byte(strings.Replace(string(query), loopNoIndicater, strconv.Itoa(loopNo), -1))
}
