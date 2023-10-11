workspace(name = "prufen")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "56d8c5a5c91e1af73eca71a6fab2ced959b67c86d12ba37feedb0a2dfea441a6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "448e37e0dbf61d6fa8f00aaa12d191745e14f07c31cabfa731f0c8e8a4f41b97",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.28.0/bazel-gazelle-v0.28.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.28.0.tar.gz",
    ],
)

load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

new_git_repository(
    name = "com_rogchap_v8go",
    # This is external/v8go.BUILD file.
    build_file = "v8go.BUILD",
    # v0.9.0
    # Note: keep in sync with the version imported in go.mod
    commit = "0e40e6e5827ad897d25f915917d8206e9a8231db",
    remote = "https://github.com/rogchap/v8go",
    # Suggested by Bazel.
    shallow_since = "1679354875 -0400",
)

# An example of using local (and patched) version of the repository:
# new_local_repository(
#     name = "com_rogchap_v8go",
#     build_file = "external/v8go.BUILD",
#     # Path to the local fork of github.com/rogchap/v8go
#     path = "../v8go",
# )

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("//:gazelle_repos.bzl", "gazelle_repositories")

# gazelle:repository_macro gazelle_repos.bzl%gazelle_repositories
gazelle_repositories()

go_rules_dependencies()

go_register_toolchains(version = "1.19.5")

gazelle_dependencies()

http_archive(
    name = "com_google_protobuf",
    sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
    strip_prefix = "protobuf-3.14.0",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

# If protobuf stuff is put before gazelle stuff, gazelle get's confused with errors like:
#     ... every rule of type _gazelle_runner implicitly depends upon the target '@bazel_gazelle//internal:gazelle.bash.in',
#     but this target could not be found because ...
protobuf_deps()

go_repository(
    name = "org_golang_google_grpc",
    # Disable generation from proto file to overcome issues like
    # https://github.com/bazelbuild/bazel-gazelle/issues/1058
    # as documented in
    # https://github.com/bazelbuild/rules_go/blob/5d306c433cebb1ae8a7b72df2a055be2bacbb12b/go/dependencies.rst#grpc-dependencies
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    sum = "h1:BjnpXut1btbtgN/6sp+brB2Kbm2LjNXnidYujAVbSoQ=",
    version = "v1.58.3",
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "b1e80761a8a8243d03ebca8845e9cc1ba6c82ce7c5179ce2b295cd36f7e394bf",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.25.0/rules_docker-v0.25.0.tar.gz"],
)

load("@io_bazel_rules_docker//repositories:repositories.bzl", container_repositories = "repositories")

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//go:image.bzl", _go_image_repos = "repositories")

_go_image_repos()

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

# Use recent Ubuntu image as a base for JSJail (that uses C++-based V8).
# v8go builds V8 using Ubuntu 18.04 image:
# https://github.com/rogchap/v8go/blob/fc8b9f1095704bc00c8b1b065e4834cadf2802c6/.github/workflows/v8build.yml#L17
# and Bazel builds remaining CGO code using whatever C++ toolchain is configured.
# Using current Go or C++ distroless images producess errors like:
#     /app/jail/js/js: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.34' not found (required by /app/jail/js/js)
#     /app/jail/js/js: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.32' not found (required by /app/jail/js/js)
#     /app/jail/js/js: /usr/lib/x86_64-linux-gnu/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by /app/jail/js/js)
#     /app/jail/js/js: /usr/lib/x86_64-linux-gnu/libstdc++.so.6: version `GLIBCXX_3.4.30' not found (required by /app/jail/js/js)
container_pull(
    name = "jsjail_image_base",
    # Suggested by Bazel.
    digest = "sha256:a8fe6fd30333dc60fc5306982a7c51385c2091af1e0ee887166b40a905691fd0",
    registry = "index.docker.io",
    repository = "library/ubuntu",
    # Ubuntu 22.04 updated to a specific day.
    tag = "jammy-20221003",
)

# It's too hard to build nsjail under Bazel, so let's just build it in Docker (supported by upstream) and use it as a base image.
container_pull(
    name = "cjail_image_base",
    digest = "sha256:4c17e599229bc1782239b2e7db1013bf42fb7b8080d7fe0a5cb6ca06ab23f4eb",
    registry = "index.docker.io",
    repository = "rutsky/jail-base",
    tag = "6",
)
