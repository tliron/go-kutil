Logging
=======

A common API for structured and unstructured logging with support for runtime-pluggable backends.

The main design goal is to allow for libraries that can integrate well with others without having
to change their logging code. For example, a library can plug into the klog backend if it's being
used together with the Kubernetes client.

A secondary design goal is to provide a home for a full-featured unstructured logging library,
which we here call the "simple" backend. It supports rich formatting, including ANSI coloring when
logging to a terminal, including on Windows.

API features:

* Fine-grained control over verbosity via hierarchical log names. For example, "engine.parser.background"
  inherits from "engine.parser", which in turn inherits from "engine". The empty name is the root of the
  hierarchy. Each name's default verbosity is that of its parent, which you can then override with
  `logging.SetMaxLevel`.
* Support for call stack depth. This can be used by a backend (for example, klog) to find out where in
  the code the logging happened.
* No need to create loggers. The true API entrypoint is the global function `logging.NewMessage`, which
  you provide with a name and a level. The default logger type is just a convenient wrapper around it that
  provides the familiar unstructured functions.
* That said, the `logging.Logger` type is an interface, allowing you to more easily switch
  implementations. For example, you can assign the `logging.MOCK_LOGGER` to disable a logger without changing
  your code. (It's unfortunate that Go's [`Logger`](https://pkg.go.dev/log#Logger) type is a struct.)

Basic Usage
-----------

The easiest way to plug in a backend is to anonymously import the correct sub-package into your program's
main package:

```go
import (
    _ "github.com/tliron/kutil/logging/simple"
)
```

This should enable the backend with sensible defaults. Specifically it will log to stderr with
verbosity at the "notice" max level.

Example of structured logging:

```go
import (
    "github.com/tliron/kutil/logging"
    _ "github.com/tliron/kutil/logging/simple"
    "github.com/tliron/kutil/util"
)

func main() {
    if m := logging.NewMessage([]string{"engine", "parser"}, logging.Error, 0); m != nil {
        m.Set("message", "Hello world!").Set("myfloat", 10.2).Send()
    }
    util.Exit(0)
}
```

Note that `logging.NewMessage` will return nil if the message cannot be created, for example if the
message level is higher than the max level for that name.

`Set` can accept any key and value, but two keys are recognized by the API:

* `message`: The main description of the message. This is the key used by unstructured logging.
* `scope`: An optional identifier that can be used to group messages, making them easier to filter
  (e.g. by grep). Backends may handle this specially. Unstructured backends may, for example, add
  it as a bracketed prefix for messages.

Also note that calling `util.Exit(0)` to exit your program is not absolutely necessary, however
it's good practice because it makes sure to flush buffered log messages for some backends.

Example of unstructured logging:

```go
import (
    "github.com/tliron/kutil/logging"
    _ "github.com/tliron/kutil/logging/simple"
    "github.com/tliron/kutil/util"
)

var log = logging.GetLogger("engine.parser")

func main() {
    log.Errorf("Hello %s!", "world")
    util.Exit(0)
}
```

Use conditional logging to optimize for costly unstructured message creation, e.g.:

```go
if log.AllowLevel(logging.Debug) {
    log.Debugf("Status is: %s", getStatusFromDatabase())
}
```

The scope logger can be used to automatically set the "scope" key for another logger. It
automatically detects nesting, in which case it appends the new scope separated by a ".",
e.g.:

```go
var log = logging.GetLogger("engine.parser")
var validationLog = logging.NewScopeLogger(log, "validation")
var syntaxLog = logging.NewScopeLogger(validationLog, "syntax")

func main() {
    // Nested scope will be "validation.syntax"
    syntaxLog.Errorf("Hello %s!", "world")
    ...
}
```

Configuration
-------------

All backends can be configured via the same API. For example, to increase verbosity and write
to a file:

```go
func main() {
    path := "myapp.log"
    logging.Configure(1, &path)
    ...
}
```

Backends may also have their own (non-portable) configuration APIs.

You can set the max level (verbosity) using either the global API or a logger. For
example, here is a way to make all logging verbose by default, except for one name:

```go
func init() {
    logging.SetMaxLevel(nil, logging.Debug) // nil = the root
    logging.SetMaxLevel([]string{"engine", "parser"}, logging.Error)
}
```

Note that descendents of "engine.parser", e.g. "engine.parser.analysis", would acquire
its log levels rather than the root's. Here's the same effect using loggers:

```go
var rootLog = logging.GetLogger("")
var parserLog = logging.GetLogger("engine.parser")

func init() {
    rootLog.SetMaxLevel(logging.Debug)
    parserLog.SetMaxLevel(logging.Error)
}
```

It's important to note that the configuration APIs are not thread safe. This includes
`Configure` and `SetMaxLevel`. Thus, make sure to get all your configuration done before
you start sending log messages. A good place for this is `init()` or `main()` functions.

Color
-----

For the simple backend you must explicitly attempt to enable ANSI color if desired. Note that
if it's unsupported by the terminal then no ANSI codes will be sent (unless you force it via
`terminal.EnableColor(true)`):

```go
import (
    "github.com/tliron/kutil/logging"
    _ "github.com/tliron/kutil/logging/simple"
    "github.com/tliron/kutil/terminal"
    "github.com/tliron/kutil/util"
)

func main() {
    terminal.EnableColor(false)
    logging.GetLogger("engine.parser").Error("Hello world!") // errors are in red
    util.Exit(0)
}
```