/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// IngressServiceTransform pins the namespace to istio-system for the service named knative-local-gateway.
func IngressServiceTransform() mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetAPIVersion() == "v1" && u.GetKind() == "Service" && u.GetName() == "knative-local-gateway" {
			u.SetNamespace("istio-system")
		}
		return nil
	}
}
