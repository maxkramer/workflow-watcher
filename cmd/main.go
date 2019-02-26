package main

import (
	"flag"
	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/project-interstellar/workflow-watcher/internal"
	log "github.com/sirupsen/logrus"
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
	informer.AddEventHandler(internal.WorkflowEventHandler{})

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
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func loadKubernetesConfiguration() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic("Unable to create cluster configuration " + err.Error())
	}
	return config
}
