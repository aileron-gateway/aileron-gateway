package encoder_test

import (
	"regexp"
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
)

var (
	regexpBase16           = regexp.MustCompile(`^[0-9a-fA-F]{2,}$`)
	regexpBase32           = regexp.MustCompile(`^[2-7A-Z].*$`)
	regexpBase32Hex        = regexp.MustCompile(`^[0-9A-V].*$`)
	regexpBase32Escaped    = regexp.MustCompile(`^[0-9A-Z^AEIO].*$`)
	regexpBase32HexEscaped = regexp.MustCompile(`^[0-9A-Z^AEIO].*$`)
	regexpBase64           = regexp.MustCompile(`^[0-9a-zA-Z+/].*$`)
	regexpBase64URL        = regexp.MustCompile(`^[0-9a-zA-Z-_].*$`)
)

func FuzzBase16Decode(f *testing.F) {
	seeds := []string{"", " ", "\n"}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, src string) {
		_, err := encoder.Base16Decode(src)
		if err == nil {
			if regexpBase16.MatchString(src) || src == "" {
				return
			}
			t.Errorf("src: %q", src)
		}
	})
}
func FuzzBase16Encode(f *testing.F) {
	seeds := [][]byte{}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, src []byte) {
		encoder.Base16Encode(src)
	})
}
