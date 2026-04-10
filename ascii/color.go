package ascii

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var Colors = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"orange": "\033[38;2;255;165;0m",
	"reset":  "\033[0m",
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	for h < 0 {
		h += 360
	}
	for h >= 360 {
		h -= 360
	}

	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hueToRGB(p, q, h/360.0+1.0/3.0)
		g = hueToRGB(p, q, h/360.0)
		b = hueToRGB(p, q, h/360.0-1.0/3.0)
	}

	return uint8(math.Round(r * 255)), uint8(math.Round(g * 255)), uint8(math.Round(b * 255))
}

func ParseColor(color string) (string, bool) {
	color = strings.ToLower(strings.TrimSpace(color))
	
	if val, ok := Colors[color]; ok {
		return val, true
	}

	if strings.HasPrefix(color, "hsl(") && strings.HasSuffix(color, ")") {
		inner := color[4 : len(color)-1]
		parts := strings.Split(inner, ",")
		if len(parts) == 3 {
			hStr := strings.TrimSpace(parts[0])
			sOrig := strings.TrimSpace(parts[1])
			lOrig := strings.TrimSpace(parts[2])
			
			sStr := strings.TrimSuffix(sOrig, "%")
			lStr := strings.TrimSuffix(lOrig, "%")

			h, errH := strconv.ParseFloat(hStr, 64)
			s, errS := strconv.ParseFloat(sStr, 64)
			l, errL := strconv.ParseFloat(lStr, 64)

			if errH == nil && errS == nil && errL == nil {
				if strings.HasSuffix(sOrig, "%") || s > 1 {
					s /= 100.0
				}
				if strings.HasSuffix(lOrig, "%") || l > 1 {
					l /= 100.0
				}
				
				if s < 0 { s = 0 }
				if s > 1 { s = 1 }
				if l < 0 { l = 0 }
				if l > 1 { l = 1 }

				r, g, b := hslToRGB(h, s, l)
				return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), true
			}
		}
	}

	return "", false
}

// GetCharColors computes the ANSI color codes for each character in the given text
func GetCharColors(text string, rules map[string]string) []string {
	charColors := make([]string, len(text))
	if len(rules) == 0 {
		return charColors
	}

	for sub, colorName := range rules {
		colorCode, ok := ParseColor(colorName)
		if !ok {
			continue // ignore invalid colors
		}

		// empty key means color everything
		if sub == "" {
			for i := range charColors {
				charColors[i] = colorCode
			}
			continue
		}

		for i := 0; i <= len(text)-len(sub); i++ {
			if text[i:i+len(sub)] == sub {
				for j := 0; j < len(sub); j++ {
					charColors[i+j] = colorCode
				}
			}
		}
	}
	return charColors
}
