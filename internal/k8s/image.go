package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/030/dip/internal/gitactions"
	"github.com/030/dip/internal/slack"
	"github.com/030/dip/pkg/dockerhub"
	"github.com/mitchellh/go-homedir"
	imagestreamv1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	log "github.com/sirupsen/logrus"
	v1k8s "k8s.io/api/apps/v1"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clusterImages []string
)

type Images struct {
	SlackChannelID, SlackToken string
	ToBeValidated              map[string]interface{}
	gitactions.Elements
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

func (i *Images) checkIfOutdatedAndSendMessageToSlackIfTrue(image, kind, name, namespace, regexToFindNewestTag, containerImageWithoutTag, containerImageTagInsideCluster string) error {
	latestTag, err := dockerhub.LatestTagBasedOnRegex(regexToFindNewestTag, containerImageWithoutTag)
	if err != nil {
		return err
	}

	log.Infof("%s %s", latestTag, containerImageTagInsideCluster)
	if latestTag != containerImageTagInsideCluster {
		if err := i.CreateMR(image, latestTag); err != nil {
			return err
		}

		msg := fmt.Sprintf("Image: '%s' with tag: '%s' in %s: '%s' in namespace: '%s' outdated. Latest tag: '%s'", image, containerImageTagInsideCluster, kind, name, namespace, latestTag)
		if err := slack.SendMessage(i.SlackChannelID, msg, i.SlackToken); err != nil {
			return err
		}

		// Ensure that only one outdated image is shown in Slack. This prevents
		// that all images have to be updated right away
		os.Exit(0)
	}

	return nil
}

func (i *Images) checkWhetherImagesAreUpToDate(containerImage, namespace, kind, name string) error {
	clusterImages = append(clusterImages, containerImage)

	imageNamesAndRegexToFindNewestTags := i.ToBeValidated
	for imageName, gitURLAndDockerTagRegexInterface := range imageNamesAndRegexToFindNewestTags {
		gitURLAndDockerTagRegexMap := gitURLAndDockerTagRegexInterface.(map[string]interface{})
		fmt.Println(gitURLAndDockerTagRegexMap)
		gitFile := gitURLAndDockerTagRegexMap["gitfile"].(string)
		gitProjectID := gitURLAndDockerTagRegexMap["gitprojectid"].(string)
		gitURL := gitURLAndDockerTagRegexMap["giturl"].(string)
		regexToFindLatestImageTag := gitURLAndDockerTagRegexMap["regextofindlatestimagetag"].(string)
		reviewer := gitURLAndDockerTagRegexMap["reviewer"].(string)
		tag := gitURLAndDockerTagRegexMap["tag"].(string)
		targetBranch := gitURLAndDockerTagRegexMap["targetbranch"].(string)

		log.Infof("regexToFindNewestTag: '%s'", tag)
		r := regexp.MustCompile("^(" + imageName + "):(" + tag + ")")
		if !r.MatchString(containerImage) {
			log.Info("no match")
			continue
		}
		group := r.FindStringSubmatch(containerImage)
		if len(group) == 0 {
			return fmt.Errorf("containerImage should not be empty")
		}

		containerImageWithoutTag := group[1]
		containerImageTagInsideCluster := group[2]
		log.Infof("containerImageWithoutTag: '%s', containerImageTagInsideCluster: '%s'", containerImageWithoutTag, containerImageTagInsideCluster)

		i.Elements.GitFile = gitFile
		i.Elements.GitProjectID = gitProjectID
		i.Elements.GitURL = gitURL
		i.Elements.RegexToFindLatestImageTag = regexToFindLatestImageTag
		i.Elements.Reviewer = reviewer
		i.Elements.Tag = tag
		i.Elements.TargetBranch = targetBranch

		return i.checkIfOutdatedAndSendMessageToSlackIfTrue(imageName, kind, name, namespace, tag, containerImageWithoutTag, containerImageTagInsideCluster)
	}
	return nil
}

func (i *Images) cronJobInitContainersAndContainers(cronJob v1beta1.CronJob, namespaceName string) error {
	kind := "CronJob"
	name := cronJob.Name
	log.Infof("%s: '%s'", kind, name)

	for _, initContainer := range cronJob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
		initContainerImage := initContainer.Image
		log.Infof("initContainer image: %s", initContainer.Image)
		if err := i.checkWhetherImagesAreUpToDate(initContainerImage, namespaceName, kind, name); err != nil {
			return err
		}
	}
	for _, container := range cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers {
		containerImage := container.Image
		log.Infof("container image: %s", container.Image)
		if err := i.checkWhetherImagesAreUpToDate(containerImage, namespaceName, kind, name); err != nil {
			return err
		}
	}
	return nil
}

