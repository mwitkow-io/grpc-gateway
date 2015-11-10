package runtime

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	metadataHeaderPrefix = "Grpc-Metadata-"
	xForwardedFor = "X-Forwarded-For"
	xForwardedHost = "X-Forwarded-Host"
)

/*
AnnotateContext adds context information such as metadata from the request.

If there are no metadata headers in the request, then the context returned
will be the same context.
*/
func AnnotateContext(ctx context.Context, req *http.Request) context.Context {
	var pairs []string
	for key, val := range req.Header {
		if strings.HasPrefix(key, metadataHeaderPrefix) {
			pairs = append(pairs, key[len(metadataHeaderPrefix):], val[0])
		}
		if key == "Authorization" {
			pairs = append(pairs, "authorization", val[0])
		}
	}
	if req.Header.Get(xForwardedHost) != "" {
		pairs = append(pairs, strings.ToLower(xForwardedHost), req.Header.Get(xForwardedHost))
	} else if req.Host != "" {
		pairs = append(pairs, strings.ToLower(xForwardedHost), req.Host)
	}
	if req.Header.Get(xForwardedFor) == "" {
		pairs = append(pairs, strings.ToLower(xForwardedFor), req.RemoteAddr)
	} else {
		pairs = append(pairs, strings.ToLower(xForwardedFor), req.Header.Get(xForwardedFor) + ", " + req.RemoteAddr)
	}

	if len(pairs) != 0 {
		ctx = metadata.NewContext(ctx, metadata.Pairs(pairs...))
	}
	return ctx
}
