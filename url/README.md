URLs
====

Usage example:

```go
import (
    urlpkg "github.com/tliron/kutil/url"
)

func ReadAll(url string, format string) ([]byte, error) {
    context := urlpkg.NewContext()
    defer context.Release()

    if url_, err = urlpkg.NewURL(url, context); err == nil {
        if reader, err := url_.Open(); err == nil {
            defer reader.Close()
            return io.ReadAll(reader)
        } else {
            return nil, err
        }
    } else {
        return nil, err
    }
}
```

file:
-----

http: and https:
----------------

zip:
----

git:
----

docker:
-------

internal:
---------
