package codec

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/require"
	"net/http/httputil"
	"bytes"
	"fmt"
	"bufio"
)

func Test_bufio(t *testing.T) {
	should := require.New(t)
	req, err := http.NewRequest("GET", "/", nil)
	should.Nil(err)
	reqBytes, err := httputil.DumpRequest(req, true)
	should.Nil(err)
	buf := &bytes.Buffer{}
	buf.Write(reqBytes)
	buf.Write(reqBytes)
	reader := bufio.NewReaderSize(buf, 2048)
	fmt.Println(http.ReadRequest(reader))
	fmt.Println(reader.Buffered())
	fmt.Println(http.ReadRequest(reader))
	fmt.Println(reader.Buffered())
}