// Package pragma is inspired by Go protobuffers pragmas: it provides embeddable types that
// force early feedback from the compiler.
package pragma

import "sync"

// NotCopiable can be embedded in a struct to force a warning in case of value copy via assignment.
// Declaring this as a 0-length array is guaranteed by the Go spec to have size 0.
type NotCopiable [0]sync.Mutex
