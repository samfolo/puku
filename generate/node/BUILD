go_test(
    name = "node_test",
    srcs = [
        "builtins_test.go",
        "import_test.go",
    ],
    deps = [
        ":node",
        "///third_party/go/github.com_stretchr_testify//assert",
        "//kinds",
    ],
)

go_library(
    name = "node",
    srcs = [
        "builtins.go",
        "import.go",
    ],
    deps = ["//kinds"],
)
