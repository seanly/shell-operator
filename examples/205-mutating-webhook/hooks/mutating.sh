#!/usr/bin/env bash

source /shell_lib.sh

function __config__(){
    cat <<EOF
configVersion: v1
kubernetesMutating:
- name: inject-pod.example.com
  namespace:
    labelSelector:
      matchLabels:
        name: example-205
  rules:
  - apiGroups:   [""]
    apiVersions: ["*"]
    operations:  ["CREATE"]
    resources:   ["pods"]
    scope:       "Namespaced"

EOF
}

function __on_mutating::inject-pod.example.com() {

  cat <<EOF > $MUTATING_RESPONSE_PATH
{ 
  "allowed":true, 
  "patchOps": [
    {"op": "add", "path": "/spec/nodeSelector", "value": {"hello": "world"}} 
  ]
}
EOF

}

hook::run $@
