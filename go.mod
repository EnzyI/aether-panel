module github.com/EnzyI/aether-panel

go 1.21

require (
	github.com/docker/docker v25.0.3+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.5.0
	github.com/docker/distribution v2.8.3+incompatible
	go.opentelemetry.io/otel v1.20.0
	go.opentelemetry.io/otel/sdk v1.20.0
	go.opentelemetry.io/otel/trace v1.20.0
	go.opentelemetry.io/otel/metric v1.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.45.0
	golang.org/x/net v0.23.0
	golang.org/x/sys v0.18.0
	golang.org/x/time v0.5.0
	golang.org/x/crypto v0.21.0
	golang.org/x/text v0.14.0
)

replace (
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.20.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.20.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.20.0
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.45.0
	golang.org/x/net => golang.org/x/net v0.23.0
	golang.org/x/sys => golang.org/x/sys v0.18.0
	golang.org/x/time => golang.org/x/time v0.5.0
	golang.org/x/crypto => golang.org/x/crypto v0.21.0
	golang.org/x/text => golang.org/x/text v0.14.0
	github.com/docker/distribution => github.com/docker/distribution v2.8.3+incompatible
)
