load("@io_bazel_rules_go//go:def.bzl", "go_context")
load("@io_bazel_rules_go//go/private:go_repository.bzl", "go_repository")

def _esc_impl(ctx):
  go = go_context(ctx)
  out = go.declare_file(go, ext=".go")
  arguments = ctx.actions.args()
  arguments.add([
      "-o", out.path,
      "-pkg", ctx.attr.package,
      "-prefix", ctx.label.package,
      "-private",
  ])
  arguments.add(ctx.files.srcs)
  ctx.actions.run(
    inputs = ctx.files.srcs,
    outputs = [out],
    mnemonic = "Esc",
    executable = ctx.file._esc,
    arguments = [arguments],
  )
  return [
    DefaultInfo(
      files = depset([out])
    )
  ]

esc = rule(
    _esc_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True, cfg = "data"),
        "package": attr.string(mandatory=True),
        "compress": attr.bool(default=True),
        "metadata": attr.bool(default=False),
        "_esc":  attr.label(allow_files=True, single_file=True, default=Label("//vendor/github.com/mjibson/esc:esc")),
        "_go_context_data": attr.label(default=Label("@io_bazel_rules_go//:go_context_data")),
    },
    toolchains = ["@io_bazel_rules_go//go:toolchain"],
)


def esc_repositories():  
  go_repository(
      name = "com_github_mjibson_esc",
      commit = "58d9cde84f237ecdd89bd7f61c2de2853f4c5c6e",
      importpath = "github.com/mjibson/esc",
  )