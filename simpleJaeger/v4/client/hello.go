package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
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


// go run hello.go Icarus Bonjour
func main() {
	if len(os.Args) != 3 {
		panic("ERROR: Expecting two arguments")
	}

	tracer, closer := Init("hello-world")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	helloTo := os.Args[1]
	greeting := os.Args[2]

	//这里的第二个参数Bonjour存入了baggage中,其中key是'greeting'
	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", helloTo)
	span.SetBaggageItem("greeting", greeting)
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://localhost:8081/format?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := Do(req)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}

	helloStr := string(resp)

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloStr", helloStr)
	url := "http://localhost:8082/publish?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	if _, err := Do(req); err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}
}

func Do(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}
