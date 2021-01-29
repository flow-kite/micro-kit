package cmdt

import (
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/o-kit/micro-kit/dist/mconf"
	"github.com/o-kit/micro-kit/misc/context"
	"github.com/o-kit/micro-kit/misc/log"
)

var rootCmd = &cobra.Command{
	Use: filepath.Base(os.Args[0]),
}

func init() {
}

// 设置Root
func SetRoot(root *cobra.Command) {
	cmds := rootCmd.Commands()
	rootCmd = root
	mconf.ParseFlag(rootCmd.PersistentFlags())
	for _, cmd := range cmds {
		addCommandOnce(rootCmd, cmd)
	}
}

// 这里是注册服务的入口，将service 的Register进行注册
func RegisterService(s Serve) *cobra.Command {
	name := GetServiceName(s)
	cmd := registerServeEx(name, s)
	if cmd.Use == "main" {
		rootCmd.Run = cmd.Run
		return rootCmd
	}

	ServiceCmd.AddCommand(cmd)
	addCommandOnce(rootCmd, ServiceCmd)
	return cmd
}

func addCommandOnce(parent, me *cobra.Command) {
	for _, sub := range parent.Commands() {
		if sub == me {
			return
		}
	}
	parent.AddCommand(me)
}

// 获取服务名称 TODO 待确认
func GetServiceName(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return path.Base(t.PkgPath())
}

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "服务列表",
}

type Serve interface {
	Register() error         // 注册服务
	Desc() string            // 服务描述
	Run(ctx context.T) error // 启动服务
	Close() error            // 关闭服务
}

func registerServe(s Serve) error {
	// name := GetServiceName(s)
	return s.Register()
}

func registerServeEx(name string, s Serve) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: s.Desc(),
		Run: func(c *cobra.Command, args []string) {
			ctx, _ := context.WithCancel(context.Dump())

			// 调用注册服务func
			if err := registerServe(s); err != nil {
				log.Error(err)
			}

			// 这里可以捕获 sentry panic
			err := s.Run(ctx)

			s.Close()

			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			return
		},
	}

	return cmd
}
