Agnostic Raw Data (ARD)
=======================

What is "agnostic raw data"?

Agnostic?
---------

Comprising primitives (string, integer, float, boolean, null) and structures (map, list). It's
agnostic because it can be trivially represented in practically any language or platform, and also
because it can be transmitted in a wide variety of formats.

Note that map keys do not have to be strings and indeed can be arbitrarily complex, not just
strings. Most map implementations in most programming languages allow for this as long as the key is
hashable.

Also not that precision of integers and floats is out of scope, and thus we normally do not have
to distinguish between signed and unsigned integers.

Raw?
----

Data validations is out of scope. There's no schema.

Caveats
-------

Some caveats and limitations for programming languages:

### Go

Unfortunately, the most popular Go YAML parser does not easily support arbitrarily complex keys
(see this [issue](https://github.com/go-yaml/yaml/issues/502)). We provide an independent library,
[yamlkeys](https://github.com/tliron/yamlkeys), to make this easier.

### JavaScript

The JavaScript language doesn't have native support for integers. However, this limitation can be
overcome if you're able to maintain precision and distinction, e.g. by wrapping the number in a
custom object.

Some caveats and limitations for transmission formats:

### YAML

YAML supports a rich set of primitive types (when it includes the common
[JSON schema](https://yaml.org/spec/1.2/spec.html#id2803231)), so ARD will survive a round trip
to YAML.

Note that some YAML 1.1 implementations support ordered maps
([!!omap](https://yaml.org/type/omap.html) vs. !!map). These will lose their order when converted
to ARD, so it's best to standardized on arbitrary order (!!map). YAML 1.2 does not support !!omap
by default, so this use case may be less and less common.

### JSON

JSON can be read into ARD. However, because JSON has fewer types and more limitations than YAML (no
integers, only floats; map keys can only be strings), ARD will lose some type information when
translated into JSON.

We can overcome this challenge by extending JSON with some conventions for encoding extra types.
See [our conventions here](json.go). Our implementation is in Go, but it should not be too difficult
to support them in another programming languages.

### XML

XML does not have a type system. Arbitrary XML cannot be parsed into ARD. 

However, we support [certain conventions](xml.go) that enforce such compatibility. Our
implementation is in Go, but it should not be too difficult to support them in another programming
languages.
