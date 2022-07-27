#!/usr/bin/env bash

NAMESPACE=example-205
SERVICE_NAME=example-205-mutating-service

COMMON_NAME=${SERVICE_NAME}.${NAMESPACE}

set -eo pipefail

echo =================================================================
echo THIS SCRIPT IS NOT SECURE! USE IT ONLY FOR DEMONSTATION PURPOSES.
echo =================================================================
echo

mkdir -p mutating-certs && cd mutating-certs

if [[ -e ca.csr  ]] ; then
  read -p "Regenerate certificates? (yes/no) [no]: "
  if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]
  then
    exit 0
  fi
fi

RM_FILES="ca* cert*"
echo ">>> Remove ${RM_FILES}"
rm -f $RM_FILES

echo ">>> Generate CA key and certificate"
cat <<EOF | cfssl gencert -initca - | cfssljson -bare ca
{
  "CN": "Shell-operator example 205-mutating-webhook Root CA",
  "key": {
    "algo": "rsa",
    "size": 2058
  }
}
EOF


CFSSL_CONFIG=$(cat <<EOF
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "server": {
        "usages": [
          "signing",
          "digital signing",
          "key encipherment",
          "server auth"
        ],
        "expiry": "8760h"
      }
    }
  }
}
EOF
)

echo ">>> Generate cert.key and cert.crt"
cat <<EOF | cfssl gencert -ca ca.pem -ca-key ca-key.pem -config <(echo "$CFSSL_CONFIG") -profile=server - | cfssljson -bare tls
{
  "CN": "${COMMON_NAME}.svc",
  "key": {
    "algo": "rsa",
    "size": 2058
  },
  "hosts": [
    "${COMMON_NAME}",
    "${COMMON_NAME}.svc",
    "${COMMON_NAME}.svc.cluster.local"
  ]
}
EOF
