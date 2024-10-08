module rye-gio

go 1.23

toolchain go1.23.2

replace github.com/refaktor/rye => ../rye
replace github.com/refaktor/ryegen => ../ryegen

require (
	gioui.org v0.7.1
	github.com/go-text/typesetting v0.1.1
	github.com/refaktor/rye v0.0.24
	github.com/refaktor/ryegen v0.1.1-0.20241004110956-28a8a8e5f0c5
	golang.org/x/image v0.18.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.8 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/RoaringBitmap/roaring v1.9.4 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/aws/aws-sdk-go-v2 v1.31.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.27.40 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.38 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/ses v1.27.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.23.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.27.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.31.4 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/bitfield/script v0.23.0 // indirect
	github.com/bits-and-blooms/bitset v1.14.3 // indirect
	github.com/blevesearch/bleve/v2 v2.4.2 // indirect
	github.com/blevesearch/bleve_index_api v1.1.12 // indirect
	github.com/blevesearch/geo v0.1.20 // indirect
	github.com/blevesearch/go-faiss v1.0.22 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/gtreap v0.1.1 // indirect
	github.com/blevesearch/mmap-go v1.0.4 // indirect
	github.com/blevesearch/scorch_segment_api/v2 v2.2.16 // indirect
	github.com/blevesearch/segment v0.9.1 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/upsidedown_store_api v1.0.2 // indirect
	github.com/blevesearch/vellum v1.0.10 // indirect
	github.com/blevesearch/zapx/v11 v11.3.10 // indirect
	github.com/blevesearch/zapx/v12 v12.3.10 // indirect
	github.com/blevesearch/zapx/v13 v13.3.10 // indirect
	github.com/blevesearch/zapx/v14 v14.3.10 // indirect
	github.com/blevesearch/zapx/v15 v15.3.13 // indirect
	github.com/blevesearch/zapx/v16 v16.1.5 // indirect
	github.com/drewlanenga/govector v0.0.0-20220726163947-b958ac08bc93 // indirect
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gliderlabs/ssh v0.3.7 // indirect
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/golang/geo v0.0.0-20230421003525-6adc56603217 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/itchyny/gojq v0.12.16 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kopoli/go-terminal-size v0.0.0-20170219200355-5c97524c8b54 // indirect
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/lufia/plan9stats v0.0.0-20240909124753-873cd0166683 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.23 // indirect
	github.com/mhale/smtpd v0.8.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mrz1836/postmark v1.6.5 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/term v1.2.0-beta.2.0.20211217091447-1a4a3b719465 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/refaktor/go-peg v0.0.0-20220116201714-31e3dfa8dc7d // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sashabaranov/go-openai v1.31.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/thomasberger/parsemail v1.2.7 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/bbolt v1.3.11 // indirect
	go.mongodb.org/mongo-driver v1.17.1 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20240707233637-46b078467d37 // indirect
	golang.org/x/exp/shiny v0.0.0-20240707233637-46b078467d37 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.24.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	mvdan.cc/sh/v3 v3.9.0 // indirect
)
