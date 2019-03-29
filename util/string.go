package util

import (
	"strings"
)

const (
	LANG_CN = "zh-CN"
	LANG_EN = "en-US"
)

const (
	RUNE_ZH_BEGIN     = 19968
	RUNE_ZH_END       = 65535
	RUNE_EN_BEGIN     = 'A'
	RUNE_EN_END       = 'z'
	RUNE_NUMBER_BEGIN = '0'
	RUNE_NUMBER_END   = '9'
)

type LangString struct {
	Str string `xml:"sub_text"`
	Lan string `xml:"language"`
}

type LangStrings []*LangString

func (ls *LangString) Trim(c string) {
	ls.Str = strings.Trim(ls.Str, c)
}

type TTSTextWhiteListItem struct {
	Tex             string       `xml:"text"`
	SubStrs         []LangString `xml:"sub_text_lang"`
	AdaptiveContext bool         `xml:"adaptive_context"`
}

type TTSTextWhiteList struct {
	List []TTSTextWhiteListItem `xml:"white_list_item"`
}

func SplitStringByLanguage(str string) LangStrings {
	lastLang := ""
	isLangChange := false
	curLang := ""
	curIndex := 0
	strs := make(LangStrings, 0)
	rs := []rune(str)
	strNumberHeader := ""
	for _, r := range rs {
		if IsRuneZH(r) {
			curLang = LANG_CN
		} else if IsRuneEN(r) {
			curLang = LANG_EN
		} else {
			if "" == lastLang {
				if IsNumber(r) || "" != strNumberHeader {
					strNumberHeader += string(r)
				}
				continue
			}
		}
		if "" == lastLang {
			lastLang = curLang
			ls := new(LangString)
			ls.Lan = curLang
			ls.Str = strNumberHeader
			strs = append(strs, ls)
		}
		if lastLang != curLang {
			isLangChange = true
		} else {
			isLangChange = false
		}
		if isLangChange {
			ls := new(LangString)
			ls.Lan = curLang
			strs = append(strs, ls)
			curIndex++
		}
		lastLang = curLang
		strs[curIndex].Str += string(r)
	}
	if 0 == len(strs) && "" != strNumberHeader {
		ls := new(LangString)
		ls.Lan = LANG_CN
		ls.Str = strNumberHeader
		strs = append(strs, ls)
	}
	for _, ls := range strs {
		ls.Trim(" ")
	}
	return strs
}

func IsRuneZH(r rune) bool {
	return r <= RUNE_ZH_END && r >= RUNE_ZH_BEGIN
}

func IsRuneEN(r rune) bool {
	return r <= RUNE_EN_END && r >= RUNE_EN_BEGIN
}

func IsNumber(r rune) bool {
	return r <= RUNE_NUMBER_END && r >= RUNE_NUMBER_BEGIN
}

func IsStringZH(str string, ignoreNumber bool) bool {
	if "" == str {
		return false
	}
	rs := []rune(str)
	for _, r := range rs {
		if IsRuneZH(r) {
			continue
		}
		if ignoreNumber && IsNumber(r) {
			continue
		} else {
			return false
		}
		if IsRuneEN(r) {
			return false
		} else {
			continue
		}
	}
	return true
}

func IsStringEN(str string, ignoreNumber bool) bool {
	if "" == str {
		return false
	}
	rs := []rune(str)
	for _, r := range rs {
		if IsRuneEN(r) {
			continue
		}
		if ignoreNumber && IsNumber(r) {
			continue
		} else {
			return false
		}
		if IsRuneZH(r) {
			return false
		} else {
			continue
		}
	}
	return true
}

func SplitStringByLanguageAndWhiteList(str string, whiteList TTSTextWhiteList) LangStrings {
	strs := SplitStringByLanguage(str)
	if nil != whiteList.List && 0 < len(whiteList.List) && len(strs) > 1 {
		for _, item := range whiteList.List {
			if strings.Contains(strings.ToLower(str), item.Tex) {
				strs = ConvertLangStrsByWhiteListItem(strs, &item)
			}
		}
	}
	for _, ls := range strs {
		ls.Trim(" ")
	}
	return strs
}

func ConvertLangStrsByWhiteListItem(langStrings LangStrings, item *TTSTextWhiteListItem) LangStrings {
	strsLen := len(langStrings)
	if strsLen == 0 {
		return nil
	}
	outStrs := make(LangStrings, 2*strsLen)
	candidate := make(LangStrings, strsLen)
	outLen := len(outStrs)
	outIndex := 0
	lastLang := ""
	for i := 0; i < strsLen; i++ {
		if 0 == i {
			outStrs[outIndex] = langStrings[i]
			lastLang = langStrings[i].Lan
			outIndex++
			continue
		}
		lowStr := strings.ToLower(langStrings[i].Str)
		index := strings.Index(lowStr, item.Tex)
		if 0 == index && LANG_CN == lastLang {
			for _, val := range item.SubStrs {
				ls := new(LangString)
				ls.Str = val.Str
				ls.Lan = val.Lan
				outStrs[outIndex] = ls
				outIndex++
				if outIndex >= outLen {
					outStrs = append(outStrs, candidate...)
					outLen = len(outStrs)
				}
			}
			if lowStr != item.Tex {
				ls := new(LangString)
				ls.Str = langStrings[i].Str[len(item.Tex):]
				ls.Lan = langStrings[i].Lan
				outStrs[outIndex] = ls
				outIndex++
				if outIndex >= outLen {
					outStrs = append(outStrs, candidate...)
					outLen = len(outStrs)
				}
			}
		} else {
			if outStrs[outIndex-1].Lan == langStrings[i].Lan {
				outStrs[outIndex-1].Str += (" " + langStrings[i].Str)
				continue
			} else {
				ls := new(LangString)
				ls.Str = langStrings[i].Str
				ls.Lan = langStrings[i].Lan
				outStrs[outIndex] = ls
				outIndex++
				if outIndex >= outLen {
					outStrs = append(outStrs, candidate...)
					outLen = len(outStrs)
				}
			}
		}
		lastLang = outStrs[outIndex-1].Lan
	}
	return outStrs[:outIndex]
}