func (i *Images) cronJobImages(kcs *kubernetes.Clientset, namespaceName string) error {
	cronJobList, err := kcs.BatchV1beta1().CronJobs(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	cronJobs := cronJobList.Items
	log.Infof("There are %d deployments in the cluster\n", len(cronJobs))
	for _, cronJob := range cronJobs {
		if err := i.cronJobInitContainersAndContainers(cronJob, namespaceName); err != nil {
			return err
		}
	}
	return nil
}

func (i *Images) deploymentInitContainersAndContainers(deployment v1k8s.Deployment, namespaceName string) error {
	kind := "Deployment"
	name := deployment.Name
	log.Infof("%s: '%s'", kind, name)

	for _, initContainer := range deployment.Spec.Template.Spec.InitContainers {
		initContainerImage := initContainer.Image
		log.Infof("initContainer image: %s", initContainer.Image)
		if err := i.checkWhetherImagesAreUpToDate(initContainerImage, namespaceName, kind, name); err != nil {
			return err
		}
	}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		containerImage := container.Image
		log.Infof("container image: %s", container.Image)
		if err := i.checkWhetherImagesAreUpToDate(containerImage, namespaceName, kind, name); err != nil {
			return err
		}
	}
	return nil
}

func (i *Images) deploymentImages(kcs *kubernetes.Clientset, namespaceName string) error {
	deploymentList, err := kcs.AppsV1().Deployments(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	deployments := deploymentList.Items
	log.Infof("There are %d deployments in the cluster\n", len(deployments))
	for _, deployment := range deployments {
		if err := i.deploymentInitContainersAndContainers(deployment, namespaceName); err != nil {
			return err
		}
	}
	return nil
}

func (i *Images) statefulSetInitContainersAndContainers(statefulSet v1k8s.StatefulSet, namespaceName string) error {
	kind := "StatefulSet"
	name := statefulSet.Name
	log.Infof("%s: '%s'", kind, name)

	for _, initContainer := range statefulSet.Spec.Template.Spec.InitContainers {
		initContainerImage := initContainer.Image
		log.Infof("initContainer image: %s", initContainer.Image)
		if err := i.checkWhetherImagesAreUpToDate(initContainerImage, namespaceName, kind, name); err != nil {
			return err
		}
	}
	for _, container := range statefulSet.Spec.Template.Spec.Containers {
		containerImage := container.Image
		log.Infof("container image: %s", container.Image)
		if err := i.checkWhetherImagesAreUpToDate(containerImage, namespaceName, kind, name); err != nil {
			return err
		}
	}
	return nil
}

func (i *Images) statefulSetImages(kcs *kubernetes.Clientset, namespaceName string) error {
	statefulSetList, err := kcs.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	statefulSets := statefulSetList.Items
	log.Infof("There are %d statefulSets in the cluster\n", len(statefulSets))
	for _, statefulSet := range statefulSets {
		if err := i.statefulSetInitContainersAndContainers(statefulSet, namespaceName); err != nil {
			return err
		}
	}
	return nil
}

func (i *Images) imageStreamImages(kcs *rest.Config, namespaceName string) error {
	imagestreamV1Client, err := imagestreamv1.NewForConfig(kcs)
	if err != nil {
		return err
	}

	imageStreamsList, err := imagestreamV1Client.ImageStreams(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	imageStreams := imageStreamsList.Items
	log.Infof("There are %d imageStreams in the cluster\n", len(imageStreams))
	for _, imageStream := range imageStreams {
		fmt.Println()
		for _, tag := range imageStream.Spec.Tags {
			name := tag.Name
			containerImage := tag.From.Name
			if err := i.checkWhetherImagesAreUpToDate(containerImage, namespaceName, "imageStream", name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Images) UpToDate() error {
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
		if err := i.cronJobImages(clientset, namespaceName); err != nil {
			return err
		}
		if err := i.deploymentImages(clientset, namespaceName); err != nil {
			return err
		}
		if err := i.statefulSetImages(clientset, namespaceName); err != nil {
			return err
		}
		if false {
			if err := i.imageStreamImages(auth, namespaceName); err != nil {
				return err
			}
		}
	}

	clusterImages = sort.StringSlice(clusterImages)
	for _, clusterImage := range clusterImages {
		fmt.Println(clusterImage)
	}
	fmt.Println("Number of images:", len(clusterImages))

	return nil
}
