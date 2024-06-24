package main

import (
	"bytes"
	"cmp"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/imports"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

var (
	delegateMethodName = "GetDelegate"
	wrapperMethodName  = "WrapperMethod"
)
var (
	inputFile = flag.String("input", "", "input file")
	typeName  = flag.String("type", "", "type")
)

func main() {
	//flag.Usage()
	flag.Parse()
	fmt.Printf("os.Args:%+v\n", os.Args)
	fmt.Printf("%+v\n", flag.Parsed())
	fmt.Printf("inputFile:%s, typeName:%s\n", *inputFile, *typeName)
	*typeName = "WrapperTest"
	file := "wrapper.go"
	dir := Must(os.Getwd())
	input := path.Join(dir, file)
	parsedFile := Must(parser.ParseFile(token.NewFileSet(), input, nil, parser.Mode(0)))

	Must(parser.ParseFile(token.NewFileSet(), path.Join(dir, "origin.go"), nil, parser.Mode(0)))
	fmt.Printf("name:%+v\n", parsedFile.Name)
	visitor := &visitor{
		typeName: *typeName,
		baseDir:  dir,
		parsedInfo: &GenData{
			PackageName:      parsedFile.Name.Name,
			TypeName:         *typeName,
			Interfaces:       make([]*InterfaceInfo, 0),
			Imports:          make([]*ImportInfo, 0),
			VendorInterfaces: make(map[string]*InterfaceInfo),
		},
	}
	ast.Walk(visitor, parsedFile)
	parsedInfo := visitor.parsedInfo
	slices.SortFunc(parsedInfo.Interfaces, func(i, j *InterfaceInfo) int {
		return -cmp.Compare(i.Name, j.Name)
	})
	pwd, _ := os.Getwd()
	parse := template.Must(template.ParseFiles(path.Join(pwd, "../cmd/delegate.tpl")))

	source := bytes.NewBufferString("")
	err := parse.Execute(source, parsedInfo)
	if err != nil {
		fmt.Printf("execute template error:%v", err)
	}

	sourceBytes := source.Bytes()
	//// gofmt source
	sourceBytes, err = imports.Process("", sourceBytes, nil)
	if err != nil {
		fmt.Printf("imports.Process error:%v", err)
	}
	output := path.Join(dir, fmt.Sprintf("%s_gen.go", removeExtension(file)))
	err = os.WriteFile(output, sourceBytes, 0644)
	if err != nil {
		fmt.Printf("write file error:%v", err)
	}
}

var _ ast.Visitor = (*visitor)(nil)

type visitor struct {
	typeName    string
	baseDir     string
	FuncTypeMap map[string]string
	parsedInfo  *GenData
}
type GenData struct {
	PackageName      string
	TypeName         string
	Wrapper          *ParamInfo
	Imports          []*ImportInfo
	Interfaces       []*InterfaceInfo
	FuncTypes        []*FuncType
	VendorInterfaces map[string]*InterfaceInfo
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return nil
	}
	switch t := node.(type) {
	case *ast.ImportSpec:
		v.parsedInfo.Imports = append(v.parsedInfo.Imports, &ImportInfo{
			Name: func() string {
				if t.Name != nil {
					return t.Name.Name
				}
				return ""
			}(),
			Path: t.Path.Value,
		})
	case *ast.TypeSpec:
		switch typeSpec := t.Type.(type) {
		case *ast.InterfaceType:
			log.Printf("interface --> %s", t.Name.Name)
			v.parseInterface("", t.Name.Name, typeSpec)
		}
	case *ast.InterfaceType:

	}
	return v
}

