package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/030/dip/pkg/dockerhub"
	sasm "github.com/030/sasm/pkg/slack"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const ext = "yml"

var clusterImages []string

func viperBase(path, filename string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(filepath.Join(home, ".dip"))

	if path != "" {
		viper.SetConfigFile(filepath.Join(path, filename+"."+ext))
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
	}
	return nil
}

func inOrOutsideCluster(kubeconfig string) (*rest.Config, error) {
	var config *rest.Config
	_, err := os.Stat(kubeconfig)
	if os.IsNotExist(err) {
		log.Info("~/.kube/config does not exist. Assuming that the program is run inside a cluster.")
		config, err = rest.InClusterConfig()
	} else {
		log.Info("~/.kube/config exists. Assuming that the program is run outside a cluster.")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
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

func checkWhetherImagesAreUpToDate(containerImage, namespace, path, kind, name string) error {
	clusterImages = append(clusterImages, containerImage)

	if err := viperBase(path, "config"); err != nil {
		return err
	}
	images := viper.GetStringMap("dip_images")
	if len(images) == 0 {
		return fmt.Errorf("no images found. Check whether the 'dip_images' variable is populated in '%s'", viper.ConfigFileUsed())
	}
	for image, tag := range images {
		tagString := fmt.Sprintf("%v", tag)
		r := regexp.MustCompile("^(" + image + "):(" + tagString + ")")
		if !r.MatchString(containerImage) {
			log.Info("no match")
			continue
		}
		group := r.FindStringSubmatch(containerImage)

		containerImageTagInsideCluster := group[2]
		containerImageWithoutTag := group[1]

		latestTag, err := dockerhub.LatestTagBasedOnRegex(false, tagString, containerImageWithoutTag)
		if err != nil {
			return err
		}
		if latestTag != tagString {
			msg := fmt.Sprintf("Image: '%s' with tag: '%s' in %s: '%s' in namespace: '%s' outdated. Latest tag: '%s'", image, containerImageTagInsideCluster, kind, name, namespace, latestTag)
			log.Info("Sending message to Slack...")
			t := sasm.Text{Type: "mrkdwn", Text: msg}
			b := []sasm.Blocks{{Type: "section", Text: &t}}
			d := sasm.Data{Blocks: b, Channel: "#dip", Icon: ":dip:", Username: "dip"}

			if err := viperBase(path, "creds"); err != nil {
				return err
			}
			slackToken := viper.GetString("slack_token")
			if slackToken == "" {
				log.Fatalf("slack_token should not be empty. Check whether these resides in: '%s'", viper.ConfigFileUsed())
			}
			if err := d.PostMessage(slackToken); err != nil {
				return err
			}
			log.Warningf("image: '%s' outdated", image)
			log.Info("other clusterImages:")
			for _, clusterImage := range clusterImages {
				log.Info(clusterImage)
			}
			os.Exit(0)
		}
	}
	return nil
}

func cronJobInitContainersAndContainers(cronJob v1beta1.CronJob, path, namespaceName string) error {
	kind := "CronJob"
	name := cronJob.Name
	log.Infof("%s: '%s'", kind, name)

	for _, initContainer := range cronJob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
		initContainerImage := initContainer.Image
		log.Infof("initContainer image: %s", initContainer.Image)
		if err := checkWhetherImagesAreUpToDate(initContainerImage, namespaceName, path, kind, name); err != nil {
			return err
		}
	}
	for _, container := range cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers {
		containerImage := container.Image
		log.Infof("container image: %s", container.Image)
		if err := checkWhetherImagesAreUpToDate(containerImage, namespaceName, path, kind, name); err != nil {
			return err
		}
	}
	return nil
}

func cronJobImages(kcs *kubernetes.Clientset, path, namespaceName string) error {
	cronJobList, err := kcs.BatchV1beta1().CronJobs(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	cronJobs := cronJobList.Items
	log.Infof("There are %d deployments in the cluster\n", len(cronJobs))
	for _, cronJob := range cronJobs {
		if err := cronJobInitContainersAndContainers(cronJob, path, namespaceName); err != nil {
			return err
		}
	}
	return nil
}

func deploymentInitContainersAndContainers(deployment v1.Deployment, path, namespaceName string) error {
	kind := "Deployment"
	name := deployment.Name
	log.Infof("%s: '%s'", kind, name)

	for _, initContainer := range deployment.Spec.Template.Spec.InitContainers {
		initContainerImage := initContainer.Image
		log.Infof("initContainer image: %s", initContainer.Image)
		if err := checkWhetherImagesAreUpToDate(initContainerImage, namespaceName, path, kind, name); err != nil {
			return err
		}
	}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		containerImage := container.Image
		log.Infof("container image: %s", container.Image)
		if err := checkWhetherImagesAreUpToDate(containerImage, namespaceName, path, kind, name); err != nil {
			return err
		}
	}
	return nil
}

func deploymentImages(kcs *kubernetes.Clientset, path, namespaceName string) error {
	deploymentList, err := kcs.AppsV1().Deployments(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	deployments := deploymentList.Items
	log.Infof("There are %d deployments in the cluster\n", len(deployments))
	for _, deployment := range deployments {
		if err := deploymentInitContainersAndContainers(deployment, path, namespaceName); err != nil {
			return err
		}
	}
	return nil
}

func statefulSetInitContainersAndContainers(statefulSet v1.StatefulSet, path, namespaceName string) error {
	kind := "StatefulSet"
	name := statefulSet.Name
	log.Infof("%s: '%s'", kind, name)

	for _, initContainer := range statefulSet.Spec.Template.Spec.InitContainers {
		initContainerImage := initContainer.Image
		log.Infof("initContainer image: %s", initContainer.Image)
		if err := checkWhetherImagesAreUpToDate(initContainerImage, namespaceName, path, kind, name); err != nil {
			return err
		}
	}
	for _, container := range statefulSet.Spec.Template.Spec.Containers {
		containerImage := container.Image
		log.Infof("container image: %s", container.Image)
		if err := checkWhetherImagesAreUpToDate(containerImage, namespaceName, path, kind, name); err != nil {
			return err
		}
	}
	return nil
}

func statefulSetImages(kcs *kubernetes.Clientset, path, namespaceName string) error {
	statefulSetList, err := kcs.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	statefulSets := statefulSetList.Items
	log.Infof("There are %d statefulSets in the cluster\n", len(statefulSets))
	for _, statefulSet := range statefulSets {
		if err := statefulSetInitContainersAndContainers(statefulSet, path, namespaceName); err != nil {
			return err
		}
	}
	return nil
}

func Images(path string) error {
	auth, err := authenticate()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(auth)
	if err != nil {
		return err
	}

	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	namespaces := namespaceList.Items
	for _, namespace := range namespaces {
		namespaceName := namespace.Name
		log.Infof("namespaceName: '%s'", namespaceName)
		if err := cronJobImages(clientset, path, namespaceName); err != nil {
			return err
		}
		if err := deploymentImages(clientset, path, namespaceName); err != nil {
			return err
		}
		if err := statefulSetImages(clientset, path, namespaceName); err != nil {
			return err
		}

	}
	return nil
}
