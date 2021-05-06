module myTracing

go 1.16

require (
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.28.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.opencensus.io v0.23.0
	go.uber.org/atomic v1.7.0 // indirect
)
