## Tracing

start Zipkin: java -jar zipkin.jar  
Zipkin Web UI: http://localhost:9411  


start Jaedger: jaeger-all-in-one --collector.zipkin.host-port=:9411  
Jaeger Web UI: http://localhost:16686  


opencensus zipkin : https://opencensus.io/quickstart/go/tracing/  
jaedger : https://pkg.go.dev/contrib.go.opencensus.io/exporter/jaeger


