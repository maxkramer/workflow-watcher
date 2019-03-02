package metrics

type Metric string;

var (
	QueueFailure         Metric = "queue.failure"
	QueueSuccess         Metric = "queue.success"
	EventHandlerAdded    Metric = "event-handler.added"
	EventHandlerModified Metric = "event-handler.modified"
	EventHandlerDeleted  Metric = "event-handler.deleted"
)
