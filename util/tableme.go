package util

import (
	"fmt"
	config "github.com/alexebird/aws-find/config"
	"github.com/alexebird/tableme/tableme"
)

func isInSlice(target string, slice []string) bool {
	for _, e := range slice {
		if target == e {
			return true
		}
	}

	return false
}

func PrintColorizedTable(bytes []byte, subcmd string, colorizeConfig []config.ColorizeConfig) {
	colorRules := make([]*tableme.ColorRule, 0)

	for _, rule := range colorizeConfig {
		if isInSlice(subcmd, rule.Subcmds) {
			colorRules = append(colorRules, &tableme.ColorRule{
				Pattern: rule.Regex,
				Color:   rule.Color,
			})
		}
	}

	colored := tableme.Colorize(bytes, colorRules)
	fmt.Print(colored.String())
}
