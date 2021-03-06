package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/release"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Services        map[string]*Service
	VolumesToCreate []string `yaml:"volumesToCreate"`
}

type Service struct {
	name, env        string
	Nodes            map[string]string
	Image            string
	Ports            []uint16
	Command, Options []string
}

var envConfs map[string]*Conf

func Get(env string) *Conf {
	if envConfs == nil {
		if content, err := ioutil.ReadFile(filepath.Join(release.Root(), `deploy.yml`)); err != nil {
			log.Panic(err)
		} else {
			envConfs = map[string]*Conf{}
			if err = yaml.Unmarshal(content, &envConfs); err != nil {
				log.Panic(err)
			}
		}
	}
	theConf := envConfs[env]
	if theConf == nil {
		log.Fatalf(`deploy.yml: %s: undefined.`, env)
	}
	for name, svc := range theConf.Services {
		svc.name = name
		svc.env = env
	}
	return theConf
}

func HasService(svcName, env string) bool {
	_, ok := Get(env).Services[svcName]
	return ok
}

func GetService(svcName, env string) *Service {
	svc, ok := Get(env).Services[svcName]
	if !ok {
		log.Fatalf(`deploy.yml: %s.services.%s: undefined.`, env, svcName)
	}
	return svc
}

func ServiceNames(env string) (names []string) {
	services := Get(env).Services
	for _, svcName := range []string{`app`, `web`, `logc`} {
		if _, ok := services[svcName]; ok {
			names = append(names, svcName)
		}
	}
	return
}

func (svc Service) ImageName() string {
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svc.name)
	}
	return svc.Image
}

func (svc Service) ImageNameWithTag(timeTag string) string {
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svc.name)
	}
	if timeTag == `` {
		return svc.Image
	} else {
		return svc.Image + `:` + svc.env + `-` + timeTag
	}
}

func TimeTag(env string) string {
	tag := time.Now().In(release.AppConf(env).TimeLocation).Format(`060102-150405`)
	log.Println(`time tag: `, color.MagentaString(tag))
	return tag
}

var rePort = regexp.MustCompile(`^\d+$`)

func (svc Service) FirstContainerName() string {
	name := release.AppConf(svc.env).DeployName() + `-` + svc.name
	if ports := svc.Ports; len(ports) > 0 {
		name += `.` + strconv.FormatInt(int64(ports[0]), 10)
	}
	return name
}
