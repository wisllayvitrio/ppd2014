package space

import (
	"github.com/wisllayvitrio/ppd2014/index"
)

type searchIndex struct {
	numAttributes int
	tupleIndex []*index.Index
}

func NewSearchIndex(numAttributes int) *searchIndex {
	search := &searchIndex{
		numAttributes,
		make([]*index.Index, numAttributes),
	}

	for i := 0; i < numAttributes; i++ {
		search.tupleIndex[i] = index.NewIndex(tupleIndexDuration)
	}

	return search
}
