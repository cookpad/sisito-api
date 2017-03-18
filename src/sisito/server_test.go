package sisito

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/ping")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"message":"pong"}`+"\n")
}

func TestRecentWithRecipient(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "RecentlyBounced",
		func(_ *Driver, name string, value string, senderdomain string) (bounced []BounceMail, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("recipient", name)
			assert.Equal("foo@example.com", value)
			assert.Equal("example.net", senderdomain)

			bounced = []BounceMail{BounceMail{Id: 1}}

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/recent?recipient=foo@example.com&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"addresser":"",`+
		`"alias":"",`+
		`"created_at":"0001-01-01T00:00:00Z",`+
		`"deliverystatus":"",`+
		`"destination":"",`+
		`"diagnosticcode":"",`+
		`"digest":"",`+
		`"lhost":"",`+
		`"messageid":"",`+
		`"reason":"",`+
		`"recipient":"",`+
		`"rhost":"",`+
		`"senderdomain":"",`+
		`"smtpagent":"",`+
		`"smtpcommand":"",`+
		`"softbounce":false,"subject":"",`+
		`"timestamp":"0001-01-01T00:00:00Z",`+
		`"timezoneoffset":"",`+
		`"updated_at":"0001-01-01T00:00:00Z",`+
		`"whitelisted":false}`+"\n")
}

func TestRecentWithDigest(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "RecentlyBounced",
		func(_ *Driver, name string, value string, senderdomain string) (bounced []BounceMail, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("digest", name)
			assert.Equal("767e74eab7081c41e0b83630511139d130249666", value)
			assert.Equal("", senderdomain)

			bounced = []BounceMail{BounceMail{Id: 1}}

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/recent?digest=767e74eab7081c41e0b83630511139d130249666")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"addresser":"",`+
		`"alias":"",`+
		`"created_at":"0001-01-01T00:00:00Z",`+
		`"deliverystatus":"",`+
		`"destination":"",`+
		`"diagnosticcode":"",`+
		`"digest":"",`+
		`"lhost":"",`+
		`"messageid":"",`+
		`"reason":"",`+
		`"recipient":"",`+
		`"rhost":"",`+
		`"senderdomain":"",`+
		`"smtpagent":"",`+
		`"smtpcommand":"",`+
		`"softbounce":false,"subject":"",`+
		`"timestamp":"0001-01-01T00:00:00Z",`+
		`"timezoneoffset":"",`+
		`"updated_at":"0001-01-01T00:00:00Z",`+
		`"whitelisted":false}`+"\n")
}

func TestRecentWithRecipientDigest(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/recent?recipient=foo@example.com&digest=767e74eab7081c41e0b83630511139d130249666&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"Cannot pass both \"recipient\" and \"digest\""}`+"\n")
}

func TestRecentWithoutRecipientDigest(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/recent?senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"\"recipient\" or \"digest\" is not present"}`+"\n")
}

func TestBouncedWithRecipient(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "IsBounced",
		func(_ *Driver, name string, value string, senderdomain string) (bounced bool, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("recipient", name)
			assert.Equal("foo@example.com", value)
			assert.Equal("example.net", senderdomain)

			bounced = true

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/bounced?recipient=foo@example.com&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"bounced":true}`+"\n")
}

func TestBouncedWithDigest(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "IsBounced",
		func(_ *Driver, name string, value string, senderdomain string) (bounced bool, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("digest", name)
			assert.Equal("767e74eab7081c41e0b83630511139d130249666", value)
			assert.Equal("example.net", senderdomain)

			bounced = false

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/bounced?digest=767e74eab7081c41e0b83630511139d130249666&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"bounced":false}`+"\n")
}

func TestBouncedWithRecipientDigest(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/bounced?recipient=foo@example.com&digest=767e74eab7081c41e0b83630511139d130249666&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"Cannot pass both \"recipient\" and \"digest\""}`+"\n")
}

func TestBouncedWithoutRecipientDigest(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/bounced?senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"\"recipient\" or \"digest\" is not present"}`+"\n")
}
