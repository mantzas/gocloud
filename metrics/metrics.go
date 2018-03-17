package metrics

// Tag definition
type Tag struct {
	Key, Value string
}

// NewTag constructor
func NewTag(k, v string) Tag {
	return Tag{k, v}
}
