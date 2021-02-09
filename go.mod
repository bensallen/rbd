module github.com/bensallen/rbd

go 1.14

require (
	github.com/google/goexpect v0.0.0-20200816234442-b5b77125c2c5 // indirect
	github.com/rekby/gpt v0.0.0-20200614112001-7da10aec5566 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/u-root/u-root v7.0.0+incompatible
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	golang.org/x/tools v0.1.0 // indirect
)

replace github.com/u-root/u-root v7.0.0+incompatible => github.com/u-root/u-root v1.0.1-0.20201119150355-04f343dd1922
