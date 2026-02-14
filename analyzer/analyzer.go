package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"regexp"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var specialCharRe = regexp.MustCompile(`[^\p{L}\p{N}\s\-_:%]`)                          // Шаблон для поиска запрещённых спецсимволов
var sensitiveRe = regexp.MustCompile(`(?i)\b(password|api[_-]?key|token|secret|key)\b`) // Шаблон для поиска чувствительных данных

// logMethods содержит список методов логгеров, которые проверяются линтером
var logMethods = map[string]struct{}{
	"Info":   {},
	"Error":  {},
	"Warn":   {},
	"Debug":  {},
	"Infow":  {},
	"Warnw":  {},
	"Debugw": {},
	"Errorw": {},
	"Fatal":  {},
	"Panic":  {},
	"DPanic": {},
}

// Analyzer - основной объект линтера
var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages according to style rules",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

// run - основной entry-point. Проходит по всем файлам пакета и ищет вызовы функций.
// Для каждого вызова функции проверяет, логгер ли, извлекает сообщение и проверяет правила
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if logFunc, msg, pos := analyzeLogCall(pass, call); logFunc != "" && msg != "" {
					checkRules(pass, logFunc, msg, pos)
				}
			}
			return true
		})
	}
	return nil, nil
}

// analyzeLogCall анализирует конкретный вызов функции.
//
// analyzeLogCall возвращает имя метода, строку из сообщения, позицию в исходнике.
// Если вызов не является логгером или нет сообщений - возвращает "", "", token.NoPos.
func analyzeLogCall(pass *analysis.Pass, call *ast.CallExpr) (string, string, token.Pos) {

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", token.NoPos
	}

	if _, ok := logMethods[sel.Sel.Name]; !ok {
		return "", "", token.NoPos
	}

	recvIdent := unwrapLogger(sel.X)
	if recvIdent == nil {
		return "", "", token.NoPos
	}

	typ := pass.TypesInfo.TypeOf(recvIdent)
	isLogger := false

	if typ != nil && isLoggerType(typ) {
		isLogger = true
	} else {
		obj := pass.TypesInfo.ObjectOf(recvIdent)
		if pkg, ok := obj.(*types.PkgName); ok {
			path := pkg.Imported().Path()
			if path == "log/slog" || path == "go.uber.org/zap" {
				isLogger = true
			}
		}
	}

	if !isLogger {
		return "", "", token.NoPos
	}

	if len(call.Args) == 0 {
		return "", "", token.NoPos
	}

	msg, pos := extractMessage(call.Args[0])
	if msg != "" {
		return sel.Sel.Name, msg, pos
	}

	return "", "", token.NoPos
}

// unwrapLogger разворачивает выражение и ищет основной идентификатор логгера
func unwrapLogger(expr ast.Expr) *ast.Ident {
	for {
		switch v := expr.(type) {
		case *ast.Ident:
			return v

		case *ast.CallExpr:
			if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
				expr = sel.X
				continue
			}
			return nil

		case *ast.SelectorExpr:
			expr = v.X
			continue

		default:
			return nil
		}
	}
}

// isLoggerType проверяет, является ли тип slog или zap
func isLoggerType(typ types.Type) bool {
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	pkgPath := obj.Pkg().Path()
	typeName := obj.Name()

	if pkgPath == "log/slog" && typeName == "Logger" {
		return true
	}

	if pkgPath == "go.uber.org/zap" && typeName == "Logger" {
		return true
	}

	return false
}

// extractMessage рекурсивно извлекает строку из аргумента вызова
func extractMessage(expr ast.Expr) (string, token.Pos) {
	switch e := expr.(type) {

	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if str, err := strconv.Unquote(e.Value); err == nil {
				return str, e.Pos()
			}
		}

	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left, lpos := extractMessage(e.X)
			right, rpos := extractMessage(e.Y)

			if lpos != token.NoPos {
				return left + right, lpos
			}
			return left + right, rpos
		}

	case *ast.ParenExpr:
		return extractMessage(e.X)

	case *ast.CallExpr:
		var combined string
		var pos token.Pos
		for _, arg := range e.Args {
			str, p := extractMessage(arg)
			if str != "" {
				combined += str
				if pos == token.NoPos {
					pos = p
				}
			}
		}
		return combined, pos
	}

	return "", token.NoPos
}

// checkRules проверяет правила стиля лог-сообщений и вызывает pass.Reportf для diagnostics
func checkRules(pass *analysis.Pass, funcName, msg string, pos token.Pos) {
	var issues []string

	if len(msg) > 0 {
		runes := []rune(msg)
		if len(runes) > 0 && unicode.IsUpper(runes[0]) {
			issues = append(issues, "message must start with lowercase letter")
		}
	}

	if hasNonEnglish(msg) {
		issues = append(issues, "message must be in English only")
	}

	if hasSpecialChars(msg) {
		issues = append(issues, "no special characters or emojis allowed")
	}

	if hasSensitiveData(msg) {
		issues = append(issues, "no sensitive data keywords allowed")
	}

	for _, issue := range issues {
		pass.Reportf(pos, "log.%s: %s", funcName, issue)
	}
}

// hasNonEnglish проверяет наличие букв не в ASCII
func hasNonEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

// hasSpecialChars проверяет наличие запрещённых символов
func hasSpecialChars(s string) bool {
	return specialCharRe.MatchString(s)
}

// hasSensitiveData проверяет наличие ключевых слов чувствительной информации
func hasSensitiveData(s string) bool {
	return sensitiveRe.MatchString(s)
}
