module github.com/nfwGytautas/mstk/gomods/microservice-api

go 1.20

replace github.com/nfwGytautas/mstk/gomods/coordinator-api => ../../gomods/coordinator-api

replace github.com/nfwGytautas/mstk/gomods/common-api => ../../gomods/common-api

require github.com/nfwGytautas/mstk/gomods/coordinator-api v0.0.0-00010101000000-000000000000

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/nfwGytautas/mstk/gomods/common-api v0.0.0-00010101000000-000000000000 // indirect
)
