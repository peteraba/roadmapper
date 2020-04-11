package main

import (
	"errors"
	"fmt"
	"image/color"
)

/**
https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717
https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54
https://coolors.co/app/0a2463-3e92cc-fffaff-d8315b-1e1b18
*/
var colors = [][]color.RGBA{
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #2
	{
		mustParseColor("#18000B"),
		mustParseColor("#2F0015"),
		mustParseColor("#460020"),
		mustParseColor("#5D002A"),
		mustParseColor("#740034"),
		mustParseColor("#8C003F"),
		mustParseColor("#A30049"),
		mustParseColor("#BA0053"),
		mustParseColor("#D1005E"),
		mustParseColor("#E80068"),
		mustParseColor("#FF0072"), // 0
		mustParseColor("#FF177E"),
		mustParseColor("#FF177E"),
		mustParseColor("#FF4598"),
		mustParseColor("#FF5CA5"),
		mustParseColor("#FF73B2"),
		mustParseColor("#FF8BBE"),
		mustParseColor("#FFA2CB"),
		mustParseColor("#FFB9D8"),
		mustParseColor("#FFD0E5"),
		// mustParseColor("#FFE7F2"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #3
	{
		mustParseColor("#01140F"),
		mustParseColor("#02271E"),
		mustParseColor("#023B2C"),
		mustParseColor("#034E3B"),
		mustParseColor("#036249"),
		mustParseColor("#047558"),
		mustParseColor("#048966"),
		mustParseColor("#059C75"),
		mustParseColor("#05B083"),
		mustParseColor("#06C392"),
		mustParseColor("#06D6A0"), // 0
		mustParseColor("#1CD9A8"),
		mustParseColor("#33DDB1"),
		mustParseColor("#49E1B9"),
		mustParseColor("#60E4C2"),
		mustParseColor("#77E8CB"),
		mustParseColor("#77E8CB"),
		mustParseColor("#8DECD3"),
		mustParseColor("#A4F0DC"),
		mustParseColor("#BBF3E5"),
		mustParseColor("#D1F7ED"),
		// mustParseColor("#E8FBF6"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #1
	{
		mustParseColor("#181200"),
		mustParseColor("#2F2300"),
		mustParseColor("#463500"),
		mustParseColor("#5D4600"),
		mustParseColor("#745700"),
		mustParseColor("#8C6900"),
		mustParseColor("#A37A00"),
		mustParseColor("#BA8B00"),
		mustParseColor("#D19D00"),
		mustParseColor("#E8AE00"),
		mustParseColor("#FFBF00"), // 0
		mustParseColor("#FFC417"),
		mustParseColor("#FFCA2E"),
		mustParseColor("#FFD045"),
		mustParseColor("#FFD65C"),
		mustParseColor("#FFDC73"),
		mustParseColor("#FFE18B"),
		mustParseColor("#FFE7A2"),
		mustParseColor("#FFEDB9"),
		mustParseColor("#FFF3D0"),
		// mustParseColor("#FFF9E7"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #4
	{
		mustParseColor("#120017"),
		mustParseColor("#24002E"),
		mustParseColor("#350045"),
		mustParseColor("#47005C"),
		mustParseColor("#590073"),
		mustParseColor("#6A0089"),
		mustParseColor("#7C00A0"),
		mustParseColor("#8E00B7"),
		mustParseColor("#9F00CE"),
		mustParseColor("#B100E5"),
		mustParseColor("#C200FB"), // 0
		mustParseColor("#C717FB"),
		mustParseColor("#CD2EFB"),
		mustParseColor("#D245FC"),
		mustParseColor("#D85CFC"),
		mustParseColor("#DD73FC"),
		mustParseColor("#E38BFD"),
		mustParseColor("#E8A2FD"),
		mustParseColor("#EEB9FD"),
		mustParseColor("#F3D0FE"),
		// mustParseColor("#F9E7FE"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #5
	{
		mustParseColor("#030303"),
		mustParseColor("#050505"),
		mustParseColor("#070707"),
		mustParseColor("#090909"),
		mustParseColor("#0B0B0B"),
		mustParseColor("#0D0D0D"),
		mustParseColor("#0F0F0F"),
		mustParseColor("#111111"),
		mustParseColor("#131313"),
		mustParseColor("#151515"),
		mustParseColor("#171717"), // 0
		mustParseColor("#2C2C2C"),
		mustParseColor("#414141"),
		mustParseColor("#565656"),
		mustParseColor("#6B6B6B"),
		mustParseColor("#808080"),
		mustParseColor("#959595"),
		mustParseColor("#AAAAAA"),
		mustParseColor("#BFBFBF"),
		mustParseColor("#D4D4D4"),
		// mustParseColor("#E9E9E9"),
	},
}

func pickFgColor(epicCount, projectCount, indentation int) *color.RGBA {
	c := colors[epicCount%len(colors)]

	switch indentation {
	case 0:
		return &c[len(c)/2]

	default:
		n := (len(c)/2 - projectCount - (indentation)*5) % len(c)
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
