module github.com/EnzyI/aether-panel

go 1.21

require (
	github.com/docker/docker v20.10.27+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.5.0
)

// Chan cac ban thu vien gay loi Go 1.24
exclude (
	github.com/docker/docker v24.0.0+incompatible
	github.com/docker/docker v25.0.0+incompatible
	github.com/docker/docker v26.0.0+incompatible
	golang.org/x/net v0.33.0
)
