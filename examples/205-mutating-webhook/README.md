# Example with mutating hooks

This is a simple example of `kubernetesMutating` binding. Read more information in [BINDING_MUTATING.md](../../BINDING_MUTATING.md).

## Run

### Generate certificates

An HTTP server behind the MutatingWebhookConfiguration requires a certificate issued by the CA. For simplicity, this process is automated with `gen-certs.sh` script. Just run it:

```
./gen-certs.sh
```

> Note: `gen-certs.sh` requires [cfssl utility](https://github.com/cloudflare/cfssl/releases/latest).

### Build and install example

Build Docker image and use helm3 to install it:

```
docker build -t localhost:5000/shell-operator:example-205 .
docker push localhost:5000/shell-operator:example-205
helm upgrade --install \
    --namespace example-205 \
    --create-namespace \
    example-205 .
```

### Cleanup

```
helm delete --namespace=example-205 example-205
kubectl delete ns example-205
kubectl delete mutatingwebhookconfiguration/example-205
```
