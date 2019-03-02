package metrics

type Agent interface {
	IncrementCounter(Metric, float64) error
}
