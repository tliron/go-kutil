package format

import (
	"github.com/tliron/kutil/ard"
)

func ToCompatibleJSON(value ard.Value) ard.Value {
	value, _ = Normalize(value)
	value = ard.ToCompatibleJSON(value)
	return value
}
