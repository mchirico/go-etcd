#!/bin/bash
etcd --name infra0 --data-dir infra0   --cert-file=./server.crt --key-file=./server.key   --advertise-client-urls=https://0.0.0.0:2379 --listen-client-urls=https://0.0.0.0:2379
