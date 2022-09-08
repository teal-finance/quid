package db

import (
	"fmt"

	"github.com/lib/pq"
)

// QueryResult :.
type QueryResult struct {
	Error    QueryError
	HasError bool
}

// QueryError :.
type QueryError struct {
	Message        string
	HasUserMessage bool
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("%t:%v: query error", e.HasUserMessage, e.Message)
}

func queryNoError() QueryResult {
	return QueryResult{HasError: false}
}

func queryError(err error) QueryResult {
	e, isPq := err.(*pq.Error)
	if isPq {
		return QueryResult{
			HasError: true,
			Error: QueryError{
				Message:        e.Message,
				HasUserMessage: true,
			},
		}
	}

	return QueryResult{
		HasError: true,
		Error: QueryError{
			Message:        e.Error(),
			HasUserMessage: false,
		},
	}
}
