package shoot

import (
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

var _ = Describe(":: extension matcher :: ", func() {
	var empty []v1beta1.Extension

	DescribeTable(
		"checking invalid args :: ",
		testExtensionInvalidArgs,
		Entry("when actual is nil", nil, empty),
		Entry("when expected is nil", "", nil),
	)

	//disabledTrue := true
	disabledFalse := false

	DescribeTable(
		"checking results :: ",
		testExtensionResults,
		Entry(
			"should match empty and zero values",
			nil,
			empty,
			true,
		),
		Entry(
			"should match identical tables",
			newExtensionTable(newExtension("test", &disabledFalse)),
			newExtensionTable(newExtension("test", &disabledFalse)),
			true,
		),
		Entry(
			"should detect extension missing",
			newExtensionTable(newExtension("test1", &disabledFalse), newExtension("test2", &disabledFalse)),
			newExtensionTable(newExtension("test", &disabledFalse)),
			false,
		),
		Entry(
			"should detect redundant extension",
			newExtensionTable(newExtension("test", &disabledFalse)),
			newExtensionTable(newExtension("test1", &disabledFalse), newExtension("test2", &disabledFalse)),
			false,
		),
	)
})

func newExtension(name string, disabled *bool) v1beta1.Extension {
	return v1beta1.Extension{
		Type:     name,
		Disabled: disabled,
	}
}

func newExtensionTable(extensions ...v1beta1.Extension) []v1beta1.Extension {
	return extensions
}

func testExtensionInvalidArgs(actual, expected interface{}) {
	matcher := NewExtensionMatcher(expected)
	_, err := matcher.Match(actual)
	Expect(err).To(HaveOccurred())
}

func testExtensionResults(actual, expected interface{}, expectedMatch bool) {
	matcher := NewExtensionMatcher(expected)
	actualMatch, err := matcher.Match(actual)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(actualMatch).Should(Equal(expectedMatch), matcher.FailureMessage(actual))
}
