package main

import (
	"errors"
	"fmt"
	"image/color"
)

var colors = [][]color.RGBA{
	{
		mustParseColor("#230B02"),
		mustParseColor("#230B02"),
		mustParseColor("#230B02"),
		mustParseColor("#561A04"),
		mustParseColor("#672005"),
		mustParseColor("#782506"),
		mustParseColor("#892A06"),
		mustParseColor("#9A2F07"),
		mustParseColor("#AB3408"),
		mustParseColor("#BC3908"),
		mustParseColor("#C24B1E"),
		mustParseColor("#C85D34"),
		mustParseColor("#CE6F4B"),
		mustParseColor("#D48161"),
		mustParseColor("#DA9378"),
		mustParseColor("#E0A58E"),
		mustParseColor("#E6B7A5"),
		mustParseColor("#ECC9BB"),
		mustParseColor("#F2DBD2"),
	},
	{
		mustParseColor("#1D1923"),
		mustParseColor("#2C2535"),
		mustParseColor("#3A3146"),
		mustParseColor("#493D58"),
		mustParseColor("#574A69"),
		mustParseColor("#66567B"),
		mustParseColor("#74628C"),
		mustParseColor("#836E9E"),
		mustParseColor("#917AAF"),
		mustParseColor("#9F86C0"),
		mustParseColor("#A791C5"),
		mustParseColor("#B09CCB"),
		mustParseColor("#B9A7D1"),
		mustParseColor("#C1B2D6"),
		mustParseColor("#CABDDC"),
		mustParseColor("#D3C8E2"),
		mustParseColor("#DCD3E8"),
		mustParseColor("#E4DEED"),
		mustParseColor("#EDE9F3"),
	},
	{
		mustParseColor("#0B0F14"),
		mustParseColor("#10161E"),
		mustParseColor("#161E27"),
		mustParseColor("#1B2531"),
		mustParseColor("#202C3B"),
		mustParseColor("#253345"),
		mustParseColor("#2B3B4E"),
		mustParseColor("#304258"),
		mustParseColor("#354962"),
		mustParseColor("#3A506B"),
		mustParseColor("#4B5F78"),
		mustParseColor("#5D6F85"),
		mustParseColor("#6F7F93"),
		mustParseColor("#818FA0"),
		mustParseColor("#939FAE"),
		mustParseColor("#A5AFBB"),
		mustParseColor("#B7BFC9"),
		mustParseColor("#C9CFD6"),
		mustParseColor("#DBDFE4"),
	},
	{
		mustParseColor("#070201"),
		mustParseColor("#0A0301"),
		mustParseColor("#0D0401"),
		mustParseColor("#100501"),
		mustParseColor("#130501"),
		mustParseColor("#160601"),
		mustParseColor("#190701"),
		mustParseColor("#1C0801"),
		mustParseColor("#1F0901"),
		mustParseColor("#220901"),
		mustParseColor("#361F18"),
		mustParseColor("#4A352F"),
		mustParseColor("#5E4C46"),
		mustParseColor("#72625D"),
		mustParseColor("#867874"),
		mustParseColor("#9A8F8B"),
		mustParseColor("#AEA5A2"),
		mustParseColor("#C2BBB9"),
		mustParseColor("#D6D2D0"),
	},
	{
		mustParseColor("#1B0503"),
		mustParseColor("#290804"),
		mustParseColor("#360A05"),
		mustParseColor("#440D06"),
		mustParseColor("#510F07"),
		mustParseColor("#5F1208"),
		mustParseColor("#6C1409"),
		mustParseColor("#7A170A"),
		mustParseColor("#87190B"),
		mustParseColor("#941B0C"),
		mustParseColor("#9D2F22"),
		mustParseColor("#A74438"),
		mustParseColor("#B1594E"),
		mustParseColor("#BA6D64"),
		mustParseColor("#C4827A"),
		mustParseColor("#CE9790"),
		mustParseColor("#D8ACA6"),
		mustParseColor("#E1C0BC"),
		mustParseColor("#EBD5D2"),
	},
	{
		mustParseColor("#C24B1E"),
		mustParseColor("#C24B1E"),
		mustParseColor("#5A3E0B"),
		mustParseColor("#5A3E0B"),
		mustParseColor("#875D10"),
		mustParseColor("#9D6D12"),
		mustParseColor("#B37C15"),
		mustParseColor("#CA8C17"),
		mustParseColor("#E09B1A"),
		mustParseColor("#E09B1A"),
		mustParseColor("#F6B130"),
		mustParseColor("#F7B945"),
		mustParseColor("#F8C159"),
		mustParseColor("#F9C86E"),
		mustParseColor("#FAD083"),
		mustParseColor("#FAD897"),
		mustParseColor("#FBE0AC"),
		mustParseColor("#FCE7C1"),
		mustParseColor("#FDEFD5"),
	},
	{
		mustParseColor("#120502"),
		mustParseColor("#1B0703"),
		mustParseColor("#240903"),
		mustParseColor("#2D0B04"),
		mustParseColor("#360D05"),
		mustParseColor("#3F0F06"),
		mustParseColor("#481106"),
		mustParseColor("#511307"),
		mustParseColor("#5A1508"),
		mustParseColor("#621708"),
		mustParseColor("#702C1E"),
		mustParseColor("#7E4134"),
		mustParseColor("#7E4134"),
		mustParseColor("#9B6B61"),
		mustParseColor("#A98078"),
		mustParseColor("#B7958E"),
		mustParseColor("#C5AAA5"),
		mustParseColor("#D4BFBB"),
		mustParseColor("#E2D4D2"),
	},
	{
		mustParseColor("#0F0F0F"),
		mustParseColor("#161716"),
		mustParseColor("#1E1E1D"),
		mustParseColor("#252524"),
		mustParseColor("#2C2D2C"),
		mustParseColor("#333433"),
		mustParseColor("#3B3B3A"),
		mustParseColor("#424341"),
		mustParseColor("#494A48"),
		mustParseColor("#50514F"),
		mustParseColor("#5F605F"),
		mustParseColor("#6F706F"),
		mustParseColor("#7F807F"),
		mustParseColor("#8F908F"),
		mustParseColor("#9FA09F"),
		mustParseColor("#AFAFAF"),
		mustParseColor("#BFBFBF"),
		mustParseColor("#CFCFCF"),
		mustParseColor("#DFDFDF"),
	},
	{
		mustParseColor("#2C1211"),
		mustParseColor("#421A1A"),
		mustParseColor("#582322"),
		mustParseColor("#6E2C2A"),
		mustParseColor("#843433"),
		mustParseColor("#9B3D3B"),
		mustParseColor("#B14643"),
		mustParseColor("#C74E4C"),
		mustParseColor("#DC5754"),
		mustParseColor("#F25F5C"),
		mustParseColor("#F36D6A"),
		mustParseColor("#F47C79"),
		mustParseColor("#F58A88"),
		mustParseColor("#F69997"),
		mustParseColor("#F7A7A6"),
		mustParseColor("#F9B6B4"),
		mustParseColor("#FAC4C3"),
		mustParseColor("#FBD3D2"),
		mustParseColor("#FCE1E1"),
	},
	{
		mustParseColor("#112323"),
		mustParseColor("#193534"),
		mustParseColor("#224646"),
		mustParseColor("#2A5857"),
		mustParseColor("#326968"),
		mustParseColor("#3A7B79"),
		mustParseColor("#438C8B"),
		mustParseColor("#4B9E9C"),
		mustParseColor("#53AFAD"),
		mustParseColor("#5BC0BE"),
		mustParseColor("#87D1CF"),
		mustParseColor("#96D6D5"),
		mustParseColor("#A5DCDB"),
		mustParseColor("#B4E2E1"),
		mustParseColor("#C3E8E7"),
		mustParseColor("#D2EDED"),
		mustParseColor("#E1F3F3"),
	},
}

func pickFgColor(epicCount, projectCount, indentation int) *color.RGBA {
	c := colors[epicCount%len(colors)]

	switch indentation {
	case 0:
		return &c[len(c)/2]

	default:
		ind := indentation - 1
		n := (len(c)/2 - projectCount - ind*5) % len(c)
		if n < 0 {
			n += len(c)
		}
		return &c[n]
	}
}

func pickBgColor(epic int) color.RGBA {
	c := colors[epic%len(colors)]

	return c[len(c)-1]
}

func mustParseColor(part string) color.RGBA {
	if len(part) != 4 && len(part) != 7 {
		panic(errors.New("invalid hexa color length"))
	}

	if part[0] != '#' {
		panic(errors.New("invalid first character for hexa color"))
	}

	s, err := charsToUint8(part[1:])
	if err != nil {
		panic(fmt.Errorf("failed to parse to uint8s: %w", err))
	}

	return color.RGBA{R: s[0], G: s[1], B: s[2], A: 255}
}
