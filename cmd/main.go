package main

import (
	"context"
	"flag"
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
	log             = logrus.New()
)

func main() {
	configureLogger()
	flag.Parse()
	config := loadKubernetesConfiguration()

	wfClient := wfclientset.NewForConfigOrDie(config)
	informer := externalversions.NewSharedInformerFactoryWithOptions(wfClient, 10*time.Minute,
		externalversions.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.ResourceVersion = *resourceVersion
		})).Argoproj().V1alpha1().Workflows().Informer()

	pubsub := pkg.PubSub{Log: log, MessageFactory: pkg.WorkflowChangedMessageFactory{}, Ctx: context.Background(), ProjectId: "", TopicName: ""}
	informer.AddEventHandler(internal.WorkflowEventHandler{Log: log, Queue: pubsub})

	stopper := make(chan struct{})
	configureGracefulExit(stopper)
	informer.Run(stopper)
}

func configureGracefulExit(stopper chan struct{}) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		close(stopper)

		log.Warnf("caught sig: %+v", sig)
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
