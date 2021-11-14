package matcher

import (
	"github.com/stretchr/testify/assert"
	"mock-server/model"
	"testing"
)

func getExpectationIDs(list []*model.Expectation) (res []string) {
	for _, item := range list {
		res = append(res, item.Id)
	}
	return res
}

func TestEngine_AddExpectations(t *testing.T) {
	e := &Engine{}

	list := []*model.Expectation{
		{Id: "2", Priority: 2},
		{Id: "1", Priority: 1},
		{Id: "3", Priority: 3},
	}
	e.AddExpectations(list)

	assert.Equal(t, []string{"3", "2", "1"}, getExpectationIDs(e.expectations), "it should sort expectations in desc order")

	list2 := []*model.Expectation{
		{Id: "5", Priority: 2},
		{Id: "6", Priority: 6},
	}
	e.AddExpectations(list2)

	assert.Equal(t, []string{"6", "3", "2", "5", "1"}, getExpectationIDs(e.expectations), "and keep it sorted in next iterations")

	list3 := []*model.Expectation{
		{Id: "7", Priority: 0},
	}
	e.AddExpectations(list3)

	assert.Equal(t, []string{"6", "3", "2", "5", "1", "7"}, getExpectationIDs(e.expectations), "expecations without priority does not change sort order")
}
