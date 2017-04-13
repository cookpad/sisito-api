package sisito

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"io"
	"log"
	"strings"
	"time"
)

type Driver struct {
	Config *Config
	DbMap  *gorp.DbMap
}

func NewDriver(config *Config, debug bool, out io.Writer) (driver *Driver, err error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database)

	var db *sql.DB
	db, err = sql.Open("mysql", url)

	if err != nil {
		return
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

	if debug {
		dbmap.TraceOn("[gorp]", log.New(out, "", log.Ldate|log.Ltime))
	}

	driver = &Driver{Config: config, DbMap: dbmap}

	return
}

func (driver *Driver) Close() {
	driver.DbMap.Db.Close()
}

type BounceMail struct {
	Id             int32
	Timestamp      time.Time
	Lhost          string
	Rhost          string
	Alias          string
	Listid         string
	Reason         string
	Action         string
	Subject        string
	Messageid      string
	Smtpagent      string
	Softbounce     uint8
	Smtpcommand    string
	Destination    string
	Senderdomain   string
	Feedbacktype   string
	Diagnosticcode string
	Deliverystatus string
	Timezoneoffset string
	Addresser      string
	Addresseralias string
	Recipient      string
	Digest         string
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	Whitelisted    uint8
}

type Recipient struct {
	Recipient string
	Alias     string
}

func (driver *Driver) appendFilter(buf *bytes.Buffer, params *[]interface{}) {
	for _, filter := range driver.Config.Filter {
		if filter.Sql != "" {
			buf.WriteString("\n       ")
			buf.WriteString(filter.Join)
			buf.WriteString(" ")
			buf.WriteString(filter.Sql)
		} else {
			buf.WriteString("\n       ")
			buf.WriteString(filter.Join)
			buf.WriteString(" bm.")
			buf.WriteString(filter.Key)
			buf.WriteString(" ")
			buf.WriteString(filter.Operator)

			if filter.Values == nil {
				buf.WriteString(" ?")
				*params = append(*params, filter.Value)
			} else {
				buf.WriteString(" (")

				for i, value := range filter.Values {
					if i > 0 {
						buf.WriteString(",")
					}

					buf.WriteString("?")
					*params = append(*params, value)
				}

				buf.WriteString(")")
			}
		}
	}
}

func (driver *Driver) RecentlyListed(name string, value string, senderdomain string, useFilter bool) (listed []BounceMail, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.%s = ?`, name)

	if name == "recipient" {
		value = NormalizeRecipient(value)
	}

	sqlBuf := bytes.NewBufferString(sqlBase)
	params := []interface{}{value}

	if senderdomain != "" {
		sqlBuf.WriteString(`
       AND bm.senderdomain = ?`)

		params = append(params, senderdomain)
	}

	if useFilter {
		driver.appendFilter(sqlBuf, &params)
	}

	sqlBuf.WriteString(`
  ORDER BY bm.id DESC
     LIMIT 1`)

	sql := sqlBuf.String()

	listed = []BounceMail{}
	_, err = driver.DbMap.Select(&listed, sql, params...)

	return
}

func (driver *Driver) Listed(name string, value string, senderdomain string, useFilter bool) (listed bool, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.%s = ?`, name)

	if name == "recipient" {
		value = NormalizeRecipient(value)
	}

	sqlBuf := bytes.NewBufferString(sqlBase)
	params := []interface{}{value}

	if senderdomain != "" {
		sqlBuf.WriteString(`
       AND bm.senderdomain = ?`)

		params = append(params, senderdomain)
	}

	if useFilter {
		driver.appendFilter(sqlBuf, &params)
	}

	sqlBuf.WriteString(`
       AND wm.id IS NULL
     LIMIT 1`)

	sql := sqlBuf.String()

	var count int64
	count, err = driver.DbMap.SelectInt(sql, params...)

	if err != nil {
		return
	}

	if count > 0 {
		listed = true
	} else {
		listed = false
	}

	return
}

func (driver *Driver) BlacklistRecipients(senderdomain string, reasons []string, softbounce *bool, limit uint64, offset uint64, useFilter bool) (recipients []string, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT bm.recipient, bm.alias
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL`)

	sqlBuf := bytes.NewBufferString(sqlBase)
	params := []interface{}{}

	if senderdomain != "" {
		sqlBuf.WriteString(`
       AND bm.senderdomain = ?`)

		params = append(params, senderdomain)
	}

	if len(reasons) > 0 {
		sqlBuf.WriteString(`
       AND bm.reason IN (`)

		phs := make([]string, len(reasons))

		for i, v := range reasons {
			params = append(params, v)
			phs[i] = "?"
		}

		sqlBuf.WriteString(strings.Join(phs, ","))
		sqlBuf.WriteString(")")
	}

	if softbounce != nil {
		sqlBuf.WriteString(`
       AND bm.softbounce = ?`)

		params = append(params, *softbounce)
	}

	if useFilter {
		driver.appendFilter(sqlBuf, &params)
	}

	sqlBuf.WriteString(`
  GROUP BY bm.recipient
  ORDER BY bm.recipient`)

	if limit > 0 {
		sqlBuf.WriteString(`
     LIMIT ?`)

		params = append(params, limit)
	}

	if offset > 0 {
		sqlBuf.WriteString(`
    OFFSET ?`)

		params = append(params, offset)
	}

	sql := sqlBuf.String()

	recipientAlieses := []Recipient{}
	_, err = driver.DbMap.Select(&recipientAlieses, sql, params...)

	if err != nil {
		return
	}

	recipients = MergeRecipientAliases(recipientAlieses)

	return
}

func NormalizeRecipient(recipient string) (normalized string) {
	normalized = strings.Replace(recipient, `"`, "", -1)
	normalized = strings.Replace(normalized, `'`, "", -1)
	return
}

func MergeRecipientAliases(recipientAlieses []Recipient) (merged []string) {
	mergedMap := make(map[string]bool, len(recipientAlieses))

	for _, v := range recipientAlieses {
		mergedMap[v.Recipient] = true
		mergedMap[v.Alias] = true
	}

	merged = make([]string, len(mergedMap))
	i := 0

	for k, _ := range mergedMap {
		merged[i] = k
		i++
	}

	return
}
