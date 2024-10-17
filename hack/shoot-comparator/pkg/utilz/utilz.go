package utilz

import (
	"fmt"
	"reflect"

	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/errors"
	"sigs.k8s.io/yaml"
)

func Get[T any](v interface{}) (T, error) {
	var result T
	if v == nil {
		return result, errors.ErrNilValue
	}

	switch typedV := v.(type) {
	case string:
		err := yaml.Unmarshal([]byte(typedV), &result)
		return result, err

	case T:
		return typedV, nil

	case *T:
		if typedV == nil {
			return result, nil
		}
		return *typedV, nil

	default:
		return result, fmt.Errorf(`%w: %s`, errors.ErrInvalidType, reflect.TypeOf(typedV))
	}
}
