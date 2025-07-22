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
    deps = [
        "///third_party/go/github.com_smacker_go-tree-sitter//:go-tree-sitter",
        "///third_party/go/github.com_smacker_go-tree-sitter//javascript",
        "///third_party/go/github.com_smacker_go-tree-sitter//typescript/tsx",
        "///third_party/go/github.com_smacker_go-tree-sitter//typescript/typescript",
        "//kinds",
    ],
)
