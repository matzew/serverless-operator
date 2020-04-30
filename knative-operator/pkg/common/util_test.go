package common_test

import (
	"testing"

	"os"
	"reflect"

	"github.com/openshift-knative/serverless-operator/knative-operator/pkg/common"

	servingv1alpha1 "knative.dev/serving-operator/pkg/apis/serving/v1alpha1"
)

func TestBuildImageOverrideMapFromEnviron(t *testing.T) {
	cases := []struct {
		name     string
		envVar   string
		value    string
		expected map[string]string
	}{
		{
			name:   "Simple container name",
			envVar: "IMAGE_foo",
			value:  "quay.io/myimage",
			expected: map[string]string{
				"foo": "quay.io/myimage",
			},
		},
		{
			name:   "Deployment+container name",
			envVar: "IMAGE_foo_bar",
			value:  "quay.io/myimage",
			expected: map[string]string{
				"foo/bar": "quay.io/myimage",
			},
		},
		{
			name:   "3 underscores",
			envVar: "IMAGE_foo_bar_baz",
			value:  "quay.io/myimage",
			expected: map[string]string{
				"foo/bar_baz": "quay.io/myimage",
			},
		},
		{
			name:     "Different prefix",
			envVar:   "X_foo",
			value:    "quay.io/myimage",
			expected: map[string]string{},
		},
		{
			name:     "No env var value",
			envVar:   "IMAGE_foo",
			value:    "",
			expected: map[string]string{},
		},
	}

	for i := range cases {
		tc := cases[i]
		os.Setenv(tc.envVar, tc.value)
		overrideMap := common.BuildImageOverrideMapFromEnviron()
		os.Unsetenv(tc.envVar)

		if !reflect.DeepEqual(overrideMap, tc.expected) {
			t.Errorf("Image override map is not equal. Case name: %q. Expected: %v, actual: %v", tc.name, tc.expected, overrideMap)
		}

	}
}

func verifyImageOverride(t *testing.T, registry *servingv1alpha1.Registry, imageName string, expected string) {
	if registry.Override[imageName] != expected {
		t.Errorf("Missing queue image. Expected a map with following override in it : %v=%v, actual: %v", imageName, expected, registry.Override)
	}
}
