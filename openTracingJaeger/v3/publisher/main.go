package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

// curl 'http://localhost:8082/publish?helloStr=hi%20there'
func main() {
	tracer, closer := Init("publisher")
	defer closer.Close()

	//这里的extract相当于接收到client端的请求,把ctx和span接收回来,这样就是从client,server端形成了一条链路进行追踪
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloStr := r.FormValue("helloStr")
/*		span.LogFields(
			otlog.String("event","string-format"),
			otlog.String("value",helloStr),
			)*/

		println(helloStr)
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}