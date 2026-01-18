module github.com/EnzyI/aether-panel

go 1.21

require (
	github.com/docker/docker v20.10.27+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.5.0
	golang.org/x/net v0.17.0
	golang.org/x/sys v0.13.0
	golang.org/x/time v0.3.0
)

exclude (
	golang.org/x/net v0.49.0
	golang.org/x/net v0.48.0
	golang.org/x/net v0.33.0
)

replace github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible
