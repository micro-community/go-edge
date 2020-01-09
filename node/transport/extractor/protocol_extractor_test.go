package extractor

import (
	"regexp"
	"testing"
)

var srcbytes = []byte(`<?xml version="1.0" encoding="gb2312"?>
<XML>
<VER>1.0</VER>
<NAME>danny</NAME>
<GENDER>MALE</GENDER>
<TYPE>1</TYPE>
<ADDR>Road.1</ADDR>
<PHONE>400-800-5555</PHONE>
<COMPANY>xxx</COMPANY>
<TIME>2019.12.1-11:11:11</TIME>
</XML>
`)

func TestRegReplaceHeader(t *testing.T) {

	reg, _ := regexp.Compile("(?i:^<\\?xml (.+?)\\?>)")

	indexs := reg.FindIndex(srcbytes)

	t.Log(indexs)
	t.Log(string(srcbytes[0:indexs[0]]))
	t.Log(string(srcbytes[0:indexs[1]]))
	t.Log(string(srcbytes[indexs[0]:indexs[1]]))
	t.Log("Done")

}

func TestRegUNFmatchFooter(t *testing.T) {

	reg, _ := regexp.Compile("(?i:</XML>)")

	targetString := reg.Find(srcbytes)

	t.Log(string(targetString))

	t.Log("Done")

}

func TestRegUNFmatchHeader(t *testing.T) {

	reg, _ := regexp.Compile("(?i:</XML>)")

	indexs := reg.FindIndex(srcbytes)

	t.Log(indexs)
	t.Log(string(srcbytes[0:indexs[0]]))
	t.Log(string(srcbytes[0:indexs[1]]))
	t.Log(string(srcbytes[indexs[0]:indexs[1]]))
	t.Log("Done")

}

func TestRegUNFmatchTypeAndFooter(t *testing.T) {

	reg, _ := regexp.Compile("(?i:<typ>(?i:EVT|RSP)</typ>)")

	indexs := reg.FindIndex(srcbytes)

	t.Log(indexs)
	t.Log(string(srcbytes[0:indexs[0]]))
	t.Log(string(srcbytes[0:indexs[1]]))
	t.Log(string(srcbytes[indexs[0]:indexs[1]]))
	t.Log("Done")

}

func TestRegUNFmatchTypeAndCatch(t *testing.T) {

	reg, _ := regexp.Compile("(?i:<typ>(?i:EVT|RSP)</typ>)")

	indexs := reg.FindIndex(srcbytes)

	t.Log(indexs)
	t.Log(string(srcbytes[0:indexs[0]]))
	t.Log(string(srcbytes[0:indexs[1]]))
	t.Log(string(srcbytes[indexs[0]:indexs[1]]))
	t.Log("Done")

}

func TestRegUNFmatchVersion(t *testing.T) {

	//	targetFormatString := "(?i:<ver>([a-zA-Z0-9\\.]+?)</ver>)"
	targetFormatString := "(?i:<ver>([a-zA-Z0-9\\.]+?)</ver>)"
	reg, _ := regexp.Compile(targetFormatString)

	indexs := reg.FindIndex(srcbytes)

	t.Log(indexs)
	t.Log(reg.SubexpNames())
	t.Log(string(srcbytes[0:indexs[0]]))
	t.Log(string(srcbytes[0:indexs[1]]))
	t.Log(string(srcbytes[indexs[0]:indexs[1]]))
	t.Log("Done")

}
