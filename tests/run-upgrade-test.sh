#!/bin/bash
OLD_VERSION=$1
set -eux

if [[ -z "${OLD_VERSION}" ]]; then
  echo "Must provide old althea-l1 version for upgrade test, make sure it matches a version at https://github.com/AltheaFoundation/althea-L1/releases"
  exit 1
fi

# Remove existing container instance
set +e
docker rm -f althea_all_up_test_instance
set -e

# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

set +u
if [[ -z ${NO_IMAGE_BUILD} ]]; then
bash $DIR/build-container.sh
fi
set -u

NODES=4

# setup for Mac apple silicon compatibility
PLATFORM_CMD=""
if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ -n $(sysctl -a | grep brand | grep "Apple") ]]; then
       echo "Setting --platform=linux/amd64 for Mac apple silicon compatibility"
       PLATFORM_CMD="--platform=linux/amd64"; fi
fi

# Run new test container instance
PORTS="-p 9090:9090 -p 26657:26657 -p 1317:1317 -p 8545:8545"
docker run --name althea_all_up_test_instance $PLATFORM_CMD --cap-add=NET_ADMIN $PORTS althea-base /bin/bash /althea/tests/container-scripts/upgrade-test-internal.sh $NODES $OLD_VERSION
