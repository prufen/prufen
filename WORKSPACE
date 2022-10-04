workspace(name = "prufen")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "099a9fb96a376ccbbb7d291ed4ecbdfd42f6bc822ab77ae6f1b5cb9e914e94fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "501deb3d5695ab658e82f6f6f549ba681ea3ca2a5fb7911154b5aa45596183fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.26.0/bazel-gazelle-v0.26.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.26.0/bazel-gazelle-v0.26.0.tar.gz",
    ],
)

load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

new_git_repository(
    name = "com_rogchap_v8go",
    # This is external/v8go.BUILD file.
    build_file = "v8go.BUILD",
    # v0.7.0
    commit = "6e4af34cf4447be859741c0719aee06a3d3e7b2a",
    remote = "https://github.com/rogchap/v8go",
    # Suggested by Bazel.
    shallow_since = "1639006196 +1100",
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

go_register_toolchains(version = "1.19.1")

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
    sum = "h1:WTLtQzmQori5FUH25Pq4WT22oCsv8USpQ+F6rqtsmxw=",
    version = "v1.49.0",
)
