Agnostic Raw Data (ARD)
=======================

What is "agnostic raw data"?

Agnostic?
---------

Comprising primitives (string, integer, float, boolean, null) and structures (map, list). It's
agnostic because it can be trivially represented in practically any language or platform, and also
because it can be transmitted in a wide variety of formats.

Note that some formats present limitations:

### YAML

YAML supports a rich set of primitive types, so ARD will survive a round trip to YAML.

One difference is that in YAML 1.1 maps can be ordered (!!omap vs. !!map) but ARD maps have
arbitrary order (always !!map) for widest compatibility. A round trip from YAML to ARD would thus
lose order. (YAML 1.2 does not include !!omap support by default.)

YAML allows for maps with arbitrary keys. This is non-trivial to support in Go, and so we provide
special functions (`MapGet`, `MapPut`, `MapDelete`, `MapMerge`) that replace the Go native
functionality with additional support for detecting and handling complex keys. This feature is
provided as an independent library, [yamlkeys](https://github.com/tliron/yamlkeys).

### JSON

JSON can be read into ARD.

However, because JSON has fewer types and more limitations than YAML (no integers, only floats; map
keys can only be strings), ARD will lose some type information when translated into JSON.

We can overcome this challenge by extending JSON with some conventions for encoding extra types.
See [our conventions](json.go).

### XML

XML does not have a type system. Arbitrary XML cannot be parsed into ARD. 

However, with a proper schema and custom reader this could be implemented in the future.

Raw?
----

The data is untreated and not validated. There's no schema.
