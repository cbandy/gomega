package matchers

import "github.com/google/go-cmp/cmp"

type CompareToMatcher struct {
	EqualMatcher
	Options cmp.Options
}

func (matcher *CompareToMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return matcher.EqualMatcher.Match(actual)
	}
	return cmp.Equal(actual, matcher.Expected, matcher.Options...), nil
}

func (matcher *CompareToMatcher) FailureMessage(actual interface{}) (message string) {
	return "--- Actual\n+++ Expected\n" + cmp.Diff(actual, matcher.Expected, matcher.Options...)
}
