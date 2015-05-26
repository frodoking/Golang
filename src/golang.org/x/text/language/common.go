// Generated by running
//		maketables -url=http://www.unicode.org/Public/cldr/26/core.zip -iana=http://www.iana.org/assignments/language-subtag-registry -tld=http://www.iana.org/domains/root/db
// automatically with go generate.
// DO NOT EDIT

package language

// This file contains code common to the maketables.go and the package code.

const (
    curDigitBits = 3
    curDigitMask = 1<<curDigitBits - 1
    curRoundBits = 0 // Appear to be always zero.
)

type currencyInfo int

func mkCurrencyInfo(round, decimal int) string {
    return string([]byte{byte(round<<curDigitBits | decimal)})
}

func (c currencyInfo) round() int {
    return int(c >> curDigitBits)
}

func (c currencyInfo) decimals() int {
    return int(c & curDigitMask)
}

// langAliasType is the type of an alias in langAliasMap.
type langAliasType int8

const (
    langDeprecated langAliasType = iota
    langMacro
    langLegacy

    langAliasTypeUnknown langAliasType = -1
)
