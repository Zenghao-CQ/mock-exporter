package collector

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type App struct {
	*kubernetes.Clientset
	namespace string
	time      int64
}

func NewApp(namespace string, time int64) *App {
	var kubeconfig *string
	var config *rest.Config
	var err error
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	if config, err = rest.InClusterConfig(); err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &App{Clientset: clientset, namespace: namespace, time: time}
}

func (app App) GetPodNames(m *Metrics) {
	go func() {
		for {
			m.mutex.Lock()
			pods, err := app.CoreV1().Pods(app.namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
			m.podNames = []string{}
			for _, pod := range pods.Items {
				m.podNames = append(m.podNames, pod.Name)
			}
			if len(m.podNames) == 0 {
				m.podNames = append(m.podNames, "busybox")
			}
			m.mutex.Unlock()
			log.Print("Get pod names:", m.podNames)
			time.Sleep(time.Duration(app.time) * 60 * time.Second)
		}
	}()
}
