package matcher

import (
	"github.com/rs/xid"
	"mock-server/model"
	"sort"
	"sync"
)

type Engine struct {
	lock         *sync.RWMutex
	expectations []*model.Expectation
}

func NewEngine() *Engine {
	return &Engine{
		lock: &sync.RWMutex{},
	}
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

func (e *Engine) ClearBy(exp *model.HttpRequest) {
	NormalizeRequest(exp)

	exclude := make(map[int]bool)
	for idx, actualExp := range e.expectations {
		if MatchRequestByRequest(exp, &actualExp.HttpRequest) {
			exclude[idx] = true
		}
	}

	newExpectations := make([]*model.Expectation, len(e.expectations)-len(exclude))

	addIdx := 0
	for idx, e := range e.expectations {
		if exclude[idx] {
			continue
		}
		newExpectations[addIdx] = e
		addIdx++
	}

	e.expectations = newExpectations
}

func (e *Engine) Reset() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.expectations = []*model.Expectation{}
}
