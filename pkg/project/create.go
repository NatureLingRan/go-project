package project

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	tpl "github.com/NatureLingRan/go-project/pkg/template"
	"github.com/NatureLingRan/go-project/pkg/tools"
	"github.com/spf13/cobra"
)

// kind 创建的文件类型 如doackerfile
// name 文件的名字
// path 文件创建你的路径，当前目录为根路径
// project 项目的名字

func importPath() string {
	path := strings.ReplaceAll(os.Getenv("PWD"), os.Getenv("GOPATH")+"/src/", "")
	return strings.Replace(path, "\\", "/", -1) //将\替换成/
}

// CreateTpl 创建模板的接口
type CreateTpl interface {
	Path() string
	Content() string
}

// Create 解析模板创建文件
func Create(c CreateTpl, name, project string) {
	if c.Path() != "." {
		log.Println("创建文件夹:", c.Path())
		tools.CheckErr(os.MkdirAll(c.Path(), 0744))
	}
	log.Println("创建文件:", name)
	f, err := os.Create(filepath.Join(c.Path(), name))
	tools.CheckErr(err)

	tmpl, err := template.New("goProject").Parse(c.Content())
	tools.CheckErr(err)

	tmpl.Execute(f, map[string]string{
		"project":    project,
		"importPath": importPath(),
	})
}

// Init  如果没有指定创建文件就创建所有文件，否则就创建指定的文件
func Init(cmd *cobra.Command, args []string) {
	project := filepath.Base(os.Getenv("PWD"))

	if len(args) == 0 {
		for kind, t := range tpl.GetDefaul() {
			Create(t, kind, project)
		}
		return
	}

	kind := args[0]
	t := tpl.DefaulKind(kind)
	Create(t, kind, project)
}