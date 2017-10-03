package workspace_godoc

import "github.com/lovego/xiaomei/utils/cmd"

func shell() error {
	_, err := cmd.Run(cmd.O{}, `docker`, `exec`, `-it`,
		`--detach-keys=ctrl-@`, `workspace-godoc`, `bash`,
	)
	return err
}
