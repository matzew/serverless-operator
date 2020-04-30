package common

import (
	"fmt"
	servingv1alpha1 "knative.dev/serving-operator/pkg/apis/serving/v1alpha1"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
)

var Log = logf.Log.WithName("knative").WithName("openshift")

// Configure is a  helper to set a value for a key, potentially overriding existing contents.
func Configure(ks *servingv1alpha1.KnativeServing, cm, key, value string) bool {
	if ks.Spec.Config == nil {
		ks.Spec.Config = map[string]map[string]string{}
	}

	old, found := ks.Spec.Config[cm][key]
	if found && value == old {
		return false
	}

	if ks.Spec.Config[cm] == nil {
		ks.Spec.Config[cm] = map[string]string{}
	}

	ks.Spec.Config[cm][key] = value
	Log.Info("Configured", "map", cm, key, value, "old value", old)
	return true
}

// IngressNamespace returns namespace where ingress is deployed.
func IngressNamespace(servingNamespace string) string {
	return servingNamespace + "-ingress"
}

// BuildImageOverrideMapFromEnviron creates a map to overrides registry images
func BuildImageOverrideMapFromEnviron() map[string]string {
	overrideMap := map[string]string{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "IMAGE_") {
			suffix := strings.SplitN(pair[0], "_", 2)[1]

			var name string

			// convert
			// "IMAGE_foo=docker.io/foo"
			// "IMAGE_eventing-controller_eventing-controller=docker.io/foo2"
			// to
			// fo0: docker.io/foo
			// eventing-controller/eventing-controller: docker.io/foo2
			suffixParts := strings.SplitN(suffix, "_", 2)
			if len(suffixParts) == 2 {
				name = fmt.Sprintf("%s/%s", suffixParts[0], suffixParts[1])
			} else {
				name = fmt.Sprintf("%s", suffix)
			}

			if pair[1] != "" {
				overrideMap[name] = pair[1]
			}
		}
	}
	return overrideMap
}
