package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var goFilePattern = regexp.MustCompile(`\.go$`)

type rule struct {
	name        string
	description string
	pattern     *regexp.Regexp
}

type violation struct {
	rule       rule
	path       string
	lineNumber int
	line       string
}

var rules = []rule{
	{
		name:        "no_direct_logx_import",
		description: "business code must not import go-zero logx directly",
		pattern:     regexp.MustCompile(`"github\.com/zeromicro/go-zero/core/logx"`),
	},
	{
		name:        "no_direct_logx_usage",
		description: "business code must not call logx directly",
		pattern:     regexp.MustCompile(`\blogx\.`),
	},
	{
		name:        "prefer_errors_is_over_errs_iscode",
		description: "new code must prefer errors.Is(err, errs.E(...)) over errs.IsCode",
		pattern:     regexp.MustCompile(`\berrs\.IsCode\(`),
	},
	{
		name:        "no_string_match_on_error_contains",
		description: "do not branch business logic by strings.Contains(err.Error(), ...)",
		pattern:     regexp.MustCompile(`strings\.Contains\(\s*err\.Error\(\)`),
	},
	{
		name:        "no_string_match_on_error_equals",
		description: "do not branch business logic by err.Error() ==",
		pattern:     regexp.MustCompile(`err\.Error\(\)\s*==`),
	},
}

var (
	scanRoots = []string{
		"services",
		"pkg",
	}
	excludedDirs = []string{
		filepath.Clean(filepath.Join("pkg", "logs")),
		filepath.Clean(filepath.Join("services", "content", "pb")),
		filepath.Clean(filepath.Join("services", "identity", "pb")),
	}
)

// Run executes the static review rule scan against business code.
// Run 对业务代码执行静态 review 规则扫描。
func Run() error {
	violations, err := collectViolations()
	if err != nil {
		return err
	}

	if len(violations) == 0 {
		fmt.Println("review rules passed")
		return nil
	}

	var builder strings.Builder
	builder.WriteString("review rule violations found:\n")
	for _, item := range violations {
		builder.WriteString(fmt.Sprintf(
			"- [%s] %s:%d: %s\n  %s\n",
			item.rule.name,
			filepath.ToSlash(item.path),
			item.lineNumber,
			item.rule.description,
			strings.TrimSpace(item.line),
		))
	}

	return fmt.Errorf("%s", builder.String())
}

func collectViolations() ([]violation, error) {
	var violations []violation

	for _, root := range scanRoots {
		if _, err := os.Stat(root); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			cleanPath := filepath.Clean(path)
			if info.IsDir() {
				if shouldSkipDir(cleanPath) {
					return filepath.SkipDir
				}
				return nil
			}

			if !goFilePattern.MatchString(cleanPath) {
				return nil
			}

			fileViolations, err := scanFile(cleanPath)
			if err != nil {
				return err
			}
			violations = append(violations, fileViolations...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		if violations[i].path == violations[j].path {
			if violations[i].lineNumber == violations[j].lineNumber {
				return violations[i].rule.name < violations[j].rule.name
			}
			return violations[i].lineNumber < violations[j].lineNumber
		}
		return violations[i].path < violations[j].path
	})

	return violations, nil
}

func scanFile(path string) ([]violation, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var violations []violation

	for index, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		for _, currentRule := range rules {
			if currentRule.pattern.MatchString(line) {
				violations = append(violations, violation{
					rule:       currentRule,
					path:       path,
					lineNumber: index + 1,
					line:       rawLine,
				})
			}
		}
	}

	return violations, nil
}

func shouldSkipDir(path string) bool {
	for _, excluded := range excludedDirs {
		if path == excluded {
			return true
		}
	}
	return false
}