func (v *visitor) parseInterface(ns string, interfaceName string, t *ast.InterfaceType) {
	inter := &InterfaceInfo{
		Name:    interfaceName,
		Methods: make([]*Method, 0),
	}
	v.parsedInfo.Interfaces = append(v.parsedInfo.Interfaces, inter)

	for _, method := range t.Methods.List {
		fmt.Printf("method:%+v\n", method)
		handleMethod := func(methodType *ast.FuncType) {
			args := make([]*ParamInfo, 0)
			if methodType.Params != nil {
				for _, field := range methodType.Params.List {
					args = append(args, v.buildParams(ns, false, field.Names, field.Type)...)
				}
			}

			returns := make([]*ParamInfo, 0)
			if methodType.Results != nil {
				for _, field := range methodType.Results.List {
					returns = append(returns, v.buildParams(ns, false, field.Names, field.Type)...)
				}
			}

			var ctxParam *ParamInfo
			for _, arg := range args {
				if arg.Typed == "context.Context" {
					ctxParam = arg
				}
			}
			inter.Methods = append(inter.Methods, &Method{
				Name:    method.Names[0].Name,
				Args:    args,
				Returns: returns,
				CtxName: func() string {
					if ctxParam == nil {
						return ""
					}
					return ctxParam.Name
				}(),
			})
		}
		switch mt := method.Type.(type) {
		case *ast.FuncType:
			handleMethod(mt)
		case *ast.Ident:
			inter := v.loadVendorAndGetInterface("", mt.Name)
			v.parsedInfo.Interfaces = append(v.parsedInfo.Interfaces, inter)
		case *ast.SelectorExpr:
			fmt.Printf("selectorExpr:%+v\n", mt)
			inter := v.loadVendorAndGetInterface(mt.X.(*ast.Ident).Name, mt.Sel.Name)
			v.parsedInfo.Interfaces = append(v.parsedInfo.Interfaces, inter)
		default:
			fmt.Printf("method :%+v\n", method)
		}
	}
}

func (v *visitor) buildParams(ns string, ptr bool, names []*ast.Ident, expr ast.Expr) []*ParamInfo {
	params := make([]*ParamInfo, 0)
	switch t := expr.(type) {
	case *ast.StarExpr:
		params = append(params, v.buildParams(ns, true, names, t.X)...)
	case *ast.Ident:
		if len(names) == 0 {
			names = append(names, &ast.Ident{Name: genVarName(t.Name)})
		}
		for _, name := range names {
			params = append(params, &ParamInfo{
				Name: name.Name,
				Typed: func() string {
					if ns == "" || !t.IsExported() {
						return t.Name
					}
					return ns + "." + t.Name
				}(),
				Ptr: ptr,
			})
		}
	case *ast.SelectorExpr:
		if len(names) == 0 {
			params = append(params, &ParamInfo{
				Name:  genVarName(t.Sel.Name),
				Typed: fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name),
				Ptr:   ptr,
			})
		}
		for _, name := range names {
			params = append(params, &ParamInfo{
				Name:  name.Name,
				Typed: fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name),
				Ptr:   ptr,
			})
		}
	case *ast.InterfaceType:
		if len(names) == 0 {
			names = append(names, &ast.Ident{Name: "_inter"})
		}
		for _, name := range names {
			params = append(params, &ParamInfo{
				Name:  name.Name,
				Typed: "interface{}",
				Ptr:   ptr,
			})
		}
	case *ast.Ellipsis:
		typed, ptr := getTyped(ns, t.Elt)
		if len(names) == 0 {
			names = append(names, &ast.Ident{Name: "_arg1"})
		}
		params = append(params, &ParamInfo{
			Name:     names[0].Name,
			Typed:    typed,
			Ptr:      ptr,
			Ellipsis: true,
		})
	case *ast.FuncType:
		log.Printf("FuncType params :%+v, %+v", t, names)
		if len(names) == 0 {
			names = append(names, &ast.Ident{Name: "_fn1"})
		}
		_args := make([]*ParamInfo, 0)
		for _, _param := range t.Params.List {
			log.Printf("_param :%+v", _param)
			_args = append(_args, v.buildParams(ns, ptr, _param.Names, _param.Type)...)
		}
		_returns := make([]*ParamInfo, 0)
		for _, _param := range t.Results.List {
			_returns = append(_returns, v.buildParams(ns, ptr, _param.Names, _param.Type)...)
		}

		toParamType := func(args []*ParamInfo) []string {
			_argTypes := make([]string, 0)
			for _, arg := range args {
				_argTypes = append(_argTypes, func() string {
					builder := &strings.Builder{}
					if arg.Ellipsis {
						builder.WriteString("...")
					}
					if arg.Ptr {
						builder.WriteString("*")
					}
					builder.WriteString(arg.Typed)
					return builder.String()
				}())
			}
			return _argTypes
		}

		funcInfo := "func(" + strings.Join(toParamType(_args), ",") + ")(" + strings.Join(toParamType(_returns), ",") + ")"
		funcTypeName, ok := v.FuncTypeMap[funcInfo]

		if !ok {
			funcType := &FuncType{
				Name: "FuncType" + strconv.Itoa(VarIndex()),
				Type: "func(" + strings.Join(toParamType(_args), ",") + ")(" + strings.Join(toParamType(_returns), ",") + ")",
			}
			v.parsedInfo.FuncTypes = append(v.parsedInfo.FuncTypes, funcType)
			v.FuncTypeMap[funcInfo] = funcType.Name
			funcTypeName = funcType.Name
		}

		for _, name := range names {
			params = append(params, &ParamInfo{
				Name:  name.Name,
				Typed: funcTypeName,
				Ptr:   ptr,
			})
		}
	case *ast.ArrayType:
		typed, _ptr := getTyped(ns, t.Elt)
		if len(names) == 0 {
			names = append(names, &ast.Ident{Name: genVarName(typed)})
		}
		for _, name := range names {

			params = append(params, &ParamInfo{
				Name:  name.Name,
				Typed: "[]" + typed,
				Ptr:   _ptr,
			})
		}
	case *ast.MapType:
		log.Printf("MapType params :%+v, %+v", t, names)
		kTyped, kPtr := getTyped(ns, t.Key)
		vTyped, vPtr := getTyped(ns, t.Value)
		if len(names) == 0 {
			names = append(names, &ast.Ident{Name: "_map"})
		}
		for _, name := range names {
			named := func(n string, _p bool) string {
				if _p {
					return "*" + n
				}
				return n
			}
			params = append(params, &ParamInfo{
				Name:  name.Name,
				Typed: "map[" + named(kTyped, kPtr) + "]" + named(vTyped, vPtr),
				Ptr:   false,
			})
		}
	default:
		log.Printf("default params :%+v, %+v", t, t.(*ast.BinaryExpr))
	}
	return params
}

