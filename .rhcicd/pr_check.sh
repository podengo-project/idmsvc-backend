#!/bin/bash

# --------------------------------------------
# Options that must be configured by app owner
# --------------------------------------------
APP_NAME="idmsvc"  # name of app-sre "application" folder this component lives in
COMPONENT_NAME="idmsvc-backend"  # name of app-sre "resourceTemplate" in deploy.yaml for this component
IMAGE="quay.io/cloudservices/idmsvc-backend"  # image location on quay
DOCKERFILE="build/package/Dockerfile"
IQE_PLUGINS="idm"  # name of the IQE plugin for this app.
IQE_MARKER_EXPRESSION="api"  # This is the value passed to pytest -m
IQE_FILTER_EXPRESSION=""  # This is the value passed to pytest -k
IQE_CJI_TIMEOUT="30m"  # This is the time to wait for smoke test to complete or fail
IQE_ENV="ephemeral"

# Install bonfire repo/initialize
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/bootstrap.sh
# This script automates the install / config of bonfire
CICD_URL=https://raw.githubusercontent.com/RedHatInsights/cicd-tools/main
curl -s $CICD_URL/bootstrap.sh > .cicd_bootstrap.sh && source .cicd_bootstrap.sh

# The contents of build.sh can be found at:
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/build.sh
# This script is used to build the image that is used in the PR Check
source $CICD_ROOT/build.sh

# The contents of this script can be found at:
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/deploy_ephemeral_env.sh
# This script is used to deploy the ephemeral environment for smoke tests.
# The manual steps for this can be found in:
# https://internal.cloud.redhat.com/docs/devprod/ephemeral/02-deploying/
source $CICD_ROOT/deploy_ephemeral_env.sh

# Run smoke tests using a ClowdJobInvocation (preferred)
# The contents of this script can be found at:
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/cji_smoke_test.sh
#source $CICD_ROOT/cji_smoke_test.sh

# Post a comment with test run IDs to the PR
# The contents of this script can be found at:
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/post_test_results.sh
#source $CICD_ROOT/post_test_results.sh
mkdir -p $WORKSPACE/artifacts
cat << EOF > $WORKSPACE/artifacts/junit-dummy.xml
<testsuite tests="1">
    <testcase classname="dummy" name="dummytest"/>
</testsuite>
EOF

RESULT=$?

if [[ $RESULT -ne 0 ]]; then
    exit $RESULT
fi

exit $RESULT
