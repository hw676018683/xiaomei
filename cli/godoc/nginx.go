package godoc

import (
	"bytes"
	"text/template"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Nginx() error {
	if err := writeNginxConfig(); err != nil {
		return err
	}

	cmd.Run(cmd.O{Panic: true}, `sudo`, `nginx`, `-t`)
	cmd.Run(cmd.O{Panic: true}, `sudo`, `service`, `nginx`, `restart`)
	return nil
}

func writeNginxConfig() error {
	tmpl := template.Must(template.New(``).Parse(nginxConfig))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getConfData()); err != nil {
		return err
	}

	cmd.SudoWriteFile(`/etc/nginx/sites-enabled/godoc`, &buf)
	return nil
}

type confData struct {
	Deploy  *config.DeployConf
	App     *config.AppConf
	Servers *config.ServerConf
	Godoc   *config.GodocConf
	Nfs     bool
}

func getConfData() confData {
	fs, _ := cmd.Run(cmd.O{Panic: true, Output: true},
		`stat`, `--file-system`, `--format`, `%T`, config.App.Root(),
	)
	return confData{
		Deploy:  &config.Deploy,
		App:     &config.App,
		Servers: &config.Servers,
		Godoc:   &config.Godoc,
		Nfs:     fs == `nfs`,
	}
}