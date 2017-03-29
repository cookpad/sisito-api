package sisito

import (
	. "."
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"
	"reflect"
	"testing"
)

func TestRecentlyListed(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyListed("recipient", "foo@example.com", "example.net", true)

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestRecentlyListedWithFilter(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Key: "recipient", Operator: "NOT LIKE", Value: "localhost.localdomain"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND bm.recipient NOT LIKE ?
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net", "localhost.localdomain"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyListed("recipient", "foo@example.com", "example.net", true)

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestRecentlyListedWithoutFilter(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Key: "recipient", Operator: "NOT LIKE", Value: "localhost.localdomain"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyListed("recipient", "foo@example.com", "example.net", false)

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestRecentlyListedWithSql(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Sql: "recipient NOT LIKE 'localhost.localdomain'"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND recipient NOT LIKE 'localhost.localdomain'
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyListed("recipient", "foo@example.com", "example.net", true)

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestRecentlyListedWithoutSenderdomain(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
  ORDER BY bm.id DESC
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com"}, args)

			rows := i.(*[]BounceMail)
			*rows = append(*rows, BounceMail{Id: 1})

			return nil, nil
		})

	rows, _ := driver.RecentlyListed("recipient", "foo@example.com", "", true)

	assert.Equal([]BounceMail{BounceMail{Id: 1}}, rows)
}

func TestListed(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			return 1, nil
		})

	count, _ := driver.Listed("recipient", "foo@example.com", "example.net", true)

	assert.Equal(count, true)
}

func TestListedWithFilter(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Key: "recipient", Operator: "NOT LIKE", Value: "localhost.localdomain"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND bm.recipient NOT LIKE ?
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net", "localhost.localdomain"}, args)

			return 1, nil
		})

	count, _ := driver.Listed("recipient", "foo@example.com", "example.net", true)

	assert.Equal(count, true)
}

func TestListedWithoutFilter(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Key: "recipient", Operator: "NOT LIKE", Value: "localhost.localdomain"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			return 1, nil
		})

	count, _ := driver.Listed("recipient", "foo@example.com", "example.net", false)

	assert.Equal(count, true)
}

func TestListedWithSql(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Sql: "recipient NOT LIKE 'localhost.localdomain'"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND bm.senderdomain = ?
       AND recipient NOT LIKE 'localhost.localdomain'
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com", "example.net"}, args)

			return 1, nil
		})

	count, _ := driver.Listed("recipient", "foo@example.com", "example.net", true)

	assert.Equal(count, true)
}

func TestListedWithoutSenderdomain(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "SelectInt",
		func(_ *gorp.DbMap, query string, args ...interface{}) (int64, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.recipient = ?
       AND wm.id IS NULL
     LIMIT 1`, query)

			assert.Equal([]interface{}{"foo@example.com"}, args)

			return 1, nil
		})

	count, _ := driver.Listed("recipient", "foo@example.com", "", true)

	assert.Equal(count, true)
}

func TestBlacklistRecipients(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.recipient
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL
       AND bm.senderdomain = ?
       AND bm.reason IN (?,?)
       AND bm.softbounce = ?
  GROUP BY bm.recipient
  ORDER BY bm.recipient
     LIMIT ?
    OFFSET ?`, query)

			assert.Equal([]interface{}{
				"example.net", "userunknown", "filtered", false, uint64(100), uint64(100)}, args)

			rows := i.(*[]string)
			*rows = append(*rows, "foo@example.com")

			return nil, nil
		})

	recipients, _ := driver.BlacklistRecipients(
		"example.net", []string{"userunknown", "filtered"}, new(bool), 100, 100, true)

	assert.Equal(recipients, []string{"foo@example.com"})
}

func TestBlacklistRecipientsWithFilter(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Key: "recipient", Operator: "NOT LIKE", Value: "localhost.localdomain"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.recipient
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL
       AND bm.senderdomain = ?
       AND bm.reason IN (?,?)
       AND bm.softbounce = ?
       AND bm.recipient NOT LIKE ?
  GROUP BY bm.recipient
  ORDER BY bm.recipient
     LIMIT ?
    OFFSET ?`, query)

			assert.Equal([]interface{}{
				"example.net", "userunknown", "filtered", false, "localhost.localdomain", uint64(100), uint64(100)}, args)

			rows := i.(*[]string)
			*rows = append(*rows, "foo@example.com")

			return nil, nil
		})

	recipients, _ := driver.BlacklistRecipients(
		"example.net", []string{"userunknown", "filtered"}, new(bool), 100, 100, true)

	assert.Equal(recipients, []string{"foo@example.com"})
}

func TestBlacklistRecipientsWithoutFilter(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Key: "recipient", Operator: "NOT LIKE", Value: "localhost.localdomain"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.recipient
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL
       AND bm.senderdomain = ?
       AND bm.reason IN (?,?)
       AND bm.softbounce = ?
  GROUP BY bm.recipient
  ORDER BY bm.recipient
     LIMIT ?
    OFFSET ?`, query)

			assert.Equal([]interface{}{
				"example.net", "userunknown", "filtered", false, uint64(100), uint64(100)}, args)

			rows := i.(*[]string)
			*rows = append(*rows, "foo@example.com")

			return nil, nil
		})

	recipients, _ := driver.BlacklistRecipients(
		"example.net", []string{"userunknown", "filtered"}, new(bool), 100, 100, false)

	assert.Equal(recipients, []string{"foo@example.com"})
}

func TestBlacklistRecipientsWithSql(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{
		Filter: []FilterConfig{
			FilterConfig{Sql: "recipient NOT LIKE 'localhost.localdomain'"},
		},
	}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.recipient
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL
       AND bm.senderdomain = ?
       AND bm.reason IN (?,?)
       AND bm.softbounce = ?
       AND recipient NOT LIKE 'localhost.localdomain'
  GROUP BY bm.recipient
  ORDER BY bm.recipient
     LIMIT ?
    OFFSET ?`, query)

			assert.Equal([]interface{}{
				"example.net", "userunknown", "filtered", false, uint64(100), uint64(100)}, args)

			rows := i.(*[]string)
			*rows = append(*rows, "foo@example.com")

			return nil, nil
		})

	recipients, _ := driver.BlacklistRecipients(
		"example.net", []string{"userunknown", "filtered"}, new(bool), 100, 100, true)

	assert.Equal(recipients, []string{"foo@example.com"})
}

func TestBlacklistRecipientsWithoutOptions(t *testing.T) {
	assert := assert.New(t)
	driver := &Driver{Config: &Config{}, DbMap: &gorp.DbMap{}}

	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(driver.DbMap), "Select",
		func(_ *gorp.DbMap, i interface{}, query string, args ...interface{}) ([]interface{}, error) {
			defer guard.Unpatch()
			guard.Restore()

			assert.Equal(`
    SELECT bm.recipient
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL
  GROUP BY bm.recipient
  ORDER BY bm.recipient`, query)

			assert.Equal([]interface{}{}, args)

			rows := i.(*[]string)
			*rows = append(*rows, "foo@example.com")

			return nil, nil
		})

	recipients, _ := driver.BlacklistRecipients("", nil, nil, 0, 0, true)

	assert.Equal(recipients, []string{"foo@example.com"})
}
