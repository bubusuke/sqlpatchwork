package sqlpatchwork

// onOffQP represets pieces of sql for OnOffPatchwork
type onOffQPs []onOffQP

// onOffQP represents a piece of sql for OnOffPatchwork
type onOffQP struct {
	//IDs is keys of query piece. Multiple keys can be assigned to a query piece.
	IDs []string
	//query is content of query piece.
	query []byte
}

// simpleQPs represents pieces of sql for SimplePatchwork
type simpleQPs map[string][]byte
