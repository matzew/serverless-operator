package knativekafka

import (
	"fmt"
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

			version := instance.Status.Version
			if job.GetName() == "" {
				job.SetName(fmt.Sprintf("%s%s", job.GetGenerateName(), version))
			} else {
				job.SetName(fmt.Sprintf("%s-%s", job.GetName(), version))
			}

			return scheme.Scheme.Convert(job, u, nil)
		}

		return nil
	}
}
