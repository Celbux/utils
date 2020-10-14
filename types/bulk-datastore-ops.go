package types

import "cloud.google.com/go/datastore"

// Maps generic interface to their type. Must be passed when writing to Datastore
const (
	TypeWalletData = iota + 1
	TypeVoucher
	TypePair
)

//ProcessRequest contains all information required for this instance to bulk write
//Kind is used for Datastore grouping
//Type is used to map interface to their type when actually writing to datastore
//EntityChunks contains all chunks this instance must process
type ProcessRequest struct {
	Kind         string
	Type         int
	EntityKeys	 []*datastore.Key
	EntityChunks []EntityChunk
}

//EntityChunk includes a chunk of entities to store. MAX 500 per
type EntityChunk struct {
	Entities interface{}
}