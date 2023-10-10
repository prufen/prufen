# Prufen

CI-like system for running checks for user-submitted code.

## Development

### `cjail`

You need a container filesystem to run commands inside.

You can download a Docker image for this, e.g. using `skopeo` and `oci-image-tool`:

```
# on Ubuntu 22.04
$ sudo apt install skopeo oci-image-tool
$ skopeo --insecure-policy copy "docker://docker.io/library/busybox" oci:.images:busybox
$ oci-image-tool unpack --ref name=busybox .images-repo .images/busybox/
$ tree -L 2 .images .images-repo/
.images
└── busybox
    ├── bin
    ├── dev
    ├── etc
    ├── home
    ├── lib
    ├── lib64 -> lib
    ├── root
    ├── tmp
    ├── usr
    └── var
.images-repo/
├── blobs
│   └── sha256
├── index.json
└── oci-layout

2 directories, 2 files
```

You need to have [`nsjail`](https://github.com/google/nsjail) binary in `$PATH`.

```
$ bazel run //cjail:cjail -- --images-ref-to-dir-json-map "{\"busybox\": \"$PWD/.images/busybox/\"}"
...
2023/10/10 20:42:35 starting CJail server...
2023/10/10 20:42:35 listening on "localhost:8080"
```

```
$ grpc_cli call localhost:8080 cjail.CJail/ListImages ''
connecting to localhost:8080
image_ref: "busybox"
Rpc succeeded with OK status
```

```
$ grpc_cli call localhost:8080 cjail.CJail/Execute 'base_image_ref: "busybox" file: "/bin/sh" args: ["-c", "id"]'
connecting to localhost:8080
stdout: "uid=1000 gid=1000 groups=65534(nobody),65534(nobody),65534(nobody),65534(nobody),65534(nobody),65534(nobody),1000,65534(nobody)\n"
...
```

## Deployment to `prufen-dev`

### `cjail`

Push `cjail` image:

```
$ bazel run //cjail:push_image
INFO: Analyzed target //cjail:push_image (1 packages loaded, 5 targets configured).
INFO: Found 1 target...
Target //cjail:push_image up-to-date:
  bazel-bin/cjail/push_image.digest
  bazel-bin/cjail/push_image
INFO: Elapsed time: 1.345s, Critical Path: 1.16s
INFO: 10 processes: 1 internal, 9 linux-sandbox.
INFO: Build completed successfully, 10 total actions
INFO: Build completed successfully, 10 total actions
2022/10/30 10:40:49 Successfully pushed Docker image to europe-west1-docker.pkg.dev/prufen-dev/docker-repo/cjail:latest - europe-west1-docker.pkg.dev/prufen-dev/docker-repo/cjail@sha256:b6ad72fcdd59d9b81f548051dd40d7724dbe39fd6780ba6915c9485465c07c63
```

Update image in Cloud run:

```
$ gcloud run services update cjail --image=europe-west1-docker.pkg.dev/prufen-dev/docker-repo/cjail
✓ Deploying... Done.
  ✓ Creating Revision...
  ✓ Routing traffic...
Done.
Service [cjail] revision [cjail-00005-sax] has been deployed and is serving 100 percent of traffic.
Service URL: https://cjail-redacted.a.run.app
```

To update service definition:

```
$ gcloud run services replace cjail/cjail-service.yaml
```

### `jsjail`

```
$ bazel run //jsjail:push_image
$ gcloud run services update jsjail --image=europe-west1-docker.pkg.dev/prufen-dev/docker-repo/jsjail
$ gcloud run services replace jsjail/jsjail-service.yaml
```
