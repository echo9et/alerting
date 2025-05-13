package main

// Package main содержит multichecker, объединяющий несколько статических анализаторов.
//
// Включенные анализаторы:
// - atomic: проверяет корректность использования пакета sync/atomic.
// - bools: обнаруживает подозрительные булевы выражения.
// - errorsas: проверяет правильность использования errors.As.
// - printf: проверяет форматирование строковых операций.
// - shadow: обнаруживает затенение переменных.
// - structtag: проверяет корректность тегов структур.
// - tests: проверяет тестовые файлы.
// - unmarshal: проверяет корректность работы с JSON/XML.
// - unreachable: обнаруживает недостижимый код.
// - unsafeptr: проверяет корректность использования unsafe.Pointer.
// - staticcheck (SA): анализаторы из пакета staticcheck.io.
// - errcheck: проверяет необработанные ошибки.
// - unused: обнаруживает неиспользуемые параметры функций.
// - noosexitinmain: собственный анализатор, запрещающий использование os.Exit в main.
//
// Запуск multichecker:
//   go run cmd/staticlint/main.go ./...

import (
	"go/ast"

	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	"honnef.co/go/tools/staticcheck"
)

var noOsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noosexitinmain",
	Doc:  "Запрещает использование os.Exit в функции main пакета main.",
	Run:  runNoOsExitInMain,
}

func runNoOsExitInMain(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Проверяем, что файл принадлежит пакету main
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "main" {
				return true
			}

			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkgName, ok := selExpr.X.(*ast.Ident)
				if !ok || pkgName.Name != "os" || selExpr.Sel.Name != "Exit" {
					return true
				}

				pass.Reportf(callExpr.Pos(), "вызов os.Exit в функции main пакета main запрещен")
				return false
			})

			return false
		})
	}

	return nil, nil
}

// запуск анализатора для всех файлов в текущей директории
//
// go run cmd/staticlint/main.go ./...
func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	for _, analyzer := range staticcheck.Analyzers {
		analyzers = append(analyzers, analyzer.Analyzer)
	}
	analyzers = append(analyzers, errcheck.Analyzer)
	analyzers = append(analyzers, noOsExitInMainAnalyzer)

	multichecker.Main(analyzers...)
}
