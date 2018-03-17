package metrics

// Tag definition
type Tag struct {
	Key, Value string
}

// NewTag constructor
func NewTag(k, v string) Tag {
	return Tag{k, v}
}

// Counter interface
type Counter interface {
	Increase(value int, tags ...Tag)
}

// NullCounter defines a counter that does nothing
type NullCounter struct {
}

// Increase the null counter
func (nc NullCounter) Increase(value int, tags ...Tag) {
}
