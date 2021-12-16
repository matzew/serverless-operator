package knativekafka

import (
	"fmt"
	"os"

	mf "github.com/manifestival/manifestival"
	serverlessoperatorv1alpha1 "github.com/openshift-knative/serverless-operator/knative-operator/pkg/apis/operator/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

func JobTransform(instance *serverlessoperatorv1alpha1.KnativeKafka) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetKind() == "Job" {
			job := &batchv1.Job{}
			if err := scheme.Scheme.Convert(u, job, nil); err != nil {
				return err
			}

			version := targetVersion(instance)
			component := "eventing-kafka"
			if job.GetName() == "" {
				job.SetName(fmt.Sprintf("%s%s-%s", job.GetGenerateName(), component, version))
			} else {
				job.SetName(fmt.Sprintf("%s-%s-%s", job.GetName(), component, version))
			}

			return scheme.Scheme.Convert(job, u, nil)
		}

		return nil
	}
}

func targetVersion(instance *serverlessoperatorv1alpha1.KnativeKafka) string {

	if version := instance.Status.Version; version != "" {
		return version
	}
	return os.Getenv("KNATIVE_EVENTING_KAFKA_VERSION")
}
