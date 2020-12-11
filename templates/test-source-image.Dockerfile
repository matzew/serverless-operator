FROM src

COPY oc /usr/bin/oc
COPY --from=registry.svc.ci.openshift.org/openshift/knative-v__SERVING_VERSION__:knative-serving-src /go/src/knative.dev/serving/ /go/src/knative.dev/serving/
COPY --from=registry.svc.ci.openshift.org/openshift/knative-v__EVENTING_VERSION__:knative-eventing-src /go/src/knative.dev/eventing/ /go/src/knative.dev/eventing/
COPY --from=registry.svc.ci.openshift.org/openshift/knative-v__EVENTING_KAFKA_VERSION__:knative-eventing-kafka-src /go/src/knative.dev/eventing-kafka/ /go/src/knative.dev/eventing-kafka/

RUN chmod g+w /go/src/knative.dev/serving/ && chmod g+w /go/src/knative.dev/eventing/ && chmod g+w /go/src/knative.dev/eventing-kafka/
