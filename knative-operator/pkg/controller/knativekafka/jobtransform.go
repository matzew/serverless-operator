package knativekafka

import (
	"fmt"
	mf "github.com/manifestival/manifestival"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

func JobTransform() mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetKind() == "Job" {
			job := &batchv1.Job{}
			if err := scheme.Scheme.Convert(u, job, nil); err != nil {
				return err
			}

			component := "eventing-kafka"
			if job.GetName() == "" {
				job.SetName(fmt.Sprintf("%s%s", job.GetGenerateName(), component))
			} else {
				job.SetName(fmt.Sprintf("%s-%s", job.GetName(), component))
			}

			return scheme.Scheme.Convert(job, u, nil)
		}

		return nil
	}
}
