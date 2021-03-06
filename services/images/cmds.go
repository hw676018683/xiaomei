package images

import (
	"github.com/lovego/xiaomei/services/deploy/conf"
	"github.com/lovego/xiaomei/services/images/registry"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds(svcName string) []*cobra.Command {
	return []*cobra.Command{
		buildCmdFor(svcName),
		pushCmdFor(svcName),
		tagsCmdFor(svcName),
	}
}

func buildCmdFor(svcName string) *cobra.Command {
	var tag, pull bool
	cmd := &cobra.Command{
		Use:   `build [<env>]`,
		Short: `[image] build  ` + imageDesc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			timeTag := ``
			if tag {
				timeTag = conf.TimeTag(env)
			}
			return Build(svcName, env, timeTag, pull)
		}),
	}
	cmd.Flags().BoolVarP(&tag, `tag`, `t`, false, `add a deploy time tag.`)
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	return cmd
}

func pushCmdFor(svcName string) *cobra.Command {
	return &cobra.Command{
		Use:   `push [<env> [<tag>]]`,
		Short: `[image] push   ` + imageDesc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			return Push(svcName, env, timeTag)
		}),
	}
}

func tagsCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `tags [<env>]`,
		Short: `[image] list time tags of ` + imageDesc(svcName) + ` in registry.`,
		RunE: release.EnvCall(func(env string) error {
			registry.ListTimeTags(svcName, env)
			return nil
		}),
	}
	cmd.AddCommand(pruneCmdFor(svcName), rmCmdFor(svcName), digestCmdFor(svcName))
	return cmd
}

func pruneCmdFor(svcName string) *cobra.Command {
	var n uint8
	cmd := &cobra.Command{
		Use:   `prune [<env>]`,
		Short: `prune registry time tags of ` + imageDesc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			registry.PruneTimeTags(svcName, env, int(n))
			return nil
		}),
	}
	cmd.Flags().Uint8VarP(&n, `number`, `n`, 10, `the number of time tags to keep.`)
	return cmd
}

func rmCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `rm [<env> [<tag>....]]`,
		Short: `remove registry time tags of ` + imageDesc(svcName) + `.`,
		RunE: release.EnvSliceCall(func(env string, tags []string) error {
			registry.RemoveTimeTags(svcName, env, tags)
			return nil
		}),
	}
	return cmd
}

func digestCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `digest [<env> [<tag>....]]`,
		Short: `get registry digest of ` + imageDesc(svcName) + `.`,
		RunE: release.EnvSliceCall(func(env string, tags []string) error {
			registry.DigestTimeTags(svcName, env, tags)
			return nil
		}),
	}
	return cmd
}

func imageDesc(svcName string) string {
	if svcName == `` {
		return `all images`
	} else {
		return `the ` + svcName + ` image`
	}
}
