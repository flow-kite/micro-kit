package mconf

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/o-kit/micro-kit/dist/mconsul"
	"github.com/o-kit/micro-kit/misc/build"
	"github.com/o-kit/micro-kit/misc/context"
	"github.com/o-kit/micro-kit/misc/stack"
)

/*
 * 测环境从测试环境consul上读取，生产从生产consul上读取
 *
 */

// 这里读取配置文件 - 需要从consul上读取
func Render(data []byte) ([]byte, error) {
	return RenderWithService(data, "")
}

func RenderWithService(data []byte, name string) ([]byte, error) {
	dcs, err := mconsul.GetDatacenters()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 取项目路径
	if name == "" {
		name = stack.GetRootService()
	}
	tmpl, err := template.New("mconf").Funcs(template.FuncMap{
		"key":     consulGetKey(dcs, ""),
		"db":      consulGetKey(dcs, "myapp/database/"),
		"service": consulGetKey(dcs, "service/"+name+"/"),
		"prefix":  getKeyPrefix(dcs, ""),
		"dbs":     getKeyPrefix(dcs, "myapp/database/"),
		"test":    getConfig(dcs, name),
	}).Parse(string(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, "asdf"); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func getConfig(dcs []string, serviceName string) func(name string) string {
	return func(name string) string {
		prefix := "/service/" + serviceName + "/"
		var pairs []mconsul.Pair
		for _, dc := range dcs {
			pairs, _ = mconsul.NewDatacenter(dc).GetPairs(prefix + name)
			if len(pairs) == 0 {
				continue
			}
			break
		}
		if len(pairs) == 0 {
			panic(fmt.Sprintf("keyPath %q is not found", prefix+name))
		}

		buf := bytes.NewBuffer(nil)
		buf.WriteString("[" + name + "]\n")
		for _, pair := range pairs {
			key := path.Base(pair.Key)
			buf.WriteString(key + " = " + string(pair.Value) + "\n")
		}
		return buf.String()
	}
}

func getKeyPrefix(dcs []string, prefix string) func(prefix2 string) func(key string) string {
	return func(prefix2 string) func(string) string {
		return func(key string) string {
			return consulGetKey(dcs, prefix)(prefix2, key)
		}
	}
}

func consulGetKey(dcs []string, prefix string) func(keys ...string) string {
	return func(keys ...string) string {
		key := strings.Join(keys, "_")
		context.LogInfof("key = %v", key)
		for _, dc := range dcs {
			val, err := mconsul.NewDatacenter(dc).GetValue(prefix + key)
			if err != nil || len(val) == 0 {
				continue
			}
			return string(val)
		}
		panic("key " + prefix + key + "is not found")
	}
}

func ReadFile(fp string, obj interface{}) error {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return errors.Wrap(err, "readFile failed")
	}
	// 利用consul对config.toml上的参数进行渲染
	data, err = Render(data)
	if err != nil {
		return errors.WithStack(err)
	}
	viper.SetConfigType("toml")
	_ = viper.ReadConfig(bytes.NewBuffer(data))
	return viper.Unmarshal(obj)
}

var defaultConf = "test/conf.toml"

var byFlagOnce sync.Once

func ByFlagOnce(obj interface{}) {
	byFlagOnce.Do(func() {
		if err := ReadFile(GetConfPath(), obj); err != nil {
			panic(err)
		}
	})
}

func GetConfPath() string {
	if !build.IsTest {
		dir, _ := os.Getwd()
		fmt.Println("dir = ", dir)
		sp := strings.Split(dir, "/")
		for idx, item := range sp {
			if item == "src" {
				dir = strings.Join(sp[:idx+3], "/")
				break
			}
		}
		return filepath.Join(dir, defaultConf)
	}

	return defaultConf
}

// 判断是否有配置文件
func ParseFlag(p *pflag.FlagSet) {
	hasConf := false
	p.VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "test" {
			hasConf = true
		}
		if flag.Shorthand == "c" {
			hasConf = true
		}
	})

	if !hasConf {
		p.StringVarP(&defaultConf, "test", "c", defaultConf, "toml path")
	}
}
