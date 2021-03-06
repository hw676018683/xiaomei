package deploy

import (
	"github.com/lovego/xiaomei/services/deploy/conf"
	"github.com/lovego/xiaomei/services/images"
	"github.com/lovego/xiaomei/services/images/registry"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

// Run, Deploy, Ps, Logs commands
func Cmds(svcName string) (cmds []*cobra.Command) {
	return []*cobra.Command{
		deployCmdFor(svcName),
		rmDeployCmdFor(svcName),
	}
}

func deployCmdFor(svcName string) *cobra.Command {
	var filter string
	var pull, push bool
	cmd := &cobra.Command{
		Use:   `deploy [<env> [<tag>]]`,
		Short: `deploy the ` + desc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			noTag := timeTag == ``
			if timeTag == `` {
				timeTag = conf.TimeTag(env)
				if err := images.Build(svcName, env, timeTag, pull); err != nil {
					return err
				}
				if push {
					if err := images.Push(svcName, env, timeTag); err != nil {
						return err
					}
				}
			}
			if err := deploy(svcName, env, timeTag, filter); err != nil {
				return err
			}
			if noTag {
				registry.PruneTimeTags(svcName, env, 10)
			}
			return nil
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	cmd.Flags().BoolVarP(&push, `push`, `P`, true, `push the built images to registry.`)
	return cmd
}

func rmDeployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `rm-deploy [<env>]`,
		Short: `remove deployment of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return rmDeploy(svcName, env, filter)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func desc(svcName string) string {
	if svcName == `` {
		return `project`
	} else {
		return svcName + ` service`
	}
}
