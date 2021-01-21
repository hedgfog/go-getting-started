// +heroku install ./cmd/...

module cilintservice

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210113205817-d3ed898aa8a3
)
