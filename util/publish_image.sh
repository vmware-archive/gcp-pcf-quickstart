#!/usr/bin/env bash

#
# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

if [ -z ${BASE_IMAGE_SELFLINK+x} ]; then
    echo "BASE_IMAGE_SELFLINK required"
    exit 1
fi

if [ -z ${DESTINATION_BUCKET+x} ]; then
    echo "DESTINATION_BUCKET required"
    exit 1
fi

if [ -z ${DESTINATION_NAME+x} ]; then
    echo "DESTINATION_NAME required"
    exit 1
fi

instance=baked-opsman
iamge_snapshot=image-snapshot
zone=us-east1-c

gcloud compute instances create --image ${BASE_IMAGE_SELFLINK} ${instance} --zone=${zone} --machine-type=n1-standard-32
gcloud compute instances stop ${instance} --zone=${zone}
gcloud compute disks snapshot ${instance} --snapshot-names ${iamge_snapshot} --zone=${zone}
gcloud compute disks create image-disk --source-snapshot ${iamge_snapshot} --zone=${zone} --type=pd-ssd

startup=$(mktemp)
cat << STARTUP_SCRIPT > ${startup}
sudo mkdir /mnt/tmp
sudo mkfs.ext4 -F /dev/disk/by-id/google-local-ssd-0
sudo mount -o discard,defaults /dev/disk/by-id/google-local-ssd-0 /mnt/tmp
sudo dd if=/dev/disk/by-id/google-image-disk of=/mnt/tmp/disk.raw bs=4096
cd /mnt/tmp
sudo tar czvf myimage.tar.gz disk.raw
gsutil -o GSUtil:parallel_composite_upload_threshold=150M cp /mnt/tmp/myimage.tar.gz gs://${DESTINATION_BUCKET}/${DESTINATION_NAME}.tar.gz
shutdown -h now 'done'
STARTUP_SCRIPT

gcloud compute instances create baked-opsman-capture --scopes storage-rw \
    --disk name=image-disk,device-name=image-disk \
    --image-family ubuntu-1604-lts \
    --image-project ubuntu-os-cloud \
    --local-ssd interface=scsi \
    --zone=${zone} \
    --machine-type=n1-standard-32 \
    --metadata-from-file startup-script=${startup}

gcloud compute instances tail-serial-port-output baked-opsman-capture --zone=${zone}
