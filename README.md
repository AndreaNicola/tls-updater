# TLS UPDATER

Stupid tool used to sync kubernetes.io/tls secrets.

# Build

Build with :

```bash
docker build -t tls-updater:0.0.1 .
```

# Deploy

Apply the k8s the manifest attached with this repo (kubernetes-manifest.yaml)

# Usage

The tool will look for all secrets (in all namespaces) labelled with the label 

```
tls-updater: 'true'
```

Add an annotation to your secret (tls-updater-dests) containing the names of the secret that must be synced. The secrets MUST be in the same namespace of the original secret.

A corrected configured tls secret should look like this one:

```
kind: Secret
apiVersion: v1
metadata:
  name: original-secret
  namespace: secret-namespace
  labels:
    tls-updater: 'true'
  annotations:
    tls-updater-dests: destination-tls-secret-0,destination-tls-secret-1
data:
  tls.crt: >-
    ...
  tls.key: >-
    ...
type: kubernetes.io/tls
```

## Sample

You can see a sample with some secrets to sync in tls-sample.secret.yaml file.

When you will apply this configuration the operator will copy the key-pairs, from secret-0 to secrets 1 and 2, if they're not equals.

You will see this logs in the tls-updater pod when keypairs are not equals:

```
Secret created: secret-1. Checking TLS Cert...
Secret changed secret-1. Updating secret secret-2!!
Secret changed secret-1. Updating secret secret-3!!
```

If the keypair are equals the log will be similar to this:

```
Secret created: secret-1. Checking TLS Cert...
Secret secret-2 is synched with secret secret-1. Not updating keypair
Secret secret-3 is synched with secret secret-1. Not updating keypair
```

Happy TLS Sync!