# Prufen

CI-like system for running checks for user-submitted code.

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
