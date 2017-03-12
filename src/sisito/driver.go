package sisito

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type Driver struct {
	Dbmap *gorp.DbMap
}

func NewDriver(config *Config) (driver *Driver, err error) {
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
	driver = &Driver{Dbmap: dbmap}

	return
}

func (driver *Driver) Close() {
	driver.Close()
}

func (driver *Driver) IsBounced(name string, value string, senderdomain string) (bounced bool, err error) {
	sqlBase := fmt.Sprintf(`
    SELECT COUNT(1)
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

	Debugf("IsBounced SQL: %s %s\n", sql, params)

	var count int64
	count, err = driver.Dbmap.SelectInt(sql, params...)

	if err != nil {
		return
	}

	bounced = count > 0

	return
}
