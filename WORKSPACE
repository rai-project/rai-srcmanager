workspace(name = "rai_srcmanager")



git_repository(
    name = "io_bazel_rules_go",
    commit = "570488593c55ad61a18c3d6095344f25da8a84e1", # master on 2018-01-04
    remote = "https://github.com/bazelbuild/rules_go",
)

git_repository(
    name = "bazel_gazelle",
    commit = "9e43c85089c3247fece397f95dabc1cb63096a59", # master on 2018-01-09
    remote = "https://github.com/bazelbuild/bazel-gazelle",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//:esc.bzl", "esc_repositories")



go_rules_dependencies()
go_register_toolchains()
gazelle_dependencies()
esc_repositories()