//go build handlers_gen/* && ./codegen api.go api_handlers.go

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func CreateHTTPwrapper(decl string) {

}
func structToString(name string, structType *ast.StructType) string {
	var buf bytes.Buffer
	// Создаем новый принтер
	printer.Fprint(&buf, token.NewFileSet(), &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent(name), Type: structType}},
	})
	return buf.String()
}
func main() {
	//реализация через go run codegen.go
	fset := token.NewFileSet()
	var code string
	filepath := filepath.Join("..", "api.go") // путь к апи
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %s\n", err)
		return
	}

	code = string(content)
	//передаем весь текст в строку и парсим
	f, err := parser.ParseFile(fset, filepath, code, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	container := ""
	bufStruct := ""
	for _, decl := range f.Decls {
		//проверяем является ли decl дженерик обьявлением  import, constant, type or variable declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if ok {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structStr := structToString(typeSpec.Name.Name, structType)
						// если в название структуры есть ключевое своло Api контейнер перезаписывается
						switch {
						case strings.Contains(typeSpec.Name.Name, "Api"):
							bufStruct = typeSpec.Name.Name
							container = structStr
						default:
							container = container + "\n" + structStr
						}

					} // нужно записать структуры относящиеся к определенному апи
				}
			}
		}
		funcDecl, ok := decl.(*ast.FuncDecl) //проверяем является ли обьявление функцией
		if ok && funcDecl.Recv != nil {
			for _, field := range funcDecl.Recv.List {
				if starExpr, ok := field.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == bufStruct {
						fmt.Printf("Method %s is associated with struct:%s (pointer receiver)\n", funcDecl.Name.Name, bufStruct)
					} // вышло отследить методы относящиеся к читаемой в данный момент структуре

				}
			}

		}
	}
}
