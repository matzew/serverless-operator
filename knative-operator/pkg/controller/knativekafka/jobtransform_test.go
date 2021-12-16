package knativekafka

import (
	"testing"

	serverlessoperatorv1alpha1 "github.com/openshift-knative/serverless-operator/knative-operator/pkg/apis/operator/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	util "knative.dev/operator/pkg/reconciler/common/testing"
)

const (
	StorageVersionMigration = "storage-version-migration"
)

func TestJobTransform(t *testing.T) {
	tests := []struct {
		name      string
		component serverlessoperatorv1alpha1.KnativeKafka
		job       batchv1.Job
		expected  string
	}{{
		name: "ChangeNameForKnativeKafka",
		component: serverlessoperatorv1alpha1.KnativeKafka{
			Status: serverlessoperatorv1alpha1.KnativeKafkaStatus{Version: "1.0.0"},
		},
		job:      createJob(StorageVersionMigration, ""),
		expected: StorageVersionMigration + "-eventing-kafka-1.0.0",
	}, {
		name: "ChangeNameWithGeneratedNameForKnativeKafka",
		component: serverlessoperatorv1alpha1.KnativeKafka{
			Status: serverlessoperatorv1alpha1.KnativeKafkaStatus{Version: "1.0.0"},
		},
		job:      createJob("", StorageVersionMigration),
		expected: StorageVersionMigration + "-eventing-kafka-1.0.0",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unstructuredJob := util.MakeUnstructured(t, &tt.job)
			transform := JobTransform(&tt.component)
			transform(&unstructuredJob)

			var job = &batchv1.Job{}
			err := scheme.Scheme.Convert(&unstructuredJob, job, nil)
			util.AssertEqual(t, err, nil)
			util.AssertDeepEqual(t, job.Name, tt.expected)
		})
	}
}

func createJob(name, gen string) batchv1.Job {
	return batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: gen + "-",
		},
	}
}
