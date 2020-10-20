#!/usr/bin/env bash

# shellcheck disable=SC1091,SC1090
source "$(dirname "${BASH_SOURCE[0]}")/lib.bash"

set -Eeuo pipefail

# Enable extra verbosity if running in CI.
if [ -n "$OPENSHIFT_CI" ]; then
  env
fi
debugging.setup

# Setup cluster for the "test-e2e" target
if [[ $TEST_KNATIVE_KAFKA == false ]]; then
  scale_up_workers || exit $?
  create_namespaces || exit $?
  create_htpasswd_users && add_roles || exit $?

  failed=0

  (( !failed )) && install_catalogsource || failed=1
  (( !failed )) && logger.success '🚀 Cluster prepared for testing.'

  # Run serverless-operator specific tests.
  (( !failed )) && serverless_operator_e2e_tests || failed=2
  (( !failed )) && ensure_serverless_installed || failed=3
  # Run Knative Serving & Eventing downstream E2E tests.
  (( !failed )) && downstream_serving_e2e_tests || failed=4
  (( !failed )) && downstream_eventing_e2e_tests || failed=5
else
  # Run Kafka things only when expclitly enabled in the "test-e2e-kafka" target
  (( !failed )) && serverless_operator_kafka_e2e_tests || failed=6
fi

(( failed )) && dump_state
(( failed )) && exit $failed

success
