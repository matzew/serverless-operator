package knativekafka

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	util "knative.dev/operator/pkg/reconciler/common/testing"
	"testing"
)

const (
	StorageVersionMigration = "storage-version-migration"
)

func TestJobTransform(t *testing.T) {
	tests := []struct {
		name    string
		version string
		//		component v1alpha1.KComponent
		job      batchv1.Job
		expected string
	}{{
		name:     "ChangeNameForServingJob",
		job:      createJob(StorageVersionMigration, ""),
		expected: StorageVersionMigration + "-eventing-kafka",
	}, {
		name:     "ChangeNameWithGeneratedNameForServingJob",
		job:      createJob("", StorageVersionMigration),
		expected: StorageVersionMigration + "-eventing-kafka",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unstructuredJob := util.MakeUnstructured(t, &tt.job)
			transform := JobTransform()
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
