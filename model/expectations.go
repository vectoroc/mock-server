package model

import (
	"encoding/json"
)

// Expectations list of expectations
type Expectations struct {
	list []*Expectation
}

func (e Expectations) ToArray() []*Expectation {
	return e.list
}

func (e *Expectations) UnmarshalJSON(data []byte) error {
	var list []json.RawMessage

	if err := json.Unmarshal(data, &list); err != nil {
		// not array
		ex, err := unmarshalOneExpectation(data)
		if err != nil {
			return err
		}
		e.list = append(e.list, ex)
		return nil
	}

	e.list = make([]*Expectation, len(list))
	for i, item := range list {
		ex, err := unmarshalOneExpectation(item)
		if err != nil {
			return err
		}

		e.list[i] = ex
	}

	return nil
}

func (e *Expectations) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(e.list)

	return
}

func unmarshalOneExpectation(data []byte) (*Expectation, error) {
	ex := &Expectation{}

	err := json.Unmarshal(data, ex)
	if err != nil {
		return nil, err
	}

	//if ex.HttpResponse == nil &&
	//	ex.HttpResponseTemplate == nil &&
	//	ex.HttpResponseClassCallback == nil &&
	//	ex.HttpResponseObjectCallback == nil &&
	//	ex.HttpForward == nil &&
	//	ex.HttpForwardTemplate == nil &&
	//	ex.HttpForwardObjectCallback == nil &&
	//	ex.HttpForwardClassCallback == nil &&
	//	ex.HttpOverrideForwardedRequest == nil &&
	//	ex.HttpError == nil {
	//	return nil, ErrBadFormat
	//}

	return ex, nil
}
