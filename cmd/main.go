package main

import (
	"context"
	"flag"
	"fmt"
	datadog "github.com/DataDog/datadog-go/statsd"
	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/project-interstellar/workflow-watcher/internal"
	"github.com/project-interstellar/workflow-watcher/pkg"
	"github.com/project-interstellar/workflow-watcher/pkg/queue"
	"github.com/project-interstellar/workflow-watcher/pkg/storage"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	kubeconfig      = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	logLevel        = flag.String("logLevel", "debug", "(optional) log level")
	statsdAddress   = flag.String("statsd-address", "0.0.0.0:8125", "(optional) statsd address")
	redisAddress    = flag.String("redis-address", "0.0.0.0:6379", "(optional) redis address")
	redisPassword   = flag.String("redis-password", "", "(optional) redis password")
	pubsubProjectId = flag.String("pubsub-project-id", "", "(optional) project-id to use for Google PubSub")
	pubsubTopicId   = flag.String("pubsub-topic-id", "", "(optional) topic id to use for Google PubSub")
)

func configureStatsdClient() *datadog.Client {
	statsd, err := datadog.New(*statsdAddress)
	if err != nil {
		panic("Failed to create statsd client " + err.Error())
	}

	statsd.Namespace = "com.maxkramer.workflow-watcher."
	err = statsd.Count("start", 1, nil, 1)
	if err != nil {
		log.Error("Failed to increment statsd counter `start`")
	}

	return statsd
}

func main() {
	configureLogger()
	flag.Parse()

	statsd := configureStatsdClient()
	stopper := make(chan struct{})
	configureGracefulExit(stopper, statsd)

	config := loadKubernetesConfiguration()

	wfClient := wfclientset.NewForConfigOrDie(config)

	redis := storage.NewRedis(*redisAddress, *redisPassword)
	resourceVersion, redisErr := redis.Get("resourceVersion")
	if redisErr != nil {
		log.Error("Failed to fetch resourceVersion from Redis ", redisErr)
		resourceVersion = ""
	}

	log.Infof("Starting informer with resourceVersion \"%s\"", resourceVersion)

	pubsub := queue.NewPubSub(context.Background(), pkg.WorkflowChangedMessageFactory{}, *pubsubProjectId, *pubsubTopicId)
	informer := externalversions.NewSharedInformerFactoryWithOptions(wfClient, 10*time.Minute,
		externalversions.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.ResourceVersion = resourceVersion.(string)
		})).Argoproj().V1alpha1().Workflows().Informer()

	resourceVersionChannel := make(chan string)
	defer close(resourceVersionChannel)

	informer.AddEventHandler(internal.NewWorkflowEventHandler(pubsub, statsd, resourceVersionChannel))
	go listenToResourceVersionUpdates(resourceVersionChannel, redis, statsd)

	informer.Run(stopper)
}

func listenToResourceVersionUpdates(channel chan string, redis *storage.Redis, statsd *datadog.Client) {
	for resourceVersion := range channel {
		log.Debugf("Updating resourceVersion in Redis to %s", resourceVersion)
		err := redis.Set("resourceVersion", resourceVersion)
		if err != nil {
			statsdErr := statsd.Count("redis.resource-version.write-error", 1, nil, 1)
			if statsdErr != nil {
				log.Error("Failed writing redis.resource-version.write-error to statsd ", statsdErr)
			}
			log.Error("Error setting resourceVersion in Redis", err)
		} else {
			statsdErr := statsd.Count("redis.resource-version.write-success", 1, nil, 1)
			if statsdErr != nil {
				log.Error("Failed writing redis.resource-version.write-success to statsd ", statsdErr)
			}
		}
	}

	log.Debug("Channel closed. Stopped listening to resourceVersion changes")
}

func configureGracefulExit(stopper chan struct{}, statsd *datadog.Client) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop

		log.Warnf("caught sig: %+v", sig)
		log.Debugf("Flushing statsd client")

		statsdErr := statsd.Count(fmt.Sprintf("stop.%+v", sig), 1, nil, 1)
		if statsdErr != nil {
			log.Error("Failed to increment statsd exit counter ", statsdErr)
		}

		statsdErr = statsd.Flush()
		if statsdErr != nil {
			log.Error("Failed to flush statsd client ", statsdErr)
		}

		statsdErr = statsd.Close()
		if statsdErr != nil {
			log.Error("Failed to close statsd client ", statsdErr)
		}

		log.Debug("Stopping informer")
		stopper <- struct{}{}

		log.Debug("Beginning 5 second grace-period to finish processing")
		time.Sleep(5 * time.Second)
		log.Debug("Grace-period ended")

		os.Exit(0)
	}()
}

func configureLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	level, err := log.ParseLevel(*logLevel)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func loadKubernetesConfiguration() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic("Unable to create cluster configuration " + err.Error())
	}
	return config
}
