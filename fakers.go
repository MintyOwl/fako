package fako

import (
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/icrowley/fake"
	"github.com/serenize/snaker"
)

var src = rand.New(&rndSrc{src: rand.NewSource(time.Now().UnixNano())})

// piece of code taken from github.com/icrowley/fake
type rndSrc struct {
	mtx sync.Mutex
	src rand.Source
}

func (s *rndSrc) Int63() int64 {
	s.mtx.Lock()
	n := s.src.Int63()
	s.mtx.Unlock()
	return n
}

func (s *rndSrc) Seed(n int64) {
	s.mtx.Lock()
	s.src.Seed(n)
	s.mtx.Unlock()
}

func Rndm(n int) int {
	var result string
	for i := 0; i < n; i++ {
		if i == 0 {
			first := src.Intn(n)
			if first == 0 {
				result += strconv.Itoa(5)
			} else {
				result += strconv.Itoa(first)
			}
		} else {
			result += strconv.Itoa(src.Intn(n))
		}

	}
	res, _ := strconv.Atoi(result)
	return res

}

var customGenerators = map[string]func() string{}
var typeMappingInt = map[string]func(int) int{
	"Int": Rndm,
}

var typeMapping = map[string]func() string{
	"Brand":                    fake.Brand,
	"Character":                fake.Character,
	"Characters":               fake.Characters,
	"City":                     fake.City,
	"Color":                    fake.Color,
	"Company":                  fake.Company,
	"Continent":                fake.Continent,
	"Country":                  fake.Country,
	"CreditCardType":           fake.CreditCardType,
	"Currency":                 fake.Currency,
	"CurrencyCode":             fake.CurrencyCode,
	"Digits":                   fake.Digits,
	"DomainName":               fake.DomainName,
	"DomainZone":               fake.DomainZone,
	"EmailAddress":             fake.EmailAddress,
	"EmailBody":                fake.EmailBody,
	"EmailSubject":             fake.EmailSubject,
	"FemaleFirstName":          fake.FemaleFirstName,
	"FemaleFullName":           fake.FemaleFullName,
	"FemaleFullNameWithPrefix": fake.FemaleFullNameWithPrefix,
	"FemaleFullNameWithSuffix": fake.FemaleFullNameWithSuffix,
	"FemaleLastName":           fake.FemaleLastName,
	"FemalePatronymic":         fake.FemalePatronymic,
	"FirstName":                fake.FirstName,
	"FullName":                 fake.FullName,
	"FullNameWithPrefix":       fake.FullNameWithPrefix,
	"FullNameWithSuffix":       fake.FullNameWithSuffix,
	"Gender":                   fake.Gender,
	"GenderAbbrev":             fake.GenderAbbrev,
	"HexColor":                 fake.HexColor,
	"HexColorShort":            fake.HexColorShort,
	"IPv4":                     fake.IPv4,
	"Industry":                 fake.Industry,
	"JobTitle":                 fake.JobTitle,
	"Language":                 fake.Language,
	"LastName":                 fake.LastName,
	"LatitudeDirection":        fake.LatitudeDirection,
	"LongitudeDirection":       fake.LongitudeDirection,
	"MaleFirstName":            fake.MaleFirstName,
	"MaleFullName":             fake.MaleFullName,
	"MaleFullNameWithPrefix":   fake.MaleFullNameWithPrefix,
	"MaleFullNameWithSuffix":   fake.MaleFullNameWithSuffix,
	"MaleLastName":             fake.MaleLastName,
	"MalePatronymic":           fake.MalePatronymic,
	"Model":                    fake.Model,
	"Month":                    fake.Month,
	"MonthShort":               fake.MonthShort,
	"Paragraph":                fake.Paragraph,
	"Paragraphs":               fake.Paragraphs,
	"Patronymic":               fake.Patronymic,
	"Phone":                    fake.Phone,
	"Product":                  fake.Product,
	"ProductName":              fake.ProductName,
	"Sentence":                 fake.Sentence,
	"Sentences":                fake.Sentences,
	"SimplePassword":           fake.SimplePassword,
	"State":                    fake.State,
	"StateAbbrev":              fake.StateAbbrev,
	"Street":                   fake.Street,
	"StreetAddress":            fake.StreetAddress,
	"Title":                    fake.Title,
	"TopLevelDomain":           fake.TopLevelDomain,
	"UserName":                 fake.UserName,
	"WeekDay":                  fake.WeekDay,
	"WeekDayShort":             fake.WeekDayShort,
	"Word":                     fake.Word,
	"Words":                    fake.Words,
	"Zip":                      fake.Zip,
}

// Register allows user to add his own data generators for special cases
// that we could not cover with the generators that fako includes by default.
func Register(identifier string, generator func() string) {
	fakeType := snaker.SnakeToCamel(identifier)
	customGenerators[fakeType] = generator
}

// Fuzz Fills passed interface with random data based on the struct field type,
// take a look at fuzzValueFor for details on supported data types.
func Fuzz(e interface{}) {
	ty := reflect.TypeOf(e)

	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}

	if ty.Kind() == reflect.Struct {
		value := reflect.ValueOf(e).Elem()
		for i := 0; i < ty.NumField(); i++ {
			field := value.Field(i)

			if field.CanSet() {
				field.Set(fuzzValueFor(field.Kind()))
			}
		}

	}
}

func allGeneratorsInt() map[string]func(int) int {
	return typeMappingInt

}

func allGenerators() map[string]func() string {
	dst := typeMapping
	for k, v := range customGenerators {
		dst[k] = v
	}

	return dst
}

//findFakeFunctionForInt returns a faker function for a fako identifier
func findFakeFunctionForInt(fako string) func(int) int {
	result := func(int) int { return 123456789 }
	for kind, function := range allGeneratorsInt() {
		if fako == kind {
			result = function
			break
		}
	}
	return result
}

//findFakeFunctionFor returns a faker function for a fako identifier
func findFakeFunctionFor(fako string) func() string {
	result := func() string { return "" }

	for kind, function := range allGenerators() {
		if fako == kind {
			result = function
			break
		}
	}

	return result
}

// fuzzValueFor Generates random values for the following types:
// string, bool, int, int32, int64, float32, float64
func fuzzValueFor(kind reflect.Kind) reflect.Value {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	switch kind {
	case reflect.String:
		return reflect.ValueOf(randomString(25))
	case reflect.Int:
		return reflect.ValueOf(r.Int())
	case reflect.Int32:
		return reflect.ValueOf(r.Int31())
	case reflect.Int64:
		return reflect.ValueOf(r.Int63())
	case reflect.Float32:
		return reflect.ValueOf(r.Float32())
	case reflect.Float64:
		return reflect.ValueOf(r.Float64())
	case reflect.Bool:
		val := r.Intn(2) > 0
		return reflect.ValueOf(val)
	}

	return reflect.ValueOf("")
}
