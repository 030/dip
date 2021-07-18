package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/030/dip/pkg/dockerhub"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func inOrOutsideCluster(kubeconfig string) (*rest.Config, error) {
	var config *rest.Config
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Info("~/.kube/config does not exist. Assuming that the program is run inside a cluster.")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("~/.kube/config exists. Assuming that the program is run outside a cluster.")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func authenticate() (*rest.Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := inOrOutsideCluster(kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func cronjobImages(kcs *kubernetes.Clientset, namespaceName string) error {
	cj, err := kcs.BatchV1beta1().CronJobs(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("There are %d cronjobs in the cluster\n", len(cj.Items))
	for _, hi := range cj.Items {
		fmt.Println(hi.Name)
		for _, c := range hi.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
			fmt.Println("init:", c.Image)
		}
		for _, c := range hi.Spec.JobTemplate.Spec.Template.Spec.Containers {
			fmt.Println("containers:", c.Image)
		}
		fmt.Println("----------")
	}
	return nil
}

func podImages(kcs *kubernetes.Clientset, namespaceName string) error {
	pods, err := kcs.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, hi := range pods.Items {
		fmt.Println(hi.Name)
		for _, c := range hi.Spec.InitContainers {
			fmt.Println("init:", c.Image)
		}
		for _, container := range hi.Spec.Containers {
			containerImage := container.Image
			fmt.Println("containerImage:", containerImage)

			r := regexp.MustCompile("^(splunk.*):(.*)")
			if !r.MatchString(containerImage) {
				return fmt.Errorf("no match")
			}
			group := r.FindStringSubmatch(containerImage)
			fmt.Println("HELLO", group[2])

			containerImageWithoutTag := group[1]

			// key lookup in ~/.dip/config to find image and regex for validation
			latestTag, err := dockerhub.LatestTagBasedOnRegex(false, "7.*", containerImageWithoutTag)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(latestTag)
		}
		fmt.Println("----------")
	}
	return nil
}

func Images() error {
	auth, err := authenticate()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(auth)
	if err != nil {
		return err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, namespace := range namespaces.Items {
		namespaceName := namespace.Name
		fmt.Println("========================")
		fmt.Println(namespaceName)
		fmt.Println("========================")
		cronjobImages(clientset, namespaceName)
		podImages(clientset, namespaceName)
	}
	return nil
}
