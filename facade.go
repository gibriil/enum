package enum

import (
	"iter"
	"reflect"

	"github.com/gibriil/enum/internal"
)

// member is the package-private implementation behind Member. Keeping the
// internal entry named and private prevents internal.Entry from becoming part
// of the exported API through embedding.
type member[T any] struct {
	m internal.Entry[T]
}

// Valid reports whether the member has been initialized and belongs to a
// registered enum definition.
func (m member[T]) Valid() bool { return m.m.Valid() }

// Name returns the registered name of the member.
func (m member[T]) Name() string { return m.m.Name() }

// String returns the registered name of the member.
func (m member[T]) String() string { return m.m.String() }

// Index returns the member's zero-based position in its namespace.
func (m member[T]) Index() int { return m.m.Index() }

// enumEntry and setEnumEntry are used by the package's initialization code so
// it does not need to reflect over private fields.
func (m Member[T]) enumEntry() internal.Entry[T] { return m.member.m }

func (m *Member[T]) setEnumEntry(entry internal.Entry[T]) {
	m.member.m = entry
}

type memberCarrier[T any] interface {
	enumEntry() internal.Entry[T]
	setEnumEntry(internal.Entry[T])
}

// namespace is the package-private implementation behind Namespace. Its
// forwarding methods preserve the public collection API without exposing
// internal.Definition in the public method set.
type namespace[T any] struct {
	ns *internal.Definition[T]
}

// Len returns the number of registered members in the namespace.
func (n namespace[T]) Len() int { return n.ns.Len() }

// ByName returns the member registered under name.
func (n namespace[T]) ByName(name string) (T, bool) { return n.ns.ByName(name) }

// ByIndex returns the member at index in namespace declaration order.
func (n namespace[T]) ByIndex(index int) (T, bool) { return n.ns.ByIndex(index) }

// Values returns a defensive copy of the registered member values.
func (n namespace[T]) Values() []T { return n.ns.Values() }

// Names returns a defensive copy of the registered member names.
func (n namespace[T]) Names() []string { return n.ns.Names() }

// All returns an iterator over registered member values in declaration order.
func (n namespace[T]) All() iter.Seq[T] { return n.ns.All() }

// Entries returns an iterator over registered member names and values in
// declaration order.
func (n namespace[T]) Entries() iter.Seq2[string, T] { return n.ns.Entries() }

// EntryType returns the type of values stored by the namespace.
func (n namespace[T]) EntryType() reflect.Type { return n.ns.EntryType() }

// Type returns the type used to identify the namespace definition.
func (n namespace[T]) Type() reflect.Type { return n.ns.Type() }

// Name returns the registered name of the namespace definition.
func (n namespace[T]) Name() string { return n.ns.Name() }
