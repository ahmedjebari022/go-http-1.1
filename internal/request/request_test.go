package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}


func TestRequestLineParse(t *testing.T){

	// Test: Standard Body
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)

	// r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// r, err = RequestFromReader(strings.NewReader("POST /coffee HTTP/1.1\r\n"))	
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "POST", r.RequestLine.Method)
	// assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1",r.RequestLine.HttpVersion)

	// _, err = RequestFromReader(strings.NewReader("FALSE /coffee HTTP/1.1"))
	// require.Error(t, err)


	// _, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.2"))
	// require.Error(t, err)

	// _, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	// require.Error(t, err)
	// Test: Good GET Request line
	// reader := &chunkReader{
	// 	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 	numBytesPerRead: 3,
	// }
	// r, err := RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// // Test: Good GET Request line with path
	// reader = &chunkReader{
	// 	data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 	numBytesPerRead: 1,
	// }
	// r, err = RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// reader.data = "FOK /coffee HTTP/1.1\r\nHost: localhost:42069"
	// r, err = RequestFromReader(reader)
	// require.Error(t, err)

	// reader.data = "POST /coffee HTTP/1.0\r\nHost: localhost:42069\r\n"
	// r, err = RequestFromReader(reader)
	// require.Error(t, err)

	// reader.data = "POST HTTP/1.0\r\nHost: localhost:42069\r\n"
	// r, err = RequestFromReader(reader)
	// require.Error(t, err)

	
	// reader.data = "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	// reader.pos = 0
	// r, err = RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "POST", r.RequestLine.Method)
	// assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	// assert.Equal(t, "curl/7.81.0",r.Header["user-agent"])
	// assert.Equal(t, "localhost:42069", r.Header["host"])

	// reader = &chunkReader{
	// 	data:  "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nHost: curl/7.81.0\r\nHost: ddaadarl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 	pos: 0,
	// 	numBytesPerRead: 3,		
	// }
	// r, err = RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "POST", r.RequestLine.Method)
	// assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	// assert.Equal(t, "localhost:42069, curl/7.81.0, ddaadarl/7.81.0",r.Header["host"])

	// reader.data = "POST /coffee HTTP/1.1`\r\nHost :localhost:42069\r\n\r\n"
	// reader.pos = 0
	// r, err = RequestFromReader(reader)
	// require.Error(t, err)

	// reader.data = "POST /coffee HTTP/1.1`\r\nHost: localhost:42069"
	// reader.pos = 0
	// require.Error(t, err)
	
}