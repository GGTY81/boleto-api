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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var funcMap = template.FuncMap{
	"today":                              today,
	"todayCiti":                          todayCiti,
	"brdate":                             brDate,
	"replace":                            replace,
	"docType":                            docType,
	"trim":                               trim,
	"padLeft":                            padLeft,
	"clearString":                        clearString,
	"toString":                           toString,
	"toString64":                         toString64,
	"fmtDigitableLine":                   fmtDigitableLine,
	"fmtCNPJ":                            fmtCNPJ,
	"fmtCPF":                             fmtCPF,
	"fmtDoc":                             fmtDoc,
	"truncate":                           truncateString,
	"fmtNumber":                          fmtNumber,
	"joinSpace":                          joinSpace,
	"brDateWithoutDelimiter":             brDateWithoutDelimiter,
	"enDateWithoutDelimiter":             enDateWithoutDelimiter,
	"fullDate":                           fulldate,
	"enDate":                             enDate,
	"hasErrorTags":                       hasErrorTags,
	"toFloatStr":                         toFloatStr,
	"concat":                             concat,
	"base64":                             base64,
	"unscape":                            unscape,
	"unescapeHtmlString":                 unescapeHtmlString,
	"trimLeft":                           trimLeft,
	"santanderNSUPrefix":                 santanderNSUPrefix,
	"santanderEnv":                       santanderEnv,
	"formatSingleLine":                   formatSingleLine,
	"diff":                               diff,
	"mod11dv":                            calculateOurNumberMod11,
	"mod10ItauDv":                        mod10Itau,
	"printIfNotProduction":               printIfNotProduction,
	"itauEnv":                            itauEnv,
	"caixaEnv":                           caixaEnv,
	"extractNumbers":                     extractNumbers,
	"splitValues":                        splitValues,
	"brDateDelimiter":                    brDateDelimiter,
	"brDateDelimiterTime":                brDateDelimiterTime,
	"toString16":                         toString16,
	"mod11BradescoShopFacilDv":           mod11BradescoShopFacilDv,
	"bsonMongoToString":                  bsonMongoToString,
	"escapeStringOnJson":                 escapeStringOnJson,
	"removeSpecialCharacter":             removeSpecialCharacter,
	"sanitizeCitibankSpecialCharacteres": sanitizeCitibankSpecialCharacteres,
	"clearStringCaixa":                   clearStringCaixa,
	"truncateOnly":                       truncateOnly,
	"toUint":                             toUint,
	"strToFloat":                         strToFloat,
	"float64ToString":                    float64ToString,
	"datePlusDays":                       datePlusDays,
	"datePlusDaysConsideringZeroAsStart": datePlusDaysConsideringZeroAsStart,
	"getInterestInstruction":             getInterestInstruction,
	"getFineInstruction":                 getFineInstruction,
	"datePlusDaysLocalTime":              datePlusDaysLocalTime,
	"calculateFees":                      calculateFees,
	"calculateInterestByDay":             calculateInterestByDay,
	"onlyAlphabetics":                    onlyAlphabetics,
	"onlyAlphanumerics":                  onlyAlphanumerics,
	"onlyOneSpace":                       onlyOneSpace,
	"removeAllSpaces":                    removeAllSpaces,
}

func GetFuncMaps() template.FuncMap {
	return funcMap
}

func santanderNSUPrefix(number string) string {
	if config.Get().SantanderEnv == "T" {
		return "TST" + number
	}
	return number
}

func santanderEnv() string {
	return config.Get().SantanderEnv
}

func diff(a string, b string) bool {
	return a != b
}

