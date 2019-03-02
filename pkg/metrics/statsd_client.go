package metrics

type statsdClient interface {
	Incr(name string, tags []string, rate float64) error
}
