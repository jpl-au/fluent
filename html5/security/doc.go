// Package security holds the cross-package attribute-safety probe tests: proof
// that hostile runtime strings are neutralised end-to-end by the generated
// setters (escaping) and URL sinks (scheme filtering), and that the raw hatch
// deliberately does not escape. The behaviour is generated from spec YAML by
// fluent-generator; these tests pin it against regressions.
package security
