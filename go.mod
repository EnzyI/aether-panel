module github.com/EnzyI/aether-panel

go 1.21

require (
	github.com/docker/docker v20.10.27+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.5.0
)

// Ep dung phien ban thu vien cu de khong bi loi SplitHostname
replace github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible

exclude (
	golang.org/x/net v0.33.0
	golang.org/x/sys v0.28.0
)
