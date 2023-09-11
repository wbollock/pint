package checks

import (
	"context"
	"fmt"

	"github.com/prometheus/common/model"

	"github.com/cloudflare/pint/internal/discovery"
	"github.com/cloudflare/pint/internal/parser"
)

const (
	AlertForCheckName = "alerts/for"
)

func NewAlertsForCheck() AlertsForChecksFor {
	return AlertsForChecksFor{}
}

type AlertsForChecksFor struct{}

func (c AlertsForChecksFor) Meta() CheckMeta {
	return CheckMeta{IsOnline: false}
}

func (c AlertsForChecksFor) String() string {
	return AlertForCheckName
}

func (c AlertsForChecksFor) Reporter() string {
	return AlertForCheckName
}

func (c AlertsForChecksFor) Check(_ context.Context, _ string, rule parser.Rule, _ []discovery.Entry) (problems []Problem) {
	if rule.AlertingRule == nil {
		return problems
	}

	if rule.AlertingRule.For != nil {
		problems = append(problems, c.checkField(rule.AlertingRule.For.Key.Value, rule.AlertingRule.For.Value.Value, rule.AlertingRule.For.Lines())...)
	}
	if rule.AlertingRule.KeepFiringFor != nil {
		problems = append(problems, c.checkField(rule.AlertingRule.KeepFiringFor.Key.Value, rule.AlertingRule.KeepFiringFor.Value.Value, rule.AlertingRule.KeepFiringFor.Lines())...)
	}

	return problems
}

func (c AlertsForChecksFor) checkField(name, value string, lines []int) (problems []Problem) {
	d, err := model.ParseDuration(value)
	if err != nil {
		problems = append(problems, Problem{
			Fragment: value,
			Lines:    lines,
			Reporter: c.Reporter(),
			Text:     fmt.Sprintf("invalid duration: %s", err),
			Severity: Bug,
		})
		return problems
	}

	if d == 0 {
		problems = append(problems, Problem{
			Fragment: value,
			Lines:    lines,
			Reporter: c.Reporter(),
			Text:     fmt.Sprintf("%q is the default value of %q, consider removing this line", value, name),
			Severity: Information,
		})
	}

	return problems
}
