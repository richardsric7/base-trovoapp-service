package sms

import (
	"bytes"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

func TestTermiiAPIKeyIsNotLogged(t *testing.T) {
	// a port nobody listens on: the request fails with a transport error,
	// which normally repeats the whole URL
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()

	const key = "secret-termii-key-123"
	t.Setenv("TERMII_SMS_URL", addr)
	t.Setenv("TERMII_SMS_API_KEY", key)
	t.Setenv("TERMII_SMS_SENDER_ID", "Trovo")
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	if err := SendSMSWithTermiiGateway("+44-7700900000", "code 123456"); err == nil {
		t.Fatal("expected the send to fail")
	}
	if strings.Contains(buf.String(), key) {
		t.Fatalf("the API key was logged:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "api_key=REDACTED") || !strings.Contains(buf.String(), "send sms has error") {
		t.Fatalf("expected the redacted URL and the error to be logged:\n%s", buf.String())
	}
}
