#!/bin/bash
set -o errexit
source shell/docker_img.sh
IMG=$(SetDockerImg $1)
if [[ ${IMG} == error:* ]]; then
    echo "${IMG}"
    exit 1
fi
docker build --force-rm=true -t ${IMG} .
docker push ${IMG}
docker images|grep none|awk '{print $3 }'|xargs docker rmi
docker rmi ${IMG}