func formatSingleLine(s string) string {
	s1 := strings.Replace(s, "\r", "", -1)
	return strings.Replace(s1, "\n", "; ", -1)
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
	bnoden := removeSpecialCharacter(str)

	if len(bnoden) > num {
		bnoden = str[0:num]
	}
	//Support extended ASCII
	return string([]rune(bnoden))
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

func strToFloat(n string) float64 {
	s, _ := strconv.ParseFloat(n, 64)
	return s
}

func float64ToString(format string, value float64) string {
	return fmt.Sprintf(format, value)
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

func toUint(number string) uint {
	value, _ := strconv.Atoi(number)
	return uint(value)
}

func toString16(number uint16) string {
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

func datePlusDays(date time.Time, days uint) time.Time {
	timeToPlus := time.Hour * 24 * time.Duration(days)
	return date.UTC().Add(timeToPlus)
}

func datePlusDaysConsideringZeroAsStart(date time.Time, days uint) time.Time {
	daysConsideringZeroAsStart := days - 1
	return datePlusDays(date, daysConsideringZeroAsStart)
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

func calculateOurNumberMod11(number uint, onlyDigit bool) uint {

	ourNumberDigit := util.OurNumberDv(strconv.Itoa(int(number)), util.MOD11)

	if onlyDigit {
		value, _ := strconv.Atoi(ourNumberDigit)
		return uint(value)
	}

	ourNumberWithDigit := strconv.Itoa(int(number)) + ourNumberDigit
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
	return config.Get().ItauEnv
}

func caixaEnv() string {
	return config.Get().CaixaEnv
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

func mod11BradescoShopFacilDv(number string, wallet string) string {
	var buffer bytes.Buffer
	buffer.WriteString(wallet)
	buffer.WriteString(number)
	return util.OurNumberDv(buffer.String(), util.MOD11, 7)
}

func bsonMongoToString(bsonId primitive.ObjectID) string {
	return bsonId.Hex()
}

func escapeStringOnJson(field string) string {
	field = strings.Replace(field, "\b", "", -1)
	return regexp.MustCompile(`[\t\f\r\\]`).ReplaceAllString(field, "")
}

func removeSpecialCharacter(str string) string {
	return regexp.MustCompile("[^a-zA-Z0-9ÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç,.\\-\\s]+").ReplaceAllString(str, "")
}

func sanitizeCitibankSpecialCharacteres(str string, num int) string {
	str = regexp.MustCompile("[^a-zA-Z0-9.;@\\-\\/\\s]+").ReplaceAllString(clearString(str), "")

	if len(str) > num {
		str = str[0:num]
	}

	return str
}

//clearStringCaixa Define os caracteres aceitos de acordo com a documentação da caixa.
func clearStringCaixa(str string) string {
	s := sanitize.Accents(str)
	var buffer bytes.Buffer
	for _, ch := range s {
		if util.IsDigit(ch) || util.IsBasicCharacter(ch) || util.IsCaixaSpecialCharacter(ch) {
			buffer.WriteString(string(ch))
		} else {
			buffer.WriteString(" ")
		}
	}
	return buffer.String()
}

//truncateOnly Realiza o truncate da string
func truncateOnly(str string, num int) string {
	if len([]rune(str)) > num {
		str = string([]rune(str)[0:num])
	}
	return str
}

//calculateFees Calcula o fees em reais por dia, sobre o valor do titulo
func calculateFees(amountFee uint64, percentageFee float64, titleAmount uint64) float64 {
	const conversionFactor = 0.01
	const defaultValueRate = 0

	if amountFee <= defaultValueRate && percentageFee <= defaultValueRate {
		return defaultValueRate
	}

	if amountFee > defaultValueRate {
		return float64(amountFee) * conversionFactor
	}

	originalAmountToReal := float64(titleAmount) * conversionFactor
	return (percentageFee * conversionFactor) * originalAmountToReal
}

//calculateInterestByDay Calcula a taxa de juros em reais por dia, sobre o valor do titulo
func calculateInterestByDay(amountFee uint64, percentageFee float64, titleAmount uint64) float64 {
	interestAmount := calculateFees(amountFee, percentageFee, titleAmount)

	if percentageFee > 0 {
		interestAmount /= 30
	}

	return interestAmount
}

func datePlusDaysLocalTime(date time.Time, days uint) time.Time {
	timeToPlus := time.Hour * 24 * time.Duration(days)
	return date.Add(timeToPlus)
}

//getFineInstruction Obtém a instrução de multa
func getFineInstruction(title models.Title) string {
	dateFine := datePlusDaysLocalTime(title.ExpireDateTime, title.Fees.Fine.DaysAfterExpirationDate)
	dateFineFormatted := brDate(dateFine)

	fineAmountInReal := calculateFees(title.Fees.Fine.AmountInCents, title.Fees.Fine.PercentageOnTotal, title.AmountInCents)

	return fmt.Sprintf("APOS %s: MULTA..........R$ %.2f", dateFineFormatted, fineAmountInReal)
}

//getInterestInstruction Obtém a instrução de juros
func getInterestInstruction(title models.Title) string {
	dateInterest := datePlusDaysLocalTime(title.ExpireDateTime, title.Fees.Interest.DaysAfterExpirationDate)
	dateInterestFormatted := brDate(dateInterest)

	interestAmountByDayInReal := calculateInterestByDay(title.Fees.Interest.AmountPerDayInCents, title.Fees.Interest.PercentagePerMonth, title.AmountInCents)

	return fmt.Sprintf("APOS %s: JUROS POR DIA DE ATRASO.........R$ %.2f", dateInterestFormatted, interestAmountByDayInReal)
}

func onlyAlphanumerics(str string) string {
	return regexp.MustCompile("[^a-zA-zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç0-9\\s]+").ReplaceAllString(str, "")
}

func onlyAlphabetics(str string) string {
	return regexp.MustCompile("[^a-zA-zÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕáéíóúàèìòùâêîôûãõç\\s]+").ReplaceAllString(str, "")
}

func onlyOneSpace(str string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")
}

func removeAllSpaces(str string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(str, "")
}
