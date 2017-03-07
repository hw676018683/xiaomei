package config

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/utils/mailer"
)

var Config Conf

type Conf struct {
	sync.Mutex

	env      string
	envs     []string
	data     *conf
	timeZone *time.Location
	mailer   struct {
		setted bool
		*mailer.Mailer
	}
}

type conf struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Secret string `yaml:"secret"`

	TimeZone TimeZoneConf    `yaml:"timeZone"`
	Mailer   MailerConf      `yaml:"mailer"`
	Keepers  []mailer.People `yaml:"keepers"`

	DataSource map[string]map[string]string `yaml:"dataSource"`
}

type TimeZoneConf struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

type MailerConf struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Sender mailer.People
	Passwd string `yaml:"passwd"`
}

func (c *Conf) Root() string {
	if root := DetectRoot(); root != `` {
		return root
	} else {
		panic(`app root not found.`)
	}
}

func (c *Conf) Name() string {
	c.Load()
	return c.data.Name
}

func (c *Conf) Env() string {
	if c.env == `` {
		c.env = detectEnv()
	}
	return c.env
}

func (c *Conf) Envs() []string {
	if c.envs == nil {
		c.envs = availableEnvs()
	}
	return c.envs
}

func (c *Conf) DeployName() string {
	return c.Name() + `_` + c.Env()
}

func (c *Conf) Bin() string {
	return filepath.Join(c.Root(), c.Name())
}

func (c *Conf) Domain() string {
	c.Load()
	return c.data.Domain
}

func (c *Conf) Secret() string {
	c.Load()
	return c.data.Secret
}

func (c *Conf) TimeZone() *time.Location {
	if c.timeZone == nil {
		c.Load()
		c.timeZone = time.FixedZone(c.data.TimeZone.Name, c.data.TimeZone.Offset)
	}
	return c.timeZone
}

func (c *Conf) Mailer() *mailer.Mailer {
	c.Lock()
	defer c.Unlock()
	if !c.mailer.setted {
		c.Load()
		m := c.data.Mailer
		c.mailer.Mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
		c.mailer.setted = true
	}
	return c.mailer.Mailer
}

func (c *Conf) Alarm(title, body string) {
	title = c.DeployName() + ` ` + title
	c.Mailer().Send(&mailer.Message{Receivers: c.Keepers(), Title: title, Body: body})
}

func (c *Conf) Keepers() []mailer.People {
	c.Load()
	return c.data.Keepers
}

func (c *Conf) DataSource(typ, key string) string {
	c.Load()
	if key == `` {
		key = `default`
	}
	return c.data.DataSource[typ][key]
}
