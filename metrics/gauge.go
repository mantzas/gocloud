package metrics

// Gauge interface
type Gauge interface {
	Set(value float64, tags ...Tag)
}

// NullGauge defines a gauge that does nothing
type NullGauge struct {
}

// Set the null gauge
func (nc NullGauge) Set(value float64, tags ...Tag) {
}