func genVarName(str string) string {
	if str == "" {
		return str
	}
	splits := strings.Split(str, ".")
	if len(splits) != 1 {
		str = splits[len(splits)-1]
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return "_" + string(runes)
}
func getTyped(ns string, expr ast.Expr) (name string, ptr bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return func() string {
			if ns == "" || !t.IsExported() {
				return t.Name
			}
			return ns + "." + t.Name
		}(), false
	case *ast.SelectorExpr:
		typed, _ := getTyped(ns, t.Sel)
		return fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, typed), false
	case *ast.InterfaceType:
		return "interface{}", false
	case *ast.Ellipsis:
		return getTyped(ns, t.Elt)
	case *ast.StarExpr:
		typed, _ := getTyped(ns, t.X)
		return typed, true
	case *ast.ArrayType:
		log.Printf("ArrayType params :%+v, %+v", t, name)
		typed, _ := getTyped(ns, t.Elt)
		return "[]" + typed, false
	default:
		log.Printf("default params :%+v", t.(*ast.BinaryExpr))
		return "", false
	}
}

func (v *visitor) loadVendorAndGetInterface(ns string, name string) *InterfaceInfo {
	inter, ok := v.parsedInfo.VendorInterfaces[name]
	if ok {
		return inter
	}
	return nil
}
func removeExtension(filePath string) string {
	// 获取文件名（包含扩展名）
	baseName := filepath.Base(filePath)
	// 获取扩展名
	ext := filepath.Ext(baseName)
	// 如果文件有扩展名，则去除扩展名返回；否则原样返回
	if ext != "" {
		return baseName[:len(baseName)-len(ext)]
	}
	return baseName
}
func Must[T any](v T, err error) T {
	if err == nil {
		return v
	} else {
		panic(err)
	}
}

var VarIndex = func() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}()

type FuncType struct {
	Name string
	Type string
}

type ImportInfo struct {
	Name string
	Path string
}

type InterfaceInfo struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Rev     *ParamInfo
	Name    string
	CtxName string
	Args    []*ParamInfo
	Returns []*ParamInfo
}

type ParamInfo struct {
	Name     string
	Typed    string
	Ptr      bool
	Ellipsis bool
}
