#!/bin/sh

#
# - update deployments running in our local namespace using our image
# - the fully qualified docker image identifier is in artifacts/LATEST
# - the current namespace (e.g from where the CI proxy is running) is set as
#   $NAMESPACE and automatically picked up by the update-images tool
# - the executing proxy will also set $STAMP with some SRE friendly hints (which is
#   also automatically picked up)
# - any container using the image will be updated with the new tag (e.g it will be
#   re-deployed)
#
if [ ! -f artifacts/LATEST ]; then
    echo "unable to find artifacts/LATEST, aborting (job configuration error ?"
    exit 1
fi

/usr/local/bin/update-images $(cat artifacts/LATEST)
if [ $? -ne 0 ]; then
    echo "unable to apply the new image, aborting"
    exit 1
fi