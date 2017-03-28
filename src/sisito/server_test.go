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
		reflect.TypeOf(driver), "RecentlyListed",
		func(_ *Driver, name string, value string, senderdomain string) (listed []BounceMail, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("recipient", name)
			assert.Equal("foo@example.com", value)
			assert.Equal("example.net", senderdomain)

			listed = []BounceMail{BounceMail{Id: 1}}

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
		`"softbounce":false,`+
		`"subject":"",`+
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
		reflect.TypeOf(driver), "RecentlyListed",
		func(_ *Driver, name string, value string, senderdomain string) (listed []BounceMail, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("digest", name)
			assert.Equal("767e74eab7081c41e0b83630511139d130249666", value)
			assert.Equal("", senderdomain)

			listed = []BounceMail{BounceMail{Id: 1}}

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
		`"softbounce":false,`+
		`"subject":"",`+
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

func TestListedWithRecipient(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "Listed",
		func(_ *Driver, name string, value string, senderdomain string) (listed bool, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("recipient", name)
			assert.Equal("foo@example.com", value)
			assert.Equal("example.net", senderdomain)

			listed = true

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/listed?recipient=foo@example.com&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"listed":true}`+"\n")
}

func TestListedWithDigest(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "Listed",
		func(_ *Driver, name string, value string, senderdomain string) (listed bool, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("digest", name)
			assert.Equal("767e74eab7081c41e0b83630511139d130249666", value)
			assert.Equal("example.net", senderdomain)

			listed = false

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/listed?digest=767e74eab7081c41e0b83630511139d130249666&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"listed":false}`+"\n")
}

func TestListedWithRecipientDigest(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/listed?recipient=foo@example.com&digest=767e74eab7081c41e0b83630511139d130249666&senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"Cannot pass both \"recipient\" and \"digest\""}`+"\n")
}

func TestListedWithoutRecipientDigest(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/listed?senderdomain=example.net")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"\"recipient\" or \"digest\" is not present"}`+"\n")
}

func TestBlacklist(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "BlacklistRecipients",
		func(_ *Driver, senderdomain string, reasons []string, softbounce *bool, limit uint64, offset uint64) (recipients []string, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("example.net", senderdomain)
			assert.Equal([]string{"userunknown", "filtered"}, reasons)
			assert.Equal(true, *softbounce)
			assert.Equal(uint64(100), limit)
			assert.Equal(uint64(100), offset)

			recipients = []string{"foo@example.com"}

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/blacklist" +
		"?senderdomain=example.net&reason=userunknown&reason=filtered&softbounce=1&limit=100&offset=100")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"recipients":["foo@example.com"]}`+"\n")
}

func TestBlacklistWithoutQuery(t *testing.T) {
	assert := assert.New(t)

	driver := &Driver{}
	server := NewServer(&Config{User: []UserConfig{}}, driver)

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver), "BlacklistRecipients",
		func(_ *Driver, senderdomain string, reasons []string, softbounce *bool, limit uint64, offset uint64) (recipients []string, err error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal("", senderdomain)
			assert.Equal([]string{}, reasons)
			assert.Equal((*bool)(nil), softbounce)
			assert.Equal(uint64(0), limit)
			assert.Equal(uint64(0), offset)

			recipients = []string{"foo@example.com"}

			return
		})

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/blacklist")
	body, status := readResponse(res)

	assert.Equal(200, status)
	assert.Equal(body, `{"recipients":["foo@example.com"]}`+"\n")
}

func TestBlacklistWithInvalidSoftbounce(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/blacklist?softbounce=x")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"strconv.ParseBool: parsing \"x\": invalid syntax"}`+"\n")
}

func TestBlacklistWithInvalidLimit(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/blacklist?limit=x")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"strconv.ParseUint: parsing \"x\": invalid syntax"}`+"\n")
}

func TestBlacklistWithInvalidOffset(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(&Config{User: []UserConfig{}}, nil)

	ts := httptest.NewServer(server.Engine)
	res, _ := http.Get(ts.URL + "/blacklist?offset=x")
	body, status := readResponse(res)

	assert.Equal(400, status)
	assert.Equal(body, `{"message":"strconv.ParseUint: parsing \"x\": invalid syntax"}`+"\n")
}
