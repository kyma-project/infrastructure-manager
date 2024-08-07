package shoot

import (
	"fmt"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type ExtensionMatcher struct {
	toMatch interface{}
	fails   []string
}

func NewExtensionMatcher(i interface{}) types.GomegaMatcher {
	return &ExtensionMatcher{
		toMatch: i,
	}
}

func (m *ExtensionMatcher) Match(actual interface{}) (success bool, err error) {

	if err != nil {
		return false, err
	}

	if actual == nil && m.toMatch != nil {
		return false, errors.New("actual is nil")
	}

	if actual != nil && m.toMatch == nil {
		return false, errors.New("expected is nil")
	}

	if actual == nil && m.toMatch == nil {
		return true, nil
	}

	aExtensions, err := getExtension(actual)
	findExtension := func(extensions []v1beta1.Extension, extensionToFind string) v1beta1.Extension {
		for _, extension := range extensions {
			if extension.Type == extensionToFind {
				return extension
			}
		}
		return v1beta1.Extension{}
	}

	eExtensions, err := getExtension(m.toMatch)
	if err != nil {
		return false, err
	}

	for i, eExtension := range eExtensions {
		aExtension := findExtension(aExtensions, eExtension.Type)
		matcher := gomega.BeComparableTo(eExtension.Type)
		ok, err := matcher.Match(aExtension.Type)
		if err != nil {
			return false, err
		}

		if !ok {
			msg := fmt.Sprintf("spec/Extensions[%d]: %s", i, matcher.FailureMessage(aExtension))
			m.fails = append(m.fails, msg)
		}

	}

	return false, nil
}

func (m *ExtensionMatcher) NegatedFailureMessage(_ interface{}) string {
	return "expected should not equal actual"
}

func (m *ExtensionMatcher) FailureMessage(_ interface{}) string {
	return strings.Join(m.fails, "\n")
}

func getExtension(i interface{}) (shoot []v1beta1.Extension, err error) {
	if i == nil {
		return []v1beta1.Extension{}, fmt.Errorf("invalid value nil")
	}

	switch v := i.(type) {
	case []v1beta1.Extension:
		return v, nil

	default:
		return []v1beta1.Extension{}, fmt.Errorf(`%w: %s`, errInvalidType, reflect.TypeOf(v))
	}
}
