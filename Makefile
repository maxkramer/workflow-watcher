generate_mocks:
	mockgen -destination=mocks/mock_queue.go -package=mocks github.com/project-interstellar/workflow-watcher/pkg/queue Queue
	mockgen -destination=mocks/mock_agent.go -package=mocks github.com/project-interstellar/workflow-watcher/pkg/metrics Agent
