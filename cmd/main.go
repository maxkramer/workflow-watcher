package main

import (
	"context"
	"flag"
	"github.com/DataDog/datadog-go/statsd"
	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/project-interstellar/workflow-watcher/internal"
	"github.com/project-interstellar/workflow-watcher/pkg"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	kubeconfig      = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	resourceVersion = flag.String("resourceVersion", "", "(optional) the resource version to begin listening from")
	logLevel        = flag.String("logLevel", "debug", "(optional) log level")
	statsdAddress   = flag.String("statsd-address", "0.0.0.0:8125", "(optional) statsd address")
	log             = logrus.New()
)

func configureStatsdClient() *statsd.Client {
	statsd, err := statsd.New(*statsdAddress)
	if err != nil {
		panic("Failed to create statsd client " + err.Error())
	}

	err = statsd.Count("start", 1, nil, 1)
	if err != nil {
		log.Error("Failed to increment statsd counter `start`")
	}

	statsd.Namespace = "com.maxkramer.workflow-watcher"
	return statsd
}

func main() {
	configureLogger()
	flag.Parse()

	statsd := configureStatsdClient()
	config := loadKubernetesConfiguration()

	wfClient := wfclientset.NewForConfigOrDie(config)
	informer := externalversions.NewSharedInformerFactoryWithOptions(wfClient, 10*time.Minute,
		externalversions.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.ResourceVersion = *resourceVersion
		})).Argoproj().V1alpha1().Workflows().Informer()

	pubsub := pkg.PubSub{Log: log, MessageFactory: pkg.WorkflowChangedMessageFactory{}, Ctx: context.Background(), ProjectId: "", TopicName: ""}
	informer.AddEventHandler(internal.WorkflowEventHandler{Log: log, Queue: pubsub})

	stopper := make(chan struct{})
	configureGracefulExit(stopper, statsd)
	informer.Run(stopper)
}

func configureGracefulExit(stopper chan struct{}, statsd *statsd.Client) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		close(stopper)

		log.Warnf("caught sig: %+v", sig)
		log.Debugf("Flushing statsd client")

		statsdErr := statsd.Count("exit", 1, nil, 1)
		if statsdErr != nil {
			log.Error("Failed to increment statsd exit counter")
		}

		statsdErr = statsd.Flush()
		if statsdErr != nil {
			log.Error("Failed to flush statsd client")
		}

		statsdErr = statsd.Close()
		if statsdErr != nil {
			log.Error("Failed to close statsd client")
		}

		log.Debugf("Wait for 2 second to finish processing")

		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}

func configureLogger() {
	log.Formatter = &logrus.JSONFormatter{}
	log.Out = os.Stdout
	level, err := logrus.ParseLevel(*logLevel)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}
}

func loadKubernetesConfiguration() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic("Unable to create cluster configuration " + err.Error())
	}
	return config
}
