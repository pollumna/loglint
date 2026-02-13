package analyzer

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var (
	specialCharRe  = regexp.MustCompile(`[^\p{L}\p{N}\s-]`)
	multipleDotsRe = regexp.MustCompile(`\.{2,}`)
	multipleExclRe = regexp.MustCompile(`!{2,}`)
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages according to style rules",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	slogAliases, zapAliases := collectLogImports(pass)

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if logFunc, msg, pos := analyzeLogCall(pass, call, slogAliases, zapAliases); logFunc != "" && msg != "" {
					checkRules(pass, logFunc, msg, pos)
				}
			}
			return true
		})
	}
	return nil, nil
}

func collectLogImports(pass *analysis.Pass) (map[string]bool, map[string]bool) {
	slogAliases, zapAliases := make(map[string]bool), make(map[string]bool)

	for _, file := range pass.Files {
		for _, imp := range file.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}

			alias := "slog"
			if imp.Name != nil {
				alias = imp.Name.Name
			}

			switch path {
			case "log/slog":
				slogAliases[alias] = true
			case "go.uber.org/zap":
				zapAliases[alias] = true
			}
		}
	}
	return slogAliases, zapAliases
}

func analyzeLogCall(pass *analysis.Pass, call *ast.CallExpr, slogAliases, zapAliases map[string]bool) (string, string, token.Pos) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", token.NoPos
	}

	logMethods := map[string]bool{"Info": true, "Error": true, "Warn": true, "Debug": true}
	if !logMethods[sel.Sel.Name] {
		return "", "", token.NoPos
	}

	recvIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", "", token.NoPos
	}

	if slogAliases[recvIdent.Name] || zapAliases[recvIdent.Name] {
		if len(call.Args) == 0 {
			return "", "", token.NoPos
		}
		msg, _ := extractMessage(call.Args[0])
		if msg != "" {
			return sel.Sel.Name, msg, sel.Pos()
		}
	}

	typ := pass.TypesInfo.TypeOf(recvIdent)
	if typ == nil {
		return "", "", token.NoPos
	}

	typeStr := typ.String()
	isSlogLogger := strings.Contains(typeStr, "slog.Logger")
	isZapLogger := strings.Contains(typeStr, "zap.Logger")

	if isSlogLogger || isZapLogger {
		if len(call.Args) == 0 {
			return "", "", token.NoPos
		}
		msg, _ := extractMessage(call.Args[0])
		if msg != "" {
			return sel.Sel.Name, msg, sel.Pos()
		}
	}

	return "", "", token.NoPos
}

func extractMessage(arg ast.Expr) (string, token.Pos) {
	switch expr := arg.(type) {
	case *ast.BasicLit:
		if expr.Kind == token.STRING {
			if str, err := strconv.Unquote(expr.Value); err == nil {
				return str, expr.Pos()
			}
		}

	case *ast.BinaryExpr:
		if expr.Op == token.ADD {
			leftText, _ := extractMessage(expr.X)
			rightText, _ := extractMessage(expr.Y)
			return leftText + rightText, expr.Pos()
		}

	case *ast.ParenExpr:
		return extractMessage(expr.X)

	default:
		return "", expr.Pos()
	}
	return "", arg.Pos()
}

func checkRules(pass *analysis.Pass, funcName, msg string, pos token.Pos) {
	issues := []string{}

	if len(msg) > 0 {
		firstRune := rune(msg[0])
		if 'A' <= firstRune && firstRune <= 'Z' {
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

func hasNonEnglish(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

func hasSpecialChars(s string) bool {
	if specialCharRe.MatchString(s) {
		return true
	}
	if multipleDotsRe.MatchString(s) {
		return true
	}
	if multipleExclRe.MatchString(s) {
		return true
	}
	return false
}

func hasSensitiveData(s string) bool {
	keywords := []string{"password", "api_key", "token", "secret", "key", "pass"}
	s = strings.ToLower(s)
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}
