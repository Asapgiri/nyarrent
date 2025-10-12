#!/bin/bash

domain=$1

private_key=$2
cert=$3
ca_bundle=$4

if [ -z "${domain}" ] || [ -z "${private_key}" ] || [ -z "${cert}" ] || [ -z "${ca_bundle}" ]; then
	echo "$0 [domain] [private-key] [cert] [ca-bundle]"
	exit 1
fi

output=${domain}.ssl.pem

cat ${private_key} 	> ${output}
cat ${cert} 		>> ${output}
echo ""			>> ${output}
cat ${ca_bundle} 	>> ${output}

echo "Created: ${output}"
