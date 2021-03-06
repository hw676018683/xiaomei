package token

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lovego/sessions/cookiestore"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func tokenGenCmd() *cobra.Command {
	var domain string
	var values = make(map[string]string)
	cmd := &cobra.Command{
		Use:   `gen <env>`,
		Short: `generate token command.`,
		RunE: release.EnvCall(func(env string) error {
			ck := newCookie(release.AppConf(env).Cookie)
			return tokenGen(ck, release.AppConf(env).Secret, domain, values)
		}),
	}
	cmd.Flags().StringVarP(&domain, `domain`, ``, ``, `specified the cookie domain.`)
	cmd.Flags().StringToStringVar(&values, `value`, nil, `specified the value of domain.`)
	return cmd
}

func tokenGen(ck *http.Cookie, secret, domain string, values map[string]string) error {
	data := make(map[string]interface{})
	for key, value := range values {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			data[key] = v
		} else {
			data[key] = value
		}
	}
	if domain != `` {
		ck.Domain = domain
	}
	sess, err := cookiestore.New(secret).EncodeData(ck.Name, data)
	if err != nil {
		return err
	}

	ck.Value = string(sess)

	fmt.Println(fmt.Sprintf(`document.cookie = "%s";`, ck.String()))
	return nil
}
