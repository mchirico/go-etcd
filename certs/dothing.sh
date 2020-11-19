#!/bin/bash
# dothing put foo stuff
# ./dothing2  --user=root:password put foo stuff
ENDPOINTS='etcd.cwxstat.io:2379'
ETCDCTL_API=3 etcdctl \
	   --endpoints=${ENDPOINTS} \
	   --cacert="./ca.crt" \
	   --cert="./client.crt" \
	   --key="./client.key" \
	   ${@}
