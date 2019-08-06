package matchers_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
)

var _ = Describe("CompareTo", func() {
	Context("when asserting that nil equals nil", func() {
		It("should error", func() {
			success, err := (&CompareToMatcher{}).Match(nil)

			Expect(success).Should(BeFalse())
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("When comparing objects", func() {
		It("should do the right thing", func() {
			Expect(5).Should(CompareTo(5))
			Expect(5.0).Should(CompareTo(5.0))

			Expect(5).ShouldNot(CompareTo("5"))
			Expect(5).ShouldNot(CompareTo(5.0))
			Expect(5).ShouldNot(CompareTo(3))

			Expect("5").Should(CompareTo("5"))
			Expect([]int{1, 2}).Should(CompareTo([]int{1, 2}))
			Expect([]int{1, 2}).ShouldNot(CompareTo([]int{2, 1}))
			Expect([]byte{'f', 'o', 'o'}).Should(CompareTo([]byte{'f', 'o', 'o'}))
			Expect([]byte{'f', 'o', 'o'}).ShouldNot(CompareTo([]byte{'b', 'a', 'r'}))
			Expect(map[string]string{"a": "b", "c": "d"}).Should(CompareTo(map[string]string{"a": "b", "c": "d"}))
			Expect(map[string]string{"a": "b", "c": "d"}).ShouldNot(CompareTo(map[string]string{"a": "b", "c": "e"}))
		})

		It("should use their Equal method", func() {
			moment := time.Now()
			Expect(moment.In(time.FixedZone("UTC+2", 2*60*60))).Should(CompareTo(moment.In(time.FixedZone("UTC-8", -8*60*60))))
			Expect(moment).ShouldNot(CompareTo(moment.Add(3)))
		})

		It("should panic on unexported fields", func() {
			Expect(func() { Expect(errors.New("foo")).Should(CompareTo(errors.New("foo"))) }).To(Panic())
			Expect(func() { Expect(myCustomType{}).Should(CompareTo(myCustomType{})) }).To(Panic())

			Context("unless using a custom comparer", func() {
				cmpErrorMessage := cmp.Comparer(func(x, y error) bool { return x.Error() == y.Error() })
				Expect(errors.New("foo")).Should(CompareTo(errors.New("foo"), cmpErrorMessage))
				Expect(errors.New("foo")).Should(CompareTo(fmt.Errorf("foo"), cmpErrorMessage))
				Expect(errors.New("foo")).ShouldNot(CompareTo("foo", cmpErrorMessage))
				Expect(errors.New("foo")).ShouldNot(CompareTo(errors.New("bar"), cmpErrorMessage))
			})

			Context("unless unexported fields are allowed", func() {
				cmpCustomType := cmp.AllowUnexported(myCustomType{})
				Expect(myCustomType{}).Should(CompareTo(myCustomType{}, cmpCustomType))
				Expect(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).Should(CompareTo(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}, cmpCustomType))
				Expect(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(CompareTo(
					myCustomType{s: "bar", n: 3, f: 2.0, arr: []string{"a", "b"}}, cmpCustomType))
				Expect(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(CompareTo(
					myCustomType{s: "foo", n: 2, f: 2.0, arr: []string{"a", "b"}}, cmpCustomType))
				Expect(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(CompareTo(
					myCustomType{s: "foo", n: 3, f: 3.0, arr: []string{"a", "b"}}, cmpCustomType))
				Expect(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b"}}).ShouldNot(CompareTo(
					myCustomType{s: "foo", n: 3, f: 2.0, arr: []string{"a", "b", "c"}}, cmpCustomType))
			})
		})
	})

	Describe("failure messages", func() {
		It("shows differing types", func() {
			subject := CompareToMatcher{EqualMatcher: EqualMatcher{Expected: int64(5)}}
			failureMessage := subject.FailureMessage(int32(5))

			Expect(failureMessage).To(MatchRegexp(
				`(?m:` +
					`[-]{3} Actual` + `[\pZ\s]*` +
					`[+]{3} Expected` +
					`).*`))
			Expect(failureMessage).To(MatchRegexp(
				`.*(?m:` +
					`[-][\pZ\s]*int32[(]5[)],` + `[\pZ\s]*` +
					`[+][\pZ\s]*int64[(]5[)],` +
					`).*`))
		})

		It("shows two strings simply when they are short", func() {
			subject := CompareToMatcher{EqualMatcher: EqualMatcher{Expected: "eric"}}
			failureMessage := subject.FailureMessage("erin")

			Expect(failureMessage).To(MatchRegexp(
				`(?m:` +
					`[-]{3} Actual` + `[\pZ\s]*` +
					`[+]{3} Expected` +
					`).*`))
			Expect(failureMessage).To(MatchRegexp(
				`.*(?m:` +
					`[-][\pZ\s]*"erin",` + `[\pZ\s]*` +
					`[+][\pZ\s]*"eric",` +
					`).*`))
		})

		It("shows the exact point where two long strings differ", func() {
			stringWithB := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			stringWithZ := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaazaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

			subject := CompareToMatcher{EqualMatcher: EqualMatcher{Expected: stringWithZ}}
			failureMessage := subject.FailureMessage(stringWithB)

			Expect(failureMessage).To(MatchRegexp(
				`.*(?m:` +
					`[-]{3} Actual` + `[\pZ\s]*` +
					`[+]{3} Expected` +
					`).*`))
			Expect(failureMessage).To(MatchRegexp(
				`.*(?m:` +
					`[-][\pZ\s]*"b",` + `[\pZ\s]*` +
					`[+][\pZ\s]*"z",` +
					`).*`))
		})

		It("shows the lines where two multi-line strings differ", func() {
			expect := "abcdef\n111111\naaaaaa\n222222\nbbbbbb\n333333\ncccccc\n444444\ndddddd\n555555"
			actual := "abcdef\n123456\naaaaaa\n222222\nbbbbbb\n333333\ncccccc\n456789\ndddddd\n555555"

			subject := CompareToMatcher{EqualMatcher: EqualMatcher{Expected: expect}}
			failureMessage := subject.FailureMessage(actual)

			Expect(failureMessage).To(MatchRegexp(
				`(?m:` +
					`[-]{3} Actual` + `[\pZ\s]*` +
					`[+]{3} Expected` +
					`).*`))
			Expect(failureMessage).To(MatchRegexp(
				`.*(?m:` +
					`[-][\pZ\s]*"123456",` + `[\pZ\s]*` +
					`[+][\pZ\s]*"111111",` +
					`).*`))
			Expect(failureMessage).To(MatchRegexp(
				`.*(?m:` +
					`[-][\pZ\s]*"456789",` + `[\pZ\s]*` +
					`[+][\pZ\s]*"444444",` +
					`).*`))
		})
	})
})
