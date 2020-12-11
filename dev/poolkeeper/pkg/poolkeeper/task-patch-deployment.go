package poolkeeper

import (
	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"
)

// PatchDeploymentTask patches
type PatchDeploymentTask struct {
	// NamespaceSelector specifies which namespace's pods should not be patched
	Namespace string `json:"namespace"`

	// Patch is the string to apply per types.MergePatchType
	Patch string `json:"patch"`
}

func (pa *PatchDeploymentTask) run(clientset *kubernetes.Clientset) {
	podList, err := clientset.CoreV1().Pods(pa.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("unable to list all pods for PatchDeploymentTask: ", err)
		return
	}
	allPods := podList.Items

	deploymentsToPatch := make(map[string]bool)
	for _, p := range allPods {
		nodeName := p.Spec.NodeName
		if nodeName == "" {
			continue
		}
		if p.ObjectMeta.Namespace != pa.Namespace {
			continue
		}

		for _, or := range p.ObjectMeta.OwnerReferences {
			if or.Kind == "Deployment" { // Where TF is this constant??
				deploymentsToPatch[or.Name] = true
				break
			}
		}
	}

	appsv1 := clientset.AppsV1()
	for deployment := range deploymentsToPatch {
		_, err := appsv1.Deployments(pa.Namespace).Patch(deployment, types.MergePatchType, []byte(pa.Patch))
		if err != nil {
			log.WithField("deployment", deployment).WithError(err).Error("error patching deployment")
			continue
		}
	}
}
