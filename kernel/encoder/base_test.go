package encoder_test

import (
	"fmt"

	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
)

var (
	_ = encoder.EncodeToStringFunc(encoder.Base16Encode)           // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base32Encode)           // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base32EscapedEncode)    // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base32HexEncode)        // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base32HexEscapedEncode) // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base64Encode)           // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base64RawEncode)        // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base64URLEncode)        // Ensure that the function satisfy the signature.
	_ = encoder.EncodeToStringFunc(encoder.Base64RawURLEncode)     // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base16Decode)             // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base32Decode)             // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base32EscapedDecode)      // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base32HexDecode)          // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base32HexEscapedDecode)   // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base64Decode)             // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base64RawDecode)          // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base64URLDecode)          // Ensure that the function satisfy the signature.
	_ = encoder.DecodeStringFunc(encoder.Base64RawURLDecode)       // Ensure that the function satisfy the signature.
)

func ExampleBase16Encode() {
	data := []byte("encoding example")
	encoded := encoder.Base16Encode(data)

	fmt.Println(encoded)
	// Output:
	// 656e636f64696e67206578616d706c65
}

func ExampleBase16Decode() {
	encoded := "656e636f64696e67206578616d706c65"
	decoded, err := encoder.Base16Decode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase32Encode() {
	data := []byte("encoding example")
	encoded := encoder.Base32Encode(data)

	fmt.Println(encoded)
	// Output:
	// MVXGG33ENFXGOIDFPBQW24DMMU======
}

func ExampleBase32Decode() {
	encoded := "MVXGG33ENFXGOIDFPBQW24DMMU======"
	decoded, err := encoder.Base32Decode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase32EscapedEncode() {
	data := []byte("encoding example")
	encoded := encoder.Base32EscapedEncode(data)

	fmt.Println(encoded)
	// Output:
	// QZ1JJ55GRH1JSLFHTCU046FQQY======
}

func ExampleBase32EscapedDecode() {
	encoded := "QZ1JJ55GRH1JSLFHTCU046FQQY======"
	decoded, err := encoder.Base32EscapedDecode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase32HexEncode() {
	data := []byte("encoding example")
	encoded := encoder.Base32HexEncode(data)

	fmt.Println(encoded)
	// Output:
	// CLN66RR4D5N6E835F1GMQS3CCK======
}

func ExampleBase32HexDecode() {
	encoded := "CLN66RR4D5N6E835F1GMQS3CCK======"
	decoded, err := encoder.Base32HexDecode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase32HexEscapedEncode() {
	data := []byte("encoding example")
	encoded := encoder.Base32HexEscapedEncode(data)

	fmt.Println(encoded)
	// Output:
	// DPR66VV4F5R6G835H1JQUW3DDN======
}

func ExampleBase32HexEscapedDecode() {
	encoded := "DPR66VV4F5R6G835H1JQUW3DDN======"
	decoded, err := encoder.Base32HexEscapedDecode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase64Encode() {
	data := []byte("encoding example")
	encoded := encoder.Base64Encode(data)

	fmt.Println(encoded)
	// Output:
	// ZW5jb2RpbmcgZXhhbXBsZQ==
}

func ExampleBase64Decode() {
	encoded := "ZW5jb2RpbmcgZXhhbXBsZQ=="
	decoded, err := encoder.Base64Decode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase64RawEncode() {
	data := []byte("encoding example")
	encoded := encoder.Base64RawEncode(data)

	fmt.Println(encoded)
	// Output:
	// ZW5jb2RpbmcgZXhhbXBsZQ
}

func ExampleBase64RawDecode() {
	encoded := "ZW5jb2RpbmcgZXhhbXBsZQ"
	decoded, err := encoder.Base64RawDecode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase64URLEncode() {
	data := []byte("encoding example")
	encoded := encoder.Base64URLEncode(data)

	fmt.Println(encoded)
	// Output:
	// ZW5jb2RpbmcgZXhhbXBsZQ==
}

func ExampleBase64URLDecode() {
	encoded := "ZW5jb2RpbmcgZXhhbXBsZQ=="
	decoded, err := encoder.Base64URLDecode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}

func ExampleBase64RawURLEncode() {
	data := []byte("encoding example")
	encoded := encoder.Base64RawURLEncode(data)

	fmt.Println(encoded)
	// Output:
	// ZW5jb2RpbmcgZXhhbXBsZQ
}

func ExampleBase64RawURLDecode() {
	encoded := "ZW5jb2RpbmcgZXhhbXBsZQ"
	decoded, err := encoder.Base64RawURLDecode(encoded)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(decoded))
	// Output:
	// encoding example
}
