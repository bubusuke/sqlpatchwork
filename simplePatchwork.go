package sqlpatchwork

import (
	"errors"
	"fmt"
	"strings"
)

//simplePatchwork imprements Sqlpatchwork.
type simplePatchwork struct {
	sqlName       string
	queryPieceIDs map[string]bool
	queryPieces   []queryPiece
	applyIDOrder  []string
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
		queryBuf = append(queryBuf, addLoopNoTobindName(spw.getQueryPieces(ID), loopNo)...)
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
	traceDesc := fmt.Sprintf(" /* %v %v */ ", spw.sqlName, spw.targetIDs())
	query = strings.Replace(query, " ", traceDesc, 1)
	return
}

//targetIDs gets BuildQuery targets.
func (spw *simplePatchwork) targetIDs() []string {
	return spw.applyIDOrder
}

//getQueryPieces gets query-piece.
//In simplePatchwork, querypiece ID is unique key (guaranteed at parse-process).
func (spw *simplePatchwork) getQueryPieces(ID string) []byte {
	for _, qp := range spw.queryPieces {
		if qp.IDs[0] == ID {
			return qp.query
		}
	}
	return nil
}
