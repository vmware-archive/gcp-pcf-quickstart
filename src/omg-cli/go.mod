module omg-cli

require (
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/gosuri/uilive v0.0.3
	github.com/iancoleman/strcase v0.0.0-20181128000000-3605ed457bf7
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.6.0
	github.com/pivotal-cf/kiln v0.0.0-20180329191310-9c0f5ac8553d // indirect
	github.com/pivotal-cf/om v0.0.0-20190816215002-d607995f0947
	github.com/shurcooL/httpfs v0.0.0-20190527155220-6a4d4a70508b // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd // indirect
	github.com/starkandwayne/om-tiler v0.0.0-20190820103743-86c0a1263e12
	github.com/tmc/scp v0.0.0-20170824174625-f7b48647feef
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.8.0
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/cheggaaa/pb => github.com/cheggaaa/pb v1.0.28 // from bosh-cli Gopkg.lock
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422 // https://github.com/golang/lint/issues/446#issuecomment-483638233
	github.com/jessevdk/go-flags => github.com/cppforlife/go-flags v0.0.0-20170707010757-351f5f310b26 // from bosh-cli Gopkg.lock
)
