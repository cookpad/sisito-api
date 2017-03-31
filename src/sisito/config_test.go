package sisito

import (
	. "."
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	tml := `
[server]
log = "sisito-api.log"

[database]
host = "localhost"
port = 3306
database = "sisito"
username = "root"
password = "pass"

[[user]]
userid = "foo"
password = "bar"

[[user]]
userid = "zoo"
password = "baz"
  `

	tempFile(tml, func(f *os.File) {
		flag := &Flags{Config: f.Name()}
		config, _ := LoadConfig(flag)

		assert.Equal(*config, Config{
			Server: ServerConfig{Log: "sisito-api.log"},
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "sisito",
				Username: "root",
				Password: "pass",
			},
			User: []UserConfig{
				UserConfig{
					Userid:   "foo",
					Password: "bar",
				},
				UserConfig{
					Userid:   "zoo",
					Password: "baz",
				},
			},
		})
	})
}

func TestLoadConfigWithFilter(t *testing.T) {
	assert := assert.New(t)

	tml := `
[database]
host = "localhost"
port = 3306
database = "sisito"
username = "root"
password = "pass"

[[user]]
userid = "foo"
password = "bar"

[[user]]
userid = "zoo"
password = "baz"

[[filter]]
key = "recipient"
value = "foo@example.com"

[[filter]]
key = "senderdomain"
operator = "<>"
value = "example.net"
  `

	tempFile(tml, func(f *os.File) {
		flag := &Flags{Config: f.Name()}
		config, _ := LoadConfig(flag)

		assert.Equal(*config, Config{
			Server: ServerConfig{},
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "sisito",
				Username: "root",
				Password: "pass",
			},
			User: []UserConfig{
				UserConfig{
					Userid:   "foo",
					Password: "bar",
				},
				UserConfig{
					Userid:   "zoo",
					Password: "baz",
				},
			},

			Filter: []FilterConfig{
				FilterConfig{
					Key:      "recipient",
					Operator: "=",
					Value:    "foo@example.com",
					Join:     "AND",
					Sql:      "",
				},
				FilterConfig{
					Key:      "senderdomain",
					Operator: "<>",
					Value:    "example.net",
					Join:     "AND",
					Sql:      "",
				},
			},
		})
	})
}

func TestLoadConfigWithSql(t *testing.T) {
	assert := assert.New(t)

	tml := `
[database]
host = "localhost"
port = 3306
database = "sisito"
username = "root"
password = "pass"

[[user]]
userid = "foo"
password = "bar"

[[user]]
userid = "zoo"
password = "baz"

[[filter]]
sql = "softbounce = 0"

[[filter]]
join = "OR"
key = "senderdomain"
operator = "<>"
value = "example.net"

[[filter]]
key = "reason"
values = ["filtered", "blocked"]
  `

	tempFile(tml, func(f *os.File) {
		flag := &Flags{Config: f.Name()}
		config, _ := LoadConfig(flag)

		assert.Equal(*config, Config{
			Server: ServerConfig{},
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "sisito",
				Username: "root",
				Password: "pass",
			},
			User: []UserConfig{
				UserConfig{
					Userid:   "foo",
					Password: "bar",
				},
				UserConfig{
					Userid:   "zoo",
					Password: "baz",
				},
			},

			Filter: []FilterConfig{
				FilterConfig{
					Key:      "",
					Operator: "",
					Value:    "",
					Join:     "AND",
					Sql:      "softbounce = 0",
				},
				FilterConfig{
					Key:      "senderdomain",
					Operator: "<>",
					Value:    "example.net",
					Join:     "OR",
					Sql:      "",
				},
				FilterConfig{
					Key:      "reason",
					Operator: "IN",
					Value:    "",
					Values:   []string{"filtered", "blocked"},
					Join:     "AND",
					Sql:      "",
				},
			},
		})
	})
}
