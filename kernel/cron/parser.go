// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cron

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

func normalize(exp string) string {
	exp = strings.ReplaceAll(exp, "CRON_TZ=", "TZ=")
	repl := map[string]string{
		"@yearly":    "0 0 1 1 *",
		"@annually":  "0 0 1 1 *",
		"@monthly":   "0 0 1 * *",
		"@weekly":    "0 0 * * 0",
		"@daily":     "0 0 * * *",
		"@hourly":    "0 * * * *",
		"@sunday":    "0 0 * * 0",
		"@monday":    "0 0 * * 1",
		"@tuesday":   "0 0 * * 2",
		"@wednesday": "0 0 * * 3",
		"@thursday":  "0 0 * * 4",
		"@friday":    "0 0 * * 5",
		"@saturday":  "0 0 * * 6",
	}
	for k, v := range repl {
		if strings.Contains(exp, k) {
			exp = strings.ReplaceAll(exp, k, v)
			return exp
		}
	}
	return exp
}

func replaceMonth(exp string) string {
	repl := map[string]string{
		"JAN": "1",
		"FEB": "2",
		"MAR": "3",
		"APR": "4",
		"MAY": "5",
		"JUN": "6",
		"JUL": "7",
		"AUG": "8",
		"SEP": "9",
		"OCT": "10",
		"NOV": "11",
		"DEC": "12",
	}
	if regexp.MustCompile(`[a-zA-Z]{4}`).MatchString(exp) {
		return exp // Contains invalid expression.
	}
	for k, v := range repl {
		exp = strings.ReplaceAll(exp, k, v)
	}
	return exp
}

func replaceWeek(exp string) string {
	repl := map[string]string{
		"SUN": "0",
		"MON": "1",
		"TUE": "2",
		"WED": "3",
		"THU": "4",
		"FRI": "5",
		"SAT": "6",
	}
	if regexp.MustCompile(`[a-zA-Z]{4}`).MatchString(exp) {
		return exp // Contains invalid expression.
	}
	for k, v := range repl {
		exp = strings.ReplaceAll(exp, k, v)
	}
	return exp
}

func mask(min, max, step int) uint64 {
	v := uint64(0)
	for i := min; i <= max; i += step {
		v |= 1 << i
	}
	return v
}

// parseValue parses each fields of cron expression.
// Allowed expressions are listed below.
// Other expression results in error, or false at second returned value.
// It is always be failure if the given min and max is min>max,
// Both min and max MUST be zero or grater than zero.
//
//   - Wildcard           : "*"
//   - Number             : "5"
//   - Range              : "10-20"
//   - Wildcard with step : "*/3"
//   - Number with step   : "5/3"
//   - Range with step    : "10-20/3"
func parseValue(exp string, min, max int) (uint64, bool) {
	exp = strings.Trim(exp, " \n\r\t\f,")

	if min > max || min < 0 || max < 0 { // This is not allowed.
		return 0, false
	}

	all := strings.Split(exp, ",")
	result := uint64(0)
	for _, e := range all {
		fields := strings.Split(e, "/")
		switch len(fields) {
		case 1:
			ini, end, ok := parseRange(fields[0], min, max)
			if !ok {
				return 0, false
			}
			result |= mask(ini, end, 1)

		case 2:
			step, err := strconv.Atoi(fields[1])
			if err != nil {
				return 0, false
			}
			if step <= 0 {
				return 0, false
			}
			ini, end, ok := parseRange(fields[0], min, max)
			if !ok {
				return 0, false
			}
			if ini == end {
				end = max
			}
			result |= mask(ini, end, step)

		default:
			return 0, false
		}
	}
	return result, true
}

// parseRange returns the value range of the given expression.
// Allowed formats are wildcard, number and range.
// It is always be failure if the given min and max is min>max,
// Both min and max MUST be zero or grater than zero.
//
//   - Wildcard : "*"      returns min,max
//   - Number   : "5"      returns 5,5
//   - Range    : "10-20"  returns 10,20
func parseRange(exp string, min, max int) (int, int, bool) {
	exp = strings.Trim(exp, " \n\r\t\f,")

	if min > max || min < 0 || max < 0 { // This is not allowed.
		return 0, 0, false
	}

	if exp == "*" {
		return min, max, true
	}

	if !strings.Contains(exp, "-") {
		val, err := strconv.Atoi(exp)
		if err != nil {
			return 0, 0, false
		}
		if val < min || val > max {
			return 0, 0, false
		}
		return val, val, true
	}

	fields := strings.Split(exp, "-")
	if len(fields) != 2 {
		return 0, 0, false
	}
	ini, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, false
	}
	end, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, false
	}
	if ini > end || ini < min || end > max {
		return 0, 0, false
	}

	return ini, end, true
}

// Parse parses the given cron expression
// and returns Crontab object.
//
//	TZ=UTC * * * * * *
//	|      | | | | | |
//	|      | | | | | |- Day of week
//	|      | | | | |--- Month
//	|      | | | |----- Day of month
//	|      | | |------- Hour
//	|      | |--------- Minute
//	|      |----------- Second (Optional)
//	|------------------ Timezone (Optional)
//
//	Field name   | Mandatory  | Values          | Special characters
//	----------   | ---------- | --------------  | -------------------
//	Timezone     | No         | Timezone name   |
//	Second       | No         | 0-59            | * / , -
//	Minute       | Yes        | 0-59            | * / , -
//	Hours        | Yes        | 0-23            | * / , -
//	Day of month | Yes        | 1-31            | * / , -
//	Month        | Yes        | 1-12 or JAN-DEC | * / , -
//	Day of week  | Yes        | 0-6 or SUN-SAT  | * / , -
//
// See the references.
//   - https://en.wikipedia.org/wiki/Cron
//   - https://crontab.guru/
//   - https://crontab.cronhub.io/
func Parse(crontab string) (*Crontab, error) {
	crontab = strings.Trim(crontab, " \n\r\t\f,")
	normalized := normalize(crontab)
	fields := strings.Split(normalized, " ")

	loc := time.Local
	if strings.HasPrefix(fields[0], "TZ=") {
		parsedLoc, err := time.LoadLocation(strings.TrimPrefix(fields[0], "TZ="))
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeParse,
				Description: ErrDscParse,
				Detail:      "`" + crontab + "`",
			}).Wrap(err)
		}
		loc = parsedLoc
		fields = fields[1:]
	}

	c := &Crontab{
		loc:   loc,
		timer: time.Now,
	}
	switch len(fields) {
	case 5:
		fields = append([]string{"0"}, fields...)
	case 6:
	default:
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid number of fields `" + crontab + "`",
		}
	}
	var ok bool
	if c.second, ok = parseValue(fields[0], 0, 59); !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid second expression `" + fields[0] + "`",
		}
	}
	if c.minute, ok = parseValue(fields[1], 0, 59); !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid minute expression `" + fields[1] + "`",
		}
	}
	if c.hour, ok = parseValue(fields[2], 0, 23); !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid hour expression `" + fields[2] + "`",
		}
	}
	if c.day, ok = parseValue(fields[3], 1, 31); !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid day of month expression `" + fields[3] + "`",
		}
	}
	if c.month, ok = parseValue(replaceMonth(fields[4]), 1, 12); !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid month expression `" + fields[4] + "`",
		}
	}
	if c.week, ok = parseValue(replaceWeek(fields[5]), 0, 6); !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "invalid day of week expression `" + fields[5] + "`",
		}
	}
	if !c.valid() {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeParse,
			Description: ErrDscParse,
			Detail:      "unschedulable expression `" + crontab + "`",
		}
	}
	return c, nil
}
