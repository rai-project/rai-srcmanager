workspace(name = "rai_srcmanager")



git_repository(
    name = "io_bazel_rules_go",
    commit = "570488593c55ad61a18c3d6095344f25da8a84e1", # master on 2018-01-04
    remote = "https://github.com/bazelbuild/rules_go",
)

git_repository(
    name = "bazel_gazelle",
    commit = "9e43c85089c3247fece397f95dabc1cb63096a59", # master on 2018-01-09
    remote = "https://github.com/bazelbuild/bazel_gazelle",
)

http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.8/bazel-gazelle-0.8.tar.gz",
    sha256 = "e3dadf036c769d1f40603b86ae1f0f90d11837116022d9b06e4cd88cae786676",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")


go_rules_dependencies()
go_register_toolchains()
gazelle_dependencies()

go_repository(
    name = "com_github_mjibson_esc",
    commit = "58d9cde84f237ecdd89bd7f61c2de2853f4c5c6e",
    importpath = "github.com/mjibson/esc",
)