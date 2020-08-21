package sqlpatchwork

// simpleQPs represents pieces of sql for SimplePatchwork
type simpleQPs map[string][]byte

// onOffQP represents a piece of sql for OnOffPatchwork
type onOffQP struct {
	//iDs is keys of query piece. Multiple keys can be assigned to a query piece.
	iDs []string
	//query is content of query piece.
	query []byte
}

// OnOffQP return a onOffQueryPiece.
// When iDs is not set, "__defalt" will be set as id and the query piece always included in build query.
func OnOffQP(query string, iDs ...string) onOffQP {
	if iDs == nil {
		return onOffQP{
			iDs:   []string{onoff_default_id},
			query: []byte(query),
		}
	}
	return onOffQP{
		iDs:   iDs,
		query: []byte(query),
	}
}

// OnOffQPs just only return []onOffQueryPiece.
// This function is used to create arguments of NewOnOffPWSkipPrs.
func OnOffQPs(onOffQPs ...onOffQP) []onOffQP {
	return onOffQPs
}
