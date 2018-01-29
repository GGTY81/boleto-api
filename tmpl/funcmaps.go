package tmpl

import (
	"bytes"
	"html"
	"html/template"
	"regexp"
	"strings"
	"time"

	"strconv"

	"fmt"

	"github.com/kennygrant/sanitize"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/util"
)

var funcMap = template.FuncMap{
	"today":                  today,
	"todayCiti":              todayCiti,
	"brdate":                 brDate,
	"replace":                replace,
	"docType":                docType,
	"trim":                   trim,
	"padLeft":                padLeft,
	"clearString":            clearString,
	"toString":               toString,
	"toString64":             toString64,
	"fmtDigitableLine":       fmtDigitableLine,
	"fmtCNPJ":                fmtCNPJ,
	"fmtCPF":                 fmtCPF,
	"fmtDoc":                 fmtDoc,
	"truncate":               truncateString,
	"fmtNumber":              fmtNumber,
	"joinSpace":              joinSpace,
	"brDateWithoutDelimiter": brDateWithoutDelimiter,
	"enDateWithoutDelimiter": enDateWithoutDelimiter,
	"fullDate":               fulldate,
	"enDate":                 enDate,
	"hasErrorTags":           hasErrorTags,
	"toFloatStr":             toFloatStr,
	"concat":                 concat,
	"base64":                 base64,
	"unscape":                unscape,
	"unescapeHtmlString":     unescapeHtmlString,
	"trimLeft":               trimLeft,
	"santanderNSUPrefix":     santanderNSUPrefix,
	"santanderEnv":           santanderEnv,
	"formatSingleLine":       formatSingleLine,
	"diff":                   diff,
	"mod11dv":                calculateOurNumberMod11,
	"mod10dv":                mod10Itau,
	"printIfNotProduction":   printIfNotProduction,
	"itauEnv":                itauEnv,
	"extractNumbers":         extractNumbers,
	"splitValues":            splitValues,
	"brDateDelimiter":        brDateDelimiter,
	"brDateDelimiterTime":    brDateDelimiterTime,
}

func GetFuncMaps() template.FuncMap {
	return funcMap
}

func santanderNSUPrefix(number string) string {
	if config.Get().DevMode {
		return "TST" + number
	}
	return number
}

func diff(a string, b string) bool {
	return a != b
}

func formatSingleLine(s string) string {
	s1 := strings.Replace(s, "\r", "", -1)
	return strings.Replace(s1, "\n", "; ", -1)
}

func santanderEnv() string {
	if config.Get().DevMode {
		return "T"
	}
	return "P"
}

func padLeft(value, char string, total uint) string {
	s := util.PadLeft(value, char, total)
	return s
}
func unscape(s string) template.HTML {
	return template.HTML(s)
}

func sanitizeHtmlString(s string) string {
	str := html.UnescapeString(s)
	return sanitize.HTML(str)
}

func unescapeHtmlString(s string) template.HTML {
	c := sanitizeHtmlString(s)
	return template.HTML(html.UnescapeString(c))
}

func trimLeft(s string, caract string) string {
	return strings.TrimLeft(s, caract)
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		bnoden = str[0:num]
	}
	return bnoden
}

func clearString(str string) string {
	s := sanitize.Accents(str)
	var buffer bytes.Buffer
	for _, ch := range s {
		if ch <= 122 && ch >= 32 {
			buffer.WriteString(string(ch))
		}
	}
	return buffer.String()
}

func joinSpace(str ...string) string {
	return strings.Join(str, " ")
}

func hasErrorTags(mapValues map[string]string, errorTags ...string) bool {
	hasError := false
	for _, v := range errorTags {
		if value, exist := mapValues[v]; exist && strings.Trim(value, " ") != "" {
			hasError = true
			break
		}
	}
	return hasError
}

func fmtNumber(n uint64) string {
	real := n / 100
	cents := n % 100
	return fmt.Sprintf("%d,%02d", real, cents)
}

func printIfNotProduction(obj string) string {
	if config.IsNotProduction() {
		return fmt.Sprintf("%s", obj)
	}
	return ""
}

