package sql

// WithPredicateExpectedError
func WithPredicateExpectedError(step, stage int, expectedErr error) bool {
	return step == stage && expectedErr != nil
}
