module github.com/gophish/gophish

go 1.23.5

require (
	bitbucket.org/liamstask/goose v0.0.0-20150115234039-8488cc47d90c
	github.com/NYTimes/gziphandler v1.1.1
	github.com/PuerkitoBio/goquery v1.10.2
	github.com/emersion/go-imap v1.2.1
	github.com/emersion/go-message v0.18.2
	github.com/go-sql-driver/mysql v1.9.0
	github.com/gophish/gomail v0.0.0-20200818021916-1f6d0dfd512e
	github.com/gorilla/context v1.1.2
	github.com/gorilla/csrf v1.7.2
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/sessions v1.4.0
	github.com/jinzhu/gorm v1.9.16
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/jordan-wright/unindexed v0.0.0-20181209214434-78fa79113c0f
	github.com/mattn/go-sqlite3 v1.14.24 //changed to v1.14.24 as v2.0.3+incompatable was retracted
	github.com/oschwald/maxminddb-golang v1.13.1
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.35.0 // CVE-2024-45337 fix
	golang.org/x/time v0.11.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require gopkg.in/alecthomas/kingpin.v2 v2.4.0

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/emersion/go-sasl v0.0.0-20241020182733-b788ff22d5a6 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/go-gypsy v1.0.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

replace gopkg.in/alecthomas/kingpin.v2 => github.com/alecthomas/kingpin/v2 v2.4.0
