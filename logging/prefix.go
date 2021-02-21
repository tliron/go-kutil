package logging

import (
	"sort"
	"strings"

	loggingpkg "github.com/op/go-logging"
)

//
// PrefixLeveledBackend
//

type PrefixLeveledBackend struct {
	wrapped      loggingpkg.LeveledBackend
	prefixLevels []prefixLevel
}

type prefixLevel struct {
	prefix string
	level  loggingpkg.Level
}

func NewPrefixLeveledBackend(wrapped loggingpkg.LeveledBackend) *PrefixLeveledBackend {
	return &PrefixLeveledBackend{
		wrapped: wrapped,
	}
}

// logging.Leveled interface

func (self *PrefixLeveledBackend) GetLevel(module string) loggingpkg.Level {
	for _, prefixLevel := range self.prefixLevels {
		if strings.HasPrefix(module, prefixLevel.prefix) {
			return prefixLevel.level
		}
	}

	return self.wrapped.GetLevel(module)
}

func (self *PrefixLeveledBackend) SetLevel(level loggingpkg.Level, module string) {
	if strings.HasSuffix(module, "*") {
		self.prefixLevels = append(self.prefixLevels, prefixLevel{
			prefix: module[:len(module)-1],
			level:  level,
		})

		// Sort in reverse so that the more specific (=longer) prefixes come first
		sort.Slice(self.prefixLevels, func(i int, j int) bool {
			return strings.Compare(self.prefixLevels[i].prefix, self.prefixLevels[j].prefix) == 1
		})
	} else {
		self.wrapped.SetLevel(level, module)
	}
}

func (self *PrefixLeveledBackend) IsEnabledFor(level loggingpkg.Level, module string) bool {
	return level <= self.GetLevel(module)
}

// logging.Backend interface

func (self *PrefixLeveledBackend) Log(level loggingpkg.Level, callDepth int, record *loggingpkg.Record) error {
	return self.wrapped.Log(level, callDepth, record)
}
