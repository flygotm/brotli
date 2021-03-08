package brotli

import (
	"bytes"
	c "github.com/billcoding/flygo/context"
	"github.com/billcoding/flygo/middleware"
	"github.com/billcoding/flygo/mime"
	"log"
	"os"
	"strings"
)

//Define brotli struct
type brotli struct {
	logger      *log.Logger
	contentType []string
	minSize     int
	quality     int
	lgWin       int
}

//New
func New() *brotli {
	return &brotli{
		logger: log.New(os.Stdout, "[BROTLI]", log.LstdFlags),
		contentType: []string{
			"application/javascript",
			"application/json",
			"application/xml",
			"text/javascript",
			"text/json",
			"text/xml",
			"text/plain",
			"text/xml",
			"html/css",
		},
		minSize: 2 << 9, //1KB
		quality: 6,
	}
}

//Name
func (b *brotli) Name() string {
	return "Brotli"
}

//Type
func (b *brotli) Type() *middleware.Type {
	return middleware.TypeAfter
}

//Method
func (b *brotli) Method() middleware.Method {
	return middleware.MethodAny
}

//Pattern
func (b *brotli) Pattern() middleware.Pattern {
	return middleware.PatternAny
}

func (b *brotli) accept(c *c.Context) bool {
	acceptEncoding := c.Request.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "br")
}

//Handler
func (b *brotli) Handler() func(c *c.Context) {
	return func(ctx *c.Context) {
		if b.accept(ctx) && ctx.Render().Rended() {
			odata := ctx.Render().Buffer
			if nil != odata && len(odata) >= b.minSize {
				ct := ctx.Render().ContentType
				if strings.Index(ct, ";") != -1 {
					ct = strings.TrimSpace(strings.Split(ct, ";")[0])
				}
				if ct == "" {
					ct = mime.BINARY
				}
				ctx.Header().Set("Vary", "Content-Encoding")
				ctx.Header().Set("Content-Encoding", "br")
				var buffers bytes.Buffer
				var bw = NewWriter(&buffers, WriterOptions{Quality: b.quality, LGWin: b.lgWin})
				defer bw.Close()
				_, werr := bw.Write(odata)
				if werr != nil {
					b.logger.Println(werr)
				}
				bw.Flush()
				ctx.Write(buffers.Bytes())
			}
		}
		ctx.Chain()
	}
}

//ContentType
func (b *brotli) ContentType(contentType ...string) *brotli {
	b.contentType = contentType
	return b
}

//MinSize
func (b *brotli) MinSize(minSize int) *brotli {
	b.minSize = minSize
	return b
}

//Quality
func (b *brotli) Quality(quality int) *brotli {
	b.quality = quality
	return b
}

//LGWin
func (b *brotli) LGWin(lgWin int) *brotli {
	b.lgWin = lgWin
	return b
}
