package sqlpatchwork

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// default id will be assigned when query piece is not described to be assigned ID.
const onoff_default_id = "__default"

//onOffpatchwork imprements Sqlpatchwork.
type onOffPatchwork struct {
	sqlName       string
	queryPieceIDs map[string]bool
	queryPieces   []queryPiece
	applyIDs      map[string]bool
}

//AddQueryPiecesToBuild adds query-pieces to BuildQuery target.
//When ID is not found, return error.
func (opw *onOffPatchwork) AddQueryPiecesToBuild(IDs ...string) error {
	//check
	for _, ID := range IDs {
		if _, ok := opw.queryPieceIDs[ID]; !ok {
			return errors.New(fmt.Sprintf("Failure to add. The queryPieceID is not exists. queryPieceID: '%v'\n", ID))
		}
	}
	//add
	for _, ID := range IDs {
		opw.applyIDs[ID] = true
	}
	return nil
}

//BuildQuery builds query to join query-pieces and return query.
func (opw *onOffPatchwork) BuildQuery() (query string) {
	queryBuf := make([]byte, 0, 4096)
	// build
	for _, qp := range opw.queryPieces {
		for _, id := range qp.IDs {
			if isApply, hit := opw.applyIDs[id]; hit && isApply {
				queryBuf = append(queryBuf, qp.query...)
				break
			}
		}
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
func (opw *onOffPatchwork) BuildQueryWithTraceDesc() (query string) {
	query = opw.BuildQuery()
	// Describe apply condition to trace.
	traceDesc := fmt.Sprintf(" /* %v %v */ ", opw.sqlName, opw.TargetIDs())
	query = strings.Replace(query, " ", traceDesc, 1)
	return
}

//targetIDs gets BuildQuery targets.
func (opw *onOffPatchwork) TargetIDs() []string {
	IDs := []string{}
	for key, isApply := range opw.applyIDs {
		if isApply {
			IDs = append(IDs, key)
		}
	}
	//Fix the order of elements.
	sort.Slice(IDs, func(i, j int) bool { return IDs[i] < IDs[j] })
	return IDs
}
