package runtime_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func TestAnnotateContext(t *testing.T) {
	ctx := context.Background()

	request, _ := http.NewRequest("GET", "http://www.example.com", nil)
	request.Header = http.Header{}
	request.Header.Add("Grpc-Metadata-FooBar", "Value1")
	request.Header.Add("Grpc-Metadata-Foo-BAZ", "Value2")
	request.RemoteAddr = "192.168.0.1:12345"
	annotated := runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != 4 {
		t.Errorf("Expected 4 metadata items in context; got %v", md)
	}
	if got, want := md["Foobar"], []string{"Value1"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"Foobar\"] = %v; want %v", got, want)
	}
	if got, want := md["Foo-Baz"], []string{"Value2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"Foo-Baz\"] = %v want %v", got, want)
	}
}

func TestAnnotateContextPassesNonGrpcMetadata(t *testing.T) {
	ctx := context.Background()

	request, _ := http.NewRequest("GET", "http://bar.foo.example.com", nil)
	request.Header = http.Header{}
	request.Header.Add("Authorization", "Bearer FAKETOKEN")
	request.Header.Add("X-Forwarded-For", "8.8.8.8")
	request.RemoteAddr = "192.168.0.100:12345"
	annotated := runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != 3 {
		t.Errorf("Expected 2 metadata items in context; got %v", md)
	}
	if got, want := md["x-forwarded-host"], []string{"bar.foo.example.com"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"host\"] = %v; want %v", got, want)
	}
	if got, want := md["authorization"], []string{"Bearer FAKETOKEN"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"authorization\"] = %v want %v", got, want)
	}
	if got, want := md["x-forwarded-for"], []string{"8.8.8.8, 192.168.0.100"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"x-forwarded-for\"] = %v want %v", got, want)
	}
}
