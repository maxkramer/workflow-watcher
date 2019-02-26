# WIP: Argo Workflow Watcher

> [Argo Workflows][1] is an open source container-native workflow engine for orchestrating parallel jobs on Kubernetes. Argo Workflows is implemented as a Kubernetes CRD (Custom Resource Definition).

This watcher uses the Kubernetes informer API to listen for changes on the Workflow CRD and will write these changes to a queue to be processed elsewhere. 

This is necessary because by implementing a non-queued  Kubernetes watcher, you run the risk of losing out on events should an exception take place while attempting to handle a change, and it will be incredibly difficult to scale as part of a bigger application.

## Architecture

The architecture for this solution would be something along the lines of:

![Architecture diagram][2]

New Workflows would be created through an API interface which communicates directly with the Kubernetes API server. 

The worker uses an [Informer][3] / Watch to listen for changes on the Workflow CRD and pushes these changes to a queueing system (SQS, PubSub, Socket etc). 

The API is then able to process these event changes asychronously, simply nack'ing any messages that failed to be processed, to be done at a later point.

This whole approach is an attempt to introduce fault-tolerance on the API level, into an async design pattern.

Checkpointing will be implemented using the `resourceVersion` on Kubernetes API objects to continue from where it left off and to ensure that no events are missed. A simple cache will be used, ideally a k8s config map with a single key will be written to on termination, but Redis or equiv could also be used. Just seems a bit over the top for storing a single key value pair.


## Running the watcher
Simply run `go build -o workflow-watcher ./cmd` to build the binary, which can then be run locally.

The binary supports two flags:
- `kubeconfig`: The path to the kubeconfig file containing your cluster credentials. (Master IP is taken from your current-context).
- `resourceVersion`: The resourceVersion of the object that the informer should begin listening from. 

```
$ ./workflow-watcher -kubeconfig $HOME/.kube/config
```

![example output][4]

## Todo:

- Implement an abstract queueing interface and initial adaptor for GCP PubSub
- Implement an abstract storage interface for managing the resourceVersion (esp between deploys)
- Write tests


[1]: https://github.com/argoproj/argo
[2]: https://i.imgur.com/pfbmvd7.png
[3]: https://medium.com/firehydrant-io/stay-informed-with-kubernetes-informers-4fda2a21da9e
[4]: https://i.imgur.com/04fR74b.png
