package sisito

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"os"
	"time"
)

type Driver struct {
	Dbmap *gorp.DbMap
}

func NewDriver(config *Config, debug bool) (driver *Driver, err error) {
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
		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "", log.Lmicroseconds))
	}

	driver = &Driver{Dbmap: dbmap}

	return
}

func (driver *Driver) Close() {
	driver.Close()
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
	Recipient      string
	Digest         string
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	Whitelisted    uint8
}

type Recipient struct {
	Recipient string
}

func (driver *Driver) recentlyBounced(name string, value string, senderdomain string) (bounced []BounceMail, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT bm.*, IF(wm.id IS NULL, 0, 1) AS whitelisted
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.%s = ?`, name)

	sqlBuf := bytes.NewBufferString(sqlBase)
	params := make([]interface{}, 1)
	params[0] = value

	if senderdomain != "" {
		sqlBuf.WriteString(`
       AND bm.senderdomain = ?`)

		params = append(params, senderdomain)
	}

	sqlBuf.WriteString(`
  ORDER BY bm.id DESC
     LIMIT 1`)

	sql := sqlBuf.String()

	bounced = []BounceMail{}
	_, err = driver.Dbmap.Select(&bounced, sql, params...)

	return
}

func (driver *Driver) isBounced(name string, value string, senderdomain string) (bounced bool, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT 1
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE bm.%s = ?`, name)

	sqlBuf := bytes.NewBufferString(sqlBase)
	params := make([]interface{}, 1)
	params[0] = value

	if senderdomain != "" {
		sqlBuf.WriteString(`
       AND bm.senderdomain = ?`)

		params = append(params, senderdomain)
	}

	sqlBuf.WriteString(`
       AND wm.id IS NULL
     LIMIT 1`)

	sql := sqlBuf.String()

	var count int64
	count, err = driver.Dbmap.SelectInt(sql, params...)

	if err != nil {
		return
	}

	if count > 0 {
		bounced = true
	} else {
		bounced = false
	}

	return
}

func (driver *Driver) blacklistRecipients(senderdomain string) (recipients []string, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT bm.recipient
      FROM bounce_mails bm LEFT JOIN whitelist_mails wm
        ON bm.recipient = wm.recipient AND bm.senderdomain = wm.senderdomain
     WHERE wm.id IS NULL`)

	sqlBuf := bytes.NewBufferString(sqlBase)
	params := make([]interface{}, 0)

	if senderdomain != "" {
		sqlBuf.WriteString(`
       AND bm.senderdomain = ?`)

		params = append(params, senderdomain)
	}

	sqlBuf.WriteString(`
  GROUP BY recipient`)

	sql := sqlBuf.String()

	recipients = []string{}
	_, err = driver.Dbmap.Select(&recipients, sql, params...)

	if err != nil {
		return
	}

	return
}
