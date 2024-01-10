##
# Example to define private data configuration
##

# you need your quay.io user to push the image to your
# repository.
# - Login as QUAY_USER at quay.io
# - You need to create a repository (QUAY_REPOSITORY) in quay.io
# - You need to create a robot account (<user>+<robot>) for the repository (QUAY_LOGIN).
# - You need to assign write permissions to the robot account.
# - Copy the token (QUAY_TOKEN) for the robot account, and fill the fields below.
#
# The link below is to regenerate the token (substitute <user> by yours):
#   https://quay.io/repository/<USER>/<REPOSITORY>?tab=settings

# QUAY_USER is used by build_deploy.sh script which is called by the jenkins job automation
# TODO Set your QUAY_USER
export QUAY_USER := myusername+robotaccount
# QUAY_TOKEN is used by build_deploy.sh script which is called by the jenkins job automation
# TODO Set your QUAY_TOKEN
export QUAY_TOKEN := 238AKJNF8234ANSJ
# QUAY_LOGIN is used to compose the default CONTAINER_IMAGE_BASE
export QUAY_LOGIN := $(firstword $(subst +, ,$(QUAY_USER)))
# QUAY_REPOSITORY is used to compose the base image name when deploying into ephemeral
# TODO Set your QUAY_REPOSITORY ; you have to create it and grant
#      write permissions to the above robot account
export QUAY_REPOSITORY := 7d

# TODO Update CONTAINER_IMAGE_BASE accoddingly to point out to your repository
# This should point out to your repository
CONTAINER_IMAGE_BASE ?= quay.io/$(QUAY_LOGIN)/$(QUAY_REPOSITORY)

# TODO Fill the values:
# You can get a token for you from:
#   https://access.redhat.com/RegistryAuthentication
#   https://access.redhat.com/RegistryAuthentication#creating-registry-service-accounts-6
# To retrieve your token or regenerate it, go to:
#   https://access.redhat.com/terms-based-registry/#/token/YOUR_USERNAME
# Use the user and token from the "Docker login" tab (e.g. `12345|nickname`).
# If the Dockerfile is using the authenticated repositories, you will need this values
# The build_deploy.sh script which is called by jenkins job automation, log in the
# redhat registry. They are specified here to use the same script that use
# the jenkins job and detect deviations at some point into the day to day.
export RH_REGISTRY_USER :=
export RH_REGISTRY_TOKEN :=

# Ephemeral pool to be used when deploying in ephemeral environment
# You can list them by 'bonfire pool list' command
POOL ?= default
# POOL ?= real-managed-kafka

# Set token expiration time, expressed in seconds (default 2 hours)
APP_TOKEN_EXPIRATION_SECONDS ?= 7200

########### NO SECRETS BUT GENERAL OVERRIDES

# Uncomment the next line if you use to need more than 1hour for your
# ephemeral namespace reservation, and update to your wished value.
# Remember that this value will be used for the initial reservation and for
# extending the reservation.
#
# DURATION ?= 4h

# Cluster to use for development purpose
CLUSTER ?= crc-eph.r9lp.p1.openshiftapps.com
