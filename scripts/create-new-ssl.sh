#!/bin/bash

mkdir -p ssl

hostname=$1
if [ "" == "$hostname" ]; then
	hostname="nyantan.net"
fi
company=$2
if [ "" == "$company" ]; then
	company="Nyantan"
fi

openssl req \
	-new \
	-newkey rsa:2048 \
	-nodes \
	-keyout ssl/$hostname.pem \
	-out 	ssl/$hostname.csr \
	-days 	365 \
	-subj 	"/C=HU/ST=Budapest/L=Budapest/O=$company/OU=None/CN=$hostname"

echo "The new ssl csr is:"
echo ""
cat ssl/$hostname.csr
