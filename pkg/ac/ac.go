package ac

type State struct {
	Power bool   `json:"power"`
	Mode  string `json:"mode"`
	Wind  int    `json:"wind"`
	Temp  int    `json:"temp"`
}

const (
	kPowerMask byte = 0x08
	kModeAuto  byte = 0
	kModeCool  byte = 1
	kModeDry   byte = 2
	kModeFan   byte = 3
	kModeHeat  byte = 4
	kTempMin   byte = 16
	kTempMax   byte = 30
)

func (s *State) Code() [8]byte {
	if s.Mode == "auto" {
		s.Temp = 25
	} else if s.Mode == "dry" {
		s.Wind = 1
	}

	code := [8]byte{}
	code[3] = 0x50

	// wind/fan
	wind := min(byte(s.Wind), 3)
	code[0] |= (wind << 4)

	// power
	if s.Power {
		code[0] |= kPowerMask
	} else {
		code[0] &= ^byte(0x08)
	}

	// mode
	var mode byte
	switch s.Mode {
	case "auto":
		mode = kModeAuto
	case "cool":
		mode = kModeCool
	case "dry":
		mode = kModeDry
	case "fan":
		mode = kModeFan
	case "heat":
		mode = kModeHeat
	}
	code[0] |= mode

	// temperature
	temp := byte(s.Temp)
	temp = max(temp, kTempMin)
	temp = min(temp, kTempMax)
	code[1] |= temp - kTempMin

	// checksum
	var sum byte = 10
	// Sum the lower half of the first 4 bytes of this block.
	for i := 0; i < 4; i++ {
		sum += (code[i] & 0x0F)
	}
	// then sum the upper half of the next 4 bytes.
	for i := 4; i < 8; i++ {
		sum += (code[i] >> 4)
	}
	// Trim it down to fit into the 4 bits allowed. i.e. Mod 16.
	code[7] = (sum << 4) | (code[7] & 0x0F)

	return code
}

func max(a, b byte) byte {
	if a > b {
		return a
	}

	return b
}

func min(a, b byte) byte {
	if a > b {
		return b
	}

	return a
}
