package colors

import (
	"errors"
	"fmt"
	"image/color"
	"strings"
)

/**
https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717
https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54
*/
var colors = [][]color.RGBA{
	// https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54 # 3 - VIVID CERULEAN
	{
		mustParseColor("#001016"),
		mustParseColor("#001F2C"),
		mustParseColor("#002E41"),
		mustParseColor("#003D57"),
		mustParseColor("#004C6C"),
		mustParseColor("#005B82"),
		mustParseColor("#006A97"),
		mustParseColor("#0079AD"),
		mustParseColor("#0088C2"),
		mustParseColor("#0097D8"),
		mustParseColor("#00A6ED"), // 0
		mustParseColor("#17AEEE"),
		mustParseColor("#2EB6F0"),
		mustParseColor("#45BEF1"),
		mustParseColor("#5CC6F3"),
		mustParseColor("#73CEF5"),
		mustParseColor("#8BD6F6"),
		mustParseColor("#A2DEF8"),
		mustParseColor("#B9E6FA"),
		mustParseColor("#D0EEFB"),
		// mustParseColor("#E7F6FD"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #2 - VIVID RASPBERRY
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
		mustParseColor("#FF2E8B"),
		mustParseColor("#FF4598"),
		mustParseColor("#FF5CA5"),
		mustParseColor("#FF73B2"),
		mustParseColor("#FF8BBE"),
		mustParseColor("#FFA2CB"),
		mustParseColor("#FFB9D8"),
		mustParseColor("#FFD0E5"),
		// mustParseColor("#FFE7F2"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #1 - FLUORESCENT ORANGE
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
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #3 - CARIBBEAN GREEN
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
		mustParseColor("#8DECD3"),
		mustParseColor("#A4F0DC"),
		mustParseColor("#BBF3E5"),
		mustParseColor("#D1F7ED"),
		// mustParseColor("#E8FBF6"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #4 - ELECTRIC PURPLE
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
	// https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54 # 1 - ORIOLES ORANGE
	{
		mustParseColor("#170803"),
		mustParseColor("#2D0F06"),
		mustParseColor("#441708"),
		mustParseColor("#5A1E0B"),
		mustParseColor("#70250E"),
		mustParseColor("#872D10"),
		mustParseColor("#9D3413"),
		mustParseColor("#B33B16"),
		mustParseColor("#CA4318"),
		mustParseColor("#E04A1B"),
		mustParseColor("#F6511D"), // 0
		mustParseColor("#F66031"),
		mustParseColor("#F77046"),
		mustParseColor("#F8805A"),
		mustParseColor("#F9906F"),
		mustParseColor("#FAA083"),
		mustParseColor("#FAAF98"),
		mustParseColor("#FBBFAC"),
		mustParseColor("#FCCFC1"),
		mustParseColor("#FDDFD5"),
		// mustParseColor("#FEEFEA"),
	},
	// https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54 # 5 - PRUSSIAN BLUE
	{
		mustParseColor("#020408"),
		mustParseColor("#030810"),
		mustParseColor("#040C17"),
		mustParseColor("#05101F"),
		mustParseColor("#061427"),
		mustParseColor("#08182E"),
		mustParseColor("#091D36"),
		mustParseColor("#0A213E"),
		mustParseColor("#0B2545"),
		mustParseColor("#0C284D"),
		mustParseColor("#0D2C54"), // 0
		mustParseColor("#233F63"),
		mustParseColor("#395273"),
		mustParseColor("#4F6582"),
		mustParseColor("#657892"),
		mustParseColor("#7B8BA1"),
		mustParseColor("#919FB1"),
		mustParseColor("#A7B2C0"),
		mustParseColor("#BDC5D0"),
		mustParseColor("#D3D8DF"),
		// mustParseColor("#E9EBEF"),
	},
	// https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54 # 2 - UCLA GOLD
	{
		mustParseColor("#181100"),
		mustParseColor("#2F2100"),
		mustParseColor("#463200"),
		mustParseColor("#5D4200"),
		mustParseColor("#745200"),
		mustParseColor("#8C6300"),
		mustParseColor("#A37300"),
		mustParseColor("#BA8300"),
		mustParseColor("#D19400"),
		mustParseColor("#E8A400"),
		mustParseColor("#FFB400"), // 0
		mustParseColor("#FFBA17"),
		mustParseColor("#FFC12E"),
		mustParseColor("#FFC845"),
		mustParseColor("#FFCF5C"),
		mustParseColor("#FFD673"),
		mustParseColor("#FFDC8B"),
		mustParseColor("#FFE3A2"),
		mustParseColor("#FFEAB9"),
		mustParseColor("#FFF1D0"),
		// mustParseColor("#FFF8E7"),
	},
	// https://coolors.co/f6511d-ffb400-00a6ed-7fb800-0d2c54 # 4 - APPLE GREEN
	{
		mustParseColor("#0C1100"),
		mustParseColor("#182200"),
		mustParseColor("#233300"),
		mustParseColor("#2F4300"),
		mustParseColor("#3A5400"),
		mustParseColor("#466500"),
		mustParseColor("#517600"),
		mustParseColor("#5D8600"),
		mustParseColor("#689700"),
		mustParseColor("#74A800"),
		mustParseColor("#7FB800"), // 0
		mustParseColor("#8ABE17"),
		mustParseColor("#96C42E"),
		mustParseColor("#A1CB45"),
		mustParseColor("#ADD15C"),
		mustParseColor("#B9D873"),
		mustParseColor("#C4DE8B"),
		mustParseColor("#D0E5A2"),
		mustParseColor("#DCEBB9"),
		mustParseColor("#E7F2D0"),
		// mustParseColor("#F3F8E7"),
	},
	// https://coolors.co/ffbf00-ff0072-06d6a0-c200fb-171717 #5 - EERIE BLACK
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

// PickFgColor will pick a color for a project based on
// the number of epic, task and indentation
func PickFgColor(epicCount, taskCount, indentation int) *color.RGBA {
	c := colors[epicCount%len(colors)]

	switch indentation {
	case 0:
		return &c[len(c)/2]

	case 1:
		n := (len(c)/2 - taskCount*2) % len(c)
		if n < 0 {
			n += len(c)
		}
		return &c[n]

	default:
		n := (len(c)/2 - taskCount*2 + (indentation)*5) % len(c)
		if n < 0 {
			n += len(c)
		}
		return &c[n]
	}
}

// PickBgColor will pick a color for main project
func PickBgColor(epicCount int) color.RGBA {
	c := colors[epicCount%len(colors)]

	return c[len(c)-1]
}

// mustParseColor turns parses a string as a hexadecimal color
// it will panic in case it is not possible
// meant to be used for hardcoded colors...
func mustParseColor(part string) color.RGBA {
	if len(part) != 4 && len(part) != 7 {
		panic(fmt.Errorf("invalid hexa color length: %d", len(part)))
	}

	if part[0] != '#' {
		panic(fmt.Errorf("invalid first character of hexa color. want: #, got: %c", part[0]))
	}

	s, err := CharsToUint8(part[1:])
	if err != nil {
		panic(fmt.Errorf("failed to parse to uint8s: %w", err))
	}

	return color.RGBA{R: s[0], G: s[1], B: s[2], A: 255}
}

// CharsToUint8 converts a string containing a hexadecimal representation of a color into decimal representation
func CharsToUint8(part string) ([3]uint8, error) {
	if len(part) != 3 && len(part) != 6 {
		return [3]uint8{}, errors.New("invalid hexadecimal color string")
	}

	part = strings.ToLower(part)

	tmp := []int{}
	for _, runeValue := range part {
		if idx := strings.IndexRune("0123456789abcdef", runeValue); idx > -1 {
			tmp = append(tmp, idx)
			if len(part) == 3 {
				tmp = append(tmp, idx)
			}
		}
	}

	res := [3]uint8{}
	res[0] = uint8(tmp[0]*16 + tmp[1])
	res[1] = uint8(tmp[2]*16 + tmp[3])
	res[2] = uint8(tmp[4]*16 + tmp[5])

	return res, nil
}

// ToHexa converts a color into a hexadecimal string representation (e.g. #fa3, #ffaa33)
func ToHexa(c color.Color) string {
	if c == nil {
		return ""
	}

	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%s%s%s", twoDigitHexa(r), twoDigitHexa(g), twoDigitHexa(b))
}

// twoDigitHexa converts a number into a hexadecimal representation of a string
func twoDigitHexa(i uint32) string {
	if i > 0xf {
		return fmt.Sprintf("%x", uint8(i))
	}

	return fmt.Sprintf("0%x", uint8(i))
}
