#!/bin/bash

# --------------------------------------------
# Options that must be configured by app owner
# --------------------------------------------
# TODO Define the APP_NAME to be used
APP_NAME="idmsvc"  # name of app-sre "application" folder this component lives in
COMPONENT_NAME="${APP_NAME}-backend"  # name of resourceTemplate component for deploy
IMAGE="quay.io/cloudservices/${COMPONENT_NAME}"  # image location on quay
CONTAINERFILE="build/Dockerfile"

# TODO Define IQE_PLUGINS value
IQE_PLUGINS="${APP_NAME}"  # name of the IQE plugin for this app.
IQE_MARKER_EXPRESSION="api"  # This is the value passed to pytest -m
IQE_FILTER_EXPRESSION=""  # This is the value passed to pytest -k
IQE_CJI_TIMEOUT="30m"  # This is the time to wait for smoke test to complete or fail


# Install bonfire repo/initialize
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/bootstrap.sh
# This script automates the install / config of bonfire
CICD_URL=https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd
curl -s $CICD_URL/bootstrap.sh > .cicd_bootstrap.sh && source .cicd_bootstrap.sh

# This script is used to build the image that is used in the PR Check
source $CICD_ROOT/build.sh

# This script is used to deploy the ephemeral environment for smoke tests.
# The manual steps for this can be found in:
# https://internal.cloud.redhat.com/docs/devprod/ephemeral/02-deploying/
source $CICD_ROOT/deploy_ephemeral_env.sh

# Run smoke tests using a ClowdJobInvocation and iqe-tests
source $CICD_ROOT/cji_smoke_test.sh

# Post a comment with test run IDs to the PR
source $CICD_ROOT/post_test_results.sh
