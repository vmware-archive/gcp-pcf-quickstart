module omg-cli

require (
	cloud.google.com/go v0.37.0 // indirect
	git.apache.org/thrift.git v0.12.0 // indirect
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/aws/aws-sdk-go v1.18.0 // indirect
	github.com/gosuri/uilive v0.0.0-20170323041506-ac356e6e42cd
	github.com/grpc-ecosystem/grpc-gateway v1.6.2 // indirect
	github.com/iancoleman/strcase v0.0.0-20181128000000-3605ed457bf7
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kylelemons/godebug v0.0.0-20170820004349-d65d576e9348 // indirect
	github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.4.3
	github.com/pivotal-cf/om v0.0.0-20190308185307-fa1f978a1ddb
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd // indirect
	github.com/starkandwayne/om-tiler v0.0.0-20190418111437-da6bd32159ca
	github.com/thadc23/yamldiff v0.0.0-20181024205526-2f3964ccb7da // indirect
	github.com/tmc/scp v0.0.0-20170824174625-f7b48647feef
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421
	golang.org/x/sys v0.0.0-20190312061237-fead79001313 // indirect
	google.golang.org/api v0.2.0
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/graymeta/stow => github.com/jtarchie/stow v0.0.0-20190209005554-0bff39424d5b
	github.com/jessevdk/go-flags => github.com/cppforlife/go-flags v0.0.0-20170707010757-351f5f310b26
	gopkg.in/mattn/go-colorable.v0 => github.com/mattn/go-colorable v0.1.1
)
