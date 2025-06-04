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
	container := make(map[string]string)
	// bufStruct := ""
	for _, decl := range f.Decls {
		//проходим по всем декларациям
		//если дженерик является структурой добавляем ключ в карту
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structType := structToString(typeSpec.Name.Name, structType)
						// println(structType, typeSpec.Name.Name)
						container[structType] = typeSpec.Name.Name
					}
				}
			}
			// println("there is an stuct")
		}
		//если функция добавляем ее по ключу в карту
		// funcDecl, ok := decl.(*ast.FuncDecl)
		// if ok {
		// }
	}
	for key := range container {
		println("key is: ", key)
		println("value is: ", container[key])
	}
	// 	//проверяем является ли decl дженерик обьявлением  import, constant, type or variable declaration
	// 	genDecl, ok := decl.(*ast.GenDecl)
	// 	if ok {
	// 		for _, spec := range genDecl.Specs {
	// 			typeSpec, ok := spec.(*ast.TypeSpec)
	// 			if ok {
	// 				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
	// 					structStr := structToString(typeSpec.Name.Name, structType)
	// 					// при обнаружении структуры идем по файлу и ищем все функции
	// 					for _, decl := range f.Decls {
	// 						funcDecl, ok := decl.(*ast.FuncDecl) //проверяем является ли обьявление функцией
	// 						if ok && funcDecl.Recv != nil {
	// 							for _, field := range funcDecl.Recv.List {
	// 								if starExpr, ok := field.Type.(*ast.StarExpr); ok {
	// 									if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == bufStruct {
	// 										fmt.Printf("Method %s is associated with struct:%s (pointer receiver)\n", funcDecl.Name.Name, bufStruct)

	// 									}

	// 								} // нужно записать структуры относящиеся к определенному апи
	// 							}
	// 						}

	// 					}

	// 					// если в название структуры есть ключевое своло Api контейнер перезаписывается
	// 					switch {
	// 					case strings.Contains(typeSpec.Name.Name, "Api"):
	// 						bufStruct = typeSpec.Name.Name
	// 						container = structStr
	// 					default:
	// 						container = container + "\n" + structStr
	// 					}
	// 				} // вышло отследить методы относящиеся к читаемой в данный момент структуре

	// 			}
	// 		}

	// 	}
	// }
}
