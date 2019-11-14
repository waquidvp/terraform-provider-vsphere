#!/bin/bash
set -e -u -o pipefail

$GOPATH/src/github.com/terraform-providers/terraform-provider-vsphere/scripts/esxi_restore_snapshot.sh $VSPHERE_ESXI_SNAPSHOT
sleep 30
until curl -k "https://$VSPHERE_SERVER/rest/vcenter/datacenter" -i -m 10 2> /dev/null | grep "401 Unauthorized" &> /dev/null; do
sleep 30
done