func toFloatStr(n uint64) string {
	real := n / 100
	cents := n % 100
	return fmt.Sprintf("%d.%02d", real, cents)
}

func fmtDoc(doc models.Document) string {
	if e := doc.ValidateCPF(); e == nil {
		return fmtCPF(doc.Number)
	}
	return fmtCNPJ(doc.Number)
}

func toString(number uint) string {
	return strconv.FormatInt(int64(number), 10)
}

func toString64(number uint64) string {
	return strconv.FormatInt(int64(number), 10)
}

func today() time.Time {
	return util.BrNow()
}

func todayCiti() time.Time {
	return util.NycNow()
}

func fulldate(t time.Time) string {
	return t.Format("20060102150405")
}

func brDate(d time.Time) string {
	return d.Format("02/01/2006")
}

func enDate(d time.Time, del string) string {
	return d.Format("2006" + del + "01" + del + "02")
}

func brDateWithoutDelimiter(d time.Time) string {
	return d.Format("02012006")
}

func enDateWithoutDelimiter(d time.Time) string {
	return d.Format("20060102")
}

func replace(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}

func docType(s models.Document) int {
	if s.IsCPF() {
		return 1
	}
	return 2
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
func fmtDigitableLine(s string) string {
	buf := bytes.Buffer{}
	for idx, c := range s {
		if idx == 5 || idx == 15 || idx == 26 {
			buf.WriteString(".")
		}
		if idx == 10 || idx == 21 || idx == 32 || idx == 33 {
			buf.WriteString(" ")
		}
		buf.WriteByte(byte(c))
	}
	return buf.String()
}

func fmtCNPJ(s string) string {
	buf := bytes.Buffer{}
	for idx, c := range s {
		if idx == 2 || idx == 5 {
			buf.WriteString(".")
		}
		if idx == 8 {
			buf.WriteString("/")
		}
		if idx == 12 {
			buf.WriteString("-")
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func fmtCPF(s string) string {
	buf := bytes.Buffer{}
	for idx, c := range s {
		if idx == 3 || idx == 6 {
			buf.WriteString(".")
		}
		if idx == 9 {
			buf.WriteString("-")
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func concat(s ...string) string {
	buf := bytes.Buffer{}
	for _, item := range s {
		buf.WriteString(item)
	}
	return buf.String()
}

func base64(s string) string {
	return util.Base64(s)
}

func calculateOurNumberMod11(number uint) uint {
	ourNumberWithDigit := strconv.Itoa(int(number)) + util.OurNumberDv(strconv.Itoa(int(number)), util.MOD11)
	value, _ := strconv.Atoi(ourNumberWithDigit)
	return uint(value)
}

func mod10Itau(number string, agency string, account string, wallet uint16) string {

	var buffer bytes.Buffer

	if wallet == 126 || wallet == 131 || wallet == 146 || wallet == 168 {

		buffer.WriteString(strconv.FormatUint(uint64(wallet), 10))
		buffer.WriteString(number)

		return util.OurNumberDv(buffer.String(), util.MOD10)
	} else {
		buffer.WriteString(agency)
		buffer.WriteString(account)
		buffer.WriteString(strconv.FormatUint(uint64(wallet), 10))
		buffer.WriteString(number)
		return util.OurNumberDv(buffer.String(), util.MOD10)
	}
}

func itauEnv() string {
	if config.Get().DevMode {
		return "1"
	}
	return "2"
}

func extractNumbers(value string) string {
	re := regexp.MustCompile("(\\D+)")
	sanitizeValue := re.ReplaceAllString(string(value), "")
	return sanitizeValue
}

func splitValues(value string, init int, end int) string {
	return value[init:end]
}

func brDateDelimiter(date string, del string) string {
	layout := "2006-01-02"
	d, err := time.Parse(layout, date)
	if err != nil {
		return date
	}

	return d.Format("02" + del + "01" + del + "2006")
}

func brDateDelimiterTime(date time.Time, del string) string {
	layout := "2006-01-02 00:00:00 +0000 UTC"

	d, err := time.Parse(layout, date.String())

	if err != nil {
		return date.String()
	}

	return d.Format("02" + del + "01" + del + "2006")
}
