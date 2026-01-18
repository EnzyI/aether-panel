module github.com/EnzyI/aether-panel

go 1.21

require (
	github.com/docker/docker v20.10.27+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.5.0
)

// Cấm tuyệt đối các bản thư viện gây lỗi Go 1.24
exclude (
	golang.org/x/net v0.49.0
	golang.org/x/net v0.48.0
	golang.org/x/sys v0.30.0
	golang.org/x/sys v0.29.0
)
