module github.com/nfwGytautas/mstk/gomods/balancer-api

go 1.20

replace github.com/nfwGytautas/mstk/gomods/coordinator-api => ../coordinator-api

require github.com/nfwGytautas/mstk/gomods/coordinator-api v0.0.0-00010101000000-000000000000

require github.com/BurntSushi/toml v1.2.1 // indirect