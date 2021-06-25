package js

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
)

type UtilAPI struct{}

func (self FormatAPI) StringToBytes(string_ string) []byte {
	return util.StringToBytes(string_)
}

// Another way to achieve this in JavaScript: String.fromCharCode.apply(null, bytes)
func (self FormatAPI) BytesToString(bytes []byte) string {
	return util.BytesToString(bytes)
}

// Encode bytes as base64
func (self FormatAPI) Btoa(bytes []byte) string {
	return util.ToBase64(bytes)
}

// Decode base64 to bytes
func (self FormatAPI) Atob(b64 string) ([]byte, error) {
	return util.FromBase64(b64)
}

func (self UtilAPI) DeepCopy(value ard.Value) ard.Value {
	return ard.Copy(value)
}

func (self UtilAPI) DeepEquals(a ard.Value, b ard.Value) bool {
	return ard.Equals(a, b)
}

func (self UtilAPI) IsType(value ard.Value, type_ string) (bool, error) {
	// Special case whereby an integer stored as a float type has been optimized to an integer type
	if (type_ == "!!float") && ard.IsInteger(value) {
		return true, nil
	}

	if validate, ok := ard.TypeValidators[ard.TypeName(type_)]; ok {
		return validate(value), nil
	} else {
		return false, fmt.Errorf("unsupported type: %s", type_)
	}
}

func (self UtilAPI) Hash(value ard.Value) (string, error) {
	if hash, err := hashstructure.Hash(value, hashstructure.FormatV2, nil); err == nil {
		return strconv.FormatUint(hash, 10), nil
	} else {
		return "", err
	}
}

func (self UtilAPI) Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (self UtilAPI) Now() time.Time {
	return time.Now()
}

func (self UtilAPI) Mutex() *sync.Mutex {
	return new(sync.Mutex)
}

var onces map[string]*sync.Once = make(map[string]*sync.Once)
var oncesLock sync.Mutex

func (self UtilAPI) Once(name string, value goja.Value) error {
	if call, ok := goja.AssertFunction(value); ok {
		var once *sync.Once

		oncesLock.Lock()
		var ok bool
		if once, ok = onces[name]; !ok {
			once = new(sync.Once)
			onces[name] = once
		}
		oncesLock.Unlock()

		once.Do(func() {
			if _, err := call(nil); err != nil {
				log.Errorf("%s", err.Error())
			}
		})
		return nil
	} else {
		return fmt.Errorf("not a \"function\": %T", value)
	}
}

// Goroutine
func (self UtilAPI) Go(value goja.Value) error {
	if call, ok := goja.AssertFunction(value); ok {
		go func() {
			if _, err := call(nil); err != nil {
				log.Errorf("%s", err.Error())
			}
		}()
		return nil
	} else {
		return fmt.Errorf("not a \"function\": %T", value)
	}
}
