load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

cc_library(
    name = "v8",
    srcs = [
        "deps/linux_x86_64/libv8.a",
    ],
    hdrs = glob([
        "deps/include/*.h",
        "deps/include/cppgc/*.h",
        "deps/include/cppgc/internal/*.h",
        "deps/include/libplatform/*.h",
    ]),
    includes = [
        "deps/include",
    ],
)

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude=[
            # Exclude tests, they are using "package v8go_test"
            "*_test.go",
            # Exclude CGO settings and vendor dependencies, they are set explicitly in Bazel.
            "cgo.go",
        ],
    ) + [
        "v8go.cc",
        "v8go.h",
        # Not a test, but v8go methods exported for tests only.
        "export_test.go",
    ],
    importpath = "rogchap.com/v8go",
    visibility = ["//visibility:public"],
    cgo = True,
    cdeps = [
        ":v8",
    ],
    cxxopts = [
        # Flags set in cgo.go (without `-I${SRCDIR}/deps/include`):
        # https://github.com/rogchap/v8go/blob/0e40e6e5827ad897d25f915917d8206e9a8231db/cgo.go#L9
        "-fno-rtti", "-fPIC", "-std=c++17", "-DV8_COMPRESS_POINTERS", "-DV8_31BIT_SMIS_ON_64BIT_ARCH", "-Wall", "-DV8_ENABLE_SANDBOX",
    ],
)

go_test(
    name = "go_default_test",
    srcs = glob(
        ["*_test.go"],
        exclude=[
            # Not a test, but v8go methods exported for tests only.
            "export_test.go",
        ],
    ),
    deps = [
        ":go_default_library",
    ],
)
