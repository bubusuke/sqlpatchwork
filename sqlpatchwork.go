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
	TargetIDs() []string

	//BuildQuery builds query to join query-pieces and return query.
	BuildQuery() string

	//BuildQueryWithTraceDesc builds query to join query-pieces and add sqlfilename and applied query piese IDs to query as comment.
	BuildQueryWithTraceDesc() string
}

const loopNoIndicater = "@@"

// LoopNoAttach converts "@@" to loopNo
func LoopNoAttach(name string, loopNo int) string {
	return strings.Replace(name, loopNoIndicater, strconv.Itoa(loopNo), -1)
}
