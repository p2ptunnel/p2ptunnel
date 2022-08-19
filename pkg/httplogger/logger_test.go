package httplogger

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLogBasic(t *testing.T) {
	output := &bytes.Buffer{}
	l := New(output)
	req := []byte("GET / HTTP/1.1\r\nHost: 127.0.0.1:8000\r\nUser-Agent: curl/7.79.1\r\nAccept: */*\r\n\r\n")
	l.Print(req)
	if !reflect.DeepEqual(req, output.Bytes()) {
		t.Errorf(" Expect: %s\n Get: %s\n", string(req), output.String())
	}

	expReplyOutput := "HTTP/1.0 200 OK\r\nContent-type: text/html\r\nContent-Length: 297\r\n\r\n"
	reply := []byte(expReplyOutput +
		`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>Directory listing for /</title>
</head>
<body>
<h1>Directory listing for /</h1>
<hr>
<ul>
</ul>
<hr>
</body>
</html>`)

	output.Reset()
	l.Reset()
	l.Print(reply)
	if !reflect.DeepEqual([]byte(expReplyOutput), output.Bytes()) {
		t.Errorf(" Expect: %v\n Get: %v\n", expReplyOutput, output.Bytes())
	}
}

func TestLogWithoutCrLf(t *testing.T) {
	output := &bytes.Buffer{}
	l := New(output)
	reply := []byte(`HTTP/1.0 200 OK
Server: SimpleHTTP/0.6 Python/3.8.9
Date: Thu, 18 Aug 2022 23:54:49 GMT
Content-type: text/html; charset=utf-8
Content-Length: 297

<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>Directory listing for /</title>
</head>
<body>
<h1>Directory listing for /</h1>
<hr>
<ul>
</ul>
<hr>
</body>
</html>`)

	output.Reset()
	l.Reset()
	l.Print(reply)
	if !reflect.DeepEqual(reply, output.Bytes()) {
		t.Errorf(" Expect: %s\n Get: %s\n", string(reply), output.String())
	}
}

func TestLogPartial(t *testing.T) {
	output := &bytes.Buffer{}
	l := New(output)
	for _, req := range [][]byte{
		[]byte("GET / HTTP/1.1"),
		[]byte("GET / HTTP/1.1\r\nHost"),
		[]byte("GET / HTTP/1.1\r\nHost: 127.0.0.1:8000\r\nUser-Agent: curl/7.79.1\r\n"),
	} {
		output.Reset()
		l.Reset()
		l.Print(req)
		if !reflect.DeepEqual(req, output.Bytes()) {
			t.Errorf(" Expect: %s\n Get: %s\n", string(req), output.String())
		}
	}
}

func TestLogMultiple(t *testing.T) {
	output := &bytes.Buffer{}
	l := New(output)
	reply1 := []byte("HTTP/1.0 200 OK\r\nContent-type: text")
	reply2 := []byte("/html\r\nContent-Length: 297\r\n\r\n")
	reply3 := []byte(`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>Directory listing for /</title>
</head>`)
	l.Print(reply1)
	l.Print(reply2)
	if !reflect.DeepEqual(append(reply1, reply2...), output.Bytes()) {
		t.Errorf(" Expect: %s\n Get: %s\n", string(reply1)+string(reply2), output.String())
	}
	l.Print(reply3)
	if !reflect.DeepEqual(append(reply1, reply2...), output.Bytes()) {
		t.Errorf(" Expect: %s\n Get: %s\n", string(reply1)+string(reply2), output.String())
	}
}
