package matcher

import (
	"github.com/rs/xid"
	"mock-server/model"
	"sort"
	"sync"
)

type Engine struct {
	lock         sync.RWMutex
	expectations []*model.Expectation
}

func (e *Engine) Expectations() []*model.Expectation {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.expectations
}

func (e *Engine) AddExpectations(list []*model.Expectation) {
	hasPriority := false
	for _, exp := range list {
		if exp.Id == "" {
			exp.Id = xid.New().String()
		}
		if exp.Priority > 0 {
			hasPriority = true
		}
		NormalizeRequest(&exp.HttpRequest)
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	e.expectations = append(e.expectations, list...)

	if hasPriority {
		sort.Slice(e.expectations, func(i, j int) bool {
			return e.expectations[i].Priority > e.expectations[j].Priority
		})
	}
}

func (e *Engine) ClearBy(exp *model.HttpRequest) []string {
	NormalizeRequest(exp)

	e.lock.Lock()
	defer e.lock.Unlock()

	var (
		newExpectations []*model.Expectation
		excludedIds     []string
	)

	for _, actualExp := range e.expectations {
		if MatchRequestByRequest(exp, &actualExp.HttpRequest) {
			excludedIds = append(excludedIds, actualExp.Id)
			continue
		}
		newExpectations = append(newExpectations, actualExp)
	}

	e.expectations = newExpectations
	return excludedIds
}

func (e *Engine) Reset() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.expectations = []*model.Expectation{}
}
