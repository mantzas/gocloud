package metrics

// Tag definition
type Tag struct {
	Key, Value string
}

// NewTag constructor
func NewTag(k, v string) Tag {
	return Tag{k, v}
}

// Metric interface
type Metric interface {
	IncreaseCounter(value int, tags ...Tag)
}
