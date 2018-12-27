load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "7949fc6cc17b5b191103e97481cf8889217263acf52e00b560683413af204fcb",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.16.0/bazel-gazelle-0.16.0.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_alecthomas_units",
    commit = "2efee857e7cfd4f3d0138cc3cbb1b4966962b93a",
    importpath = "github.com/alecthomas/units",
)

go_repository(
    name = "com_github_alecthomas_template",
    commit = "a0175ee3bccc567396460bf5acd36800cb10c49c",
    importpath = "github.com/alecthomas/template",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "3a771d992973f24aa725d07868b467d1ddfceafb",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "3ba18b2033700609303ccae90aebc49a8b460ca6",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "5c3871d89910bfb32f5fcab2aa4b9ec68e65a99f",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "67670fe90761d7ff18ec1d640135e53b9198328f",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "1dc9a6cbc91aacc3e8b2d63db4d2e957a5394ac4",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "6180652bb05834ae178c93911ddae14b689408b7",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "in_gopkg_alecthomas_kingpin_v2",
    commit = "947dcec5ba9c011838740e680966fd7087a71d0d",
    importpath = "gopkg.in/alecthomas/kingpin.v2",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "505ab145d0a99da450461ae2c1a9f6cd10d1f447",
    importpath = "golang.org/x/crypto",
)
