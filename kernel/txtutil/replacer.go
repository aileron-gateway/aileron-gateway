// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil

import (
	"bytes"
	"cmp"
	"encoding/hex"
	"regexp"
	"strings"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encoder"
	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
	"github.com/aileron-gateway/aileron-gateway/kernel/er"
	"github.com/aileron-gateway/aileron-gateway/kernel/hash"
	"github.com/aileron-gateway/aileron-gateway/kernel/mac"
)

// ReplaceFunc is the function that
// replaces the given content and  returns new object.
//
// Example to define a string replace function:
//
//	// func(string) string
//	type StringReplaceFunc ReplaceFunc[string]
type ReplaceFunc[T any] func(T) T

// Replacer is the interface that replaces a given object.
//
// Example to define a string replacer interface:
//
//	// interface{ Replace(string) string }
//	type StringReplacer Replacer[string]
type Replacer[T any] interface {
	Replace(T) T
}

func NewStringReplacers(specs ...*k.ReplacerSpec) ([]Replacer[string], error) {
	result := make([]Replacer[string], 0, len(specs))
	for _, s := range specs {
		r, err := NewStringReplacer(s)
		if err != nil {
			return nil, err // Return err as-is.
		}
		result = append(result, r)
	}
	return result, nil
}

func NewStringReplacer(spec *k.ReplacerSpec) (Replacer[string], error) {
	if spec == nil {
		return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeReplacer, Description: ErrDscNil, Detail: "NewStringReplacer"})
	}

	var replacer Replacer[string]
	var err error
	switch v := spec.Replacers.(type) {
	case *k.ReplacerSpec_Fixed:
		replacer = &fixedStringReplacer{
			value: v.Fixed.Value,
		}
	case *k.ReplacerSpec_Value:
		replacer = &valueStringReplacer{
			fromTo: v.Value.FromTo,
		}
	case *k.ReplacerSpec_Left:
		s := v.Left
		replacer = &leftStringReplacer{
			char:   s.Char,
			length: int(s.Length),
		}
	case *k.ReplacerSpec_Right:
		s := v.Right
		replacer = &rightStringReplacer{
			char:   s.Char,
			length: int(s.Length),
		}
	case *k.ReplacerSpec_Trim:
		replacer = &trimStringReplacer{
			cutSets: v.Trim.CutSets,
		}
	case *k.ReplacerSpec_TrimLeft:
		replacer = &trimLeftStringReplacer{
			cutSets: v.TrimLeft.CutSets,
		}
	case *k.ReplacerSpec_TrimRight:
		replacer = &trimRightStringReplacer{
			cutSets: v.TrimRight.CutSets,
		}
	case *k.ReplacerSpec_TrimPrefix:
		replacer = &trimPrefixStringReplacer{
			prefixes: v.TrimPrefix.Prefixes,
		}
	case *k.ReplacerSpec_TrimSuffix:
		replacer = &trimSuffixStringReplacer{
			suffixes: v.TrimSuffix.Suffixes,
		}
	case *k.ReplacerSpec_Encode:
		s := v.Encode
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		f, err2 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2)
		replacer = &encodeStringReplacer{
			pattern:    p,
			encodeFunc: f,
		}
	case *k.ReplacerSpec_Hash:
		s := v.Hash
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		hf, err2 := hashFunc(s.Alg)
		ef, err3 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2, err3)
		replacer = &hashStringReplacer{
			pattern:    p,
			hashFunc:   hf,
			encodeFunc: ef,
		}
	case *k.ReplacerSpec_Regexp:
		s := v.Regexp
		p, err1 := regexpPattern(s.Pattern, s.POSIX, false)
		err = err1
		replacer = &regexStringReplacer{
			pattern: p,
			repl:    []byte(s.Replace),
			literal: s.Literal,
		}
	case *k.ReplacerSpec_Expand:
		s := v.Expand
		p, err1 := regexpPattern(s.Pattern, s.POSIX, false)
		err = err1
		replacer = &expandStringReplacer{
			pattern:  p,
			template: []byte(s.Template),
		}
	case *k.ReplacerSpec_Encrypt:
		s := v.Encrypt
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		encf, password, err2 := encryptFunc(s.Alg, s.Password)
		ef, err3 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2, err3)
		replacer = &encryptStringReplacer{
			pattern:     p,
			encryptFunc: encf,
			encodeFunc:  ef,
			password:    password,
		}
	case *k.ReplacerSpec_HMAC:
		s := v.HMAC
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		hmacf, key, err2 := hmacFunc(s.Alg, s.Key)
		ef, err3 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2, err3)
		replacer = &hmacStringReplacer{
			pattern:    p,
			hmacFunc:   hmacf,
			encodeFunc: ef,
			key:        key,
		}
	default:
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscUnsupported,
		}
	}

	return replacer, err
}

func NewBytesReplacers(specs ...*k.ReplacerSpec) ([]Replacer[[]byte], error) {
	result := make([]Replacer[[]byte], 0, len(specs))
	for _, s := range specs {
		r, err := NewBytesReplacer(s)
		if err != nil {
			return nil, err // Return err as-is.
		}
		result = append(result, r)
	}
	return result, nil
}

func NewBytesReplacer(spec *k.ReplacerSpec) (Replacer[[]byte], error) {
	if spec == nil {
		return nil, (&er.Error{Package: ErrPkg, Type: ErrTypeReplacer, Description: ErrDscNil, Detail: "NewBytesReplacer"})
	}

	var replacer Replacer[[]byte]
	var err error
	switch v := spec.Replacers.(type) {
	case *k.ReplacerSpec_Fixed:
		replacer = &fixedBytesReplacer{
			value: []byte(v.Fixed.Value),
		}
	case *k.ReplacerSpec_Value:
		var from, to [][]byte
		for k, v := range v.Value.FromTo {
			from = append(from, []byte(k))
			to = append(to, []byte(v))
		}
		replacer = &valueBytesReplacer{
			from: from,
			to:   to,
		}
	case *k.ReplacerSpec_Left:
		s := v.Left
		replacer = &leftBytesReplacer{
			char:   []byte(s.Char),
			length: int(s.Length),
		}
	case *k.ReplacerSpec_Right:
		s := v.Right
		replacer = &rightBytesReplacer{
			char:   []byte(s.Char),
			length: int(s.Length),
		}
	case *k.ReplacerSpec_Trim:
		replacer = &trimBytesReplacer{
			cutSets: v.Trim.CutSets,
		}
	case *k.ReplacerSpec_TrimLeft:
		replacer = &trimLeftBytesReplacer{
			cutSets: v.TrimLeft.CutSets,
		}
	case *k.ReplacerSpec_TrimRight:
		replacer = &trimRightBytesReplacer{
			cutSets: v.TrimRight.CutSets,
		}
	case *k.ReplacerSpec_TrimPrefix:
		var prefixes [][]byte
		for _, v := range v.TrimPrefix.Prefixes {
			prefixes = append(prefixes, []byte(v))
		}
		replacer = &trimPrefixBytesReplacer{
			prefixes: prefixes,
		}
	case *k.ReplacerSpec_TrimSuffix:
		var suffixes [][]byte
		for _, v := range v.TrimSuffix.Suffixes {
			suffixes = append(suffixes, []byte(v))
		}
		return &trimSuffixBytesReplacer{
			suffixes: suffixes,
		}, nil
	case *k.ReplacerSpec_Encode:
		s := v.Encode
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		f, err2 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2)
		replacer = &encodeBytesReplacer{
			pattern:    p,
			encodeFunc: f,
		}
	case *k.ReplacerSpec_Hash:
		s := v.Hash
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		hf, err2 := hashFunc(s.Alg)
		ef, err3 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2, err3)
		replacer = &hashBytesReplacer{
			pattern:    p,
			hashFunc:   hf,
			encodeFunc: ef,
		}
	case *k.ReplacerSpec_Regexp:
		s := v.Regexp
		p, err1 := regexpPattern(s.Pattern, s.POSIX, false)
		err = err1
		replacer = &regexBytesReplacer{
			pattern: p,
			repl:    []byte(s.Replace),
			literal: s.Literal,
		}
	case *k.ReplacerSpec_Expand:
		s := v.Expand
		p, err1 := regexpPattern(s.Pattern, s.POSIX, false)
		err = err1
		replacer = &expandBytesReplacer{
			pattern:  p,
			template: []byte(s.Template),
		}
	case *k.ReplacerSpec_Encrypt:
		s := v.Encrypt
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		encf, password, err2 := encryptFunc(s.Alg, s.Password)
		ef, err3 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2, err3)
		replacer = &encryptBytesReplacer{
			pattern:     p,
			encryptFunc: encf,
			encodeFunc:  ef,
			password:    password,
		}
	case *k.ReplacerSpec_HMAC:
		s := v.HMAC
		p, err1 := regexpPattern(s.Pattern, s.POSIX, true)
		hmacf, key, err2 := hmacFunc(s.Alg, s.Key)
		ef, err3 := encodeFunc(s.Encoding)
		err = cmp.Or(err1, err2, err3)
		replacer = &hmacBytesReplacer{
			pattern:    p,
			hmacFunc:   hmacf,
			encodeFunc: ef,
			key:        key,
		}
	default:
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscUnsupported,
		}
	}

	return replacer, err
}

func encodeFunc(enc k.EncodingType) (encoder.EncodeToStringFunc, error) {
	_, ok := encoder.EncodeTypes[enc]
	if !ok {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscUnsupported,
			Detail:      "encoding type " + enc.String(),
		}
	}
	f, _ := encoder.EncoderDecoder(enc)
	return f, nil
}

func hashFunc(alg k.HashAlg) (hash.HashFunc, error) {
	f := hash.FromHashAlg(alg)
	if f == nil {
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscUnsupported,
			Detail:      "hash algorithm" + alg.String(),
		}
	}
	return f, nil
}

func hmacFunc(alg k.HashAlg, key string) (mac.HMACFunc, []byte, error) {
	f := mac.FromHashAlg(alg)
	if f == nil {
		return nil, nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscUnsupported,
			Detail:      "hmac with hash algorithm" + alg.String(),
		}
	}
	rawKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscPattern,
			Detail:      "hmac key must be Hex encoded",
		}).Wrap(err)
	}
	return f, rawKey, nil
}

func encryptFunc(alg k.CommonKeyCryptType, pwd string) (encrypt.EncryptFunc, []byte, error) {
	f := encrypt.EncrypterFromType(alg)
	if f == nil {
		return nil, nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscUnsupported,
			Detail:      "encryption algorithm" + alg.String(),
		}
	}
	password, err := hex.DecodeString(pwd)
	if err != nil {
		return nil, nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscPattern,
			Detail:      "encryption password must be Hex encoded",
		}).Wrap(err)
	}
	return f, password, nil
}

// fixedStringReplacer is the Fixed type replacer.
// This implements Replacer[string] interface.
type fixedStringReplacer struct {
	value string
}

func (r *fixedStringReplacer) Replace(in string) string {
	return r.value
}

// fixedBytesReplacer is the Fixed type replacer.
// This implements Replacer[[]byte] interface.
type fixedBytesReplacer struct {
	value []byte
}

func (r *fixedBytesReplacer) Replace(in []byte) []byte {
	return bytes.Clone(r.value)
}

// valueStringReplacer is the Value type replacer.
// This implements Replacer[string] interface.
type valueStringReplacer struct {
	fromTo map[string]string
}

func (r *valueStringReplacer) Replace(in string) string {
	for k, v := range r.fromTo {
		in = strings.ReplaceAll(in, k, v)
	}
	return in
}

// valueBytesReplacer is the Value type replacer.
// This implements Replacer[[]byte] interface.
type valueBytesReplacer struct {
	from [][]byte
	to   [][]byte
}

func (r *valueBytesReplacer) Replace(in []byte) []byte {
	for i := range r.from {
		in = bytes.ReplaceAll(in, r.from[i], r.to[i])
	}
	return in
}

// leftStringReplacer is the Left type replacer.
// This implements Replacer[string] interface.
type leftStringReplacer struct {
	char   string
	length int
}

func (r *leftStringReplacer) Replace(in string) string {
	rn := []rune(in)
	n := len(rn)
	if n <= r.length {
		return strings.Repeat(r.char, n)
	}
	return strings.Repeat(r.char, r.length) + string(rn[r.length:])
}

// leftBytesReplacer is the Left type replacer.
// This implements Replacer[[]byte] interface.
type leftBytesReplacer struct {
	char   []byte
	length int
}

func (r *leftBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	if n <= r.length {
		return bytes.Repeat(r.char, n)
	}
	return append(bytes.Repeat(r.char, r.length), in[r.length:]...)
}

// rightStringReplacer is the Right type replacer.
// This implements Replacer[string] interface.
type rightStringReplacer struct {
	char   string
	length int
}

func (r *rightStringReplacer) Replace(in string) string {
	rn := []rune(in)
	n := len(rn)
	if n <= r.length {
		return strings.Repeat(r.char, n)
	}
	return string(rn[:n-r.length]) + strings.Repeat(r.char, r.length)
}

// rightBytesReplacer is the Right type replacer.
// This implements Replacer[[]byte] interface.
type rightBytesReplacer struct {
	char   []byte
	length int
}

func (r *rightBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	if n <= r.length {
		return bytes.Repeat(r.char, n)
	}
	return append(in[:n-r.length], bytes.Repeat(r.char, r.length)...)
}

// trimStringReplacer is the Trim type replacer.
// This implements Replacer[string] interface.
type trimStringReplacer struct {
	cutSets []string
}

func (r *trimStringReplacer) Replace(in string) string {
	n := len(in)
	for _, s := range r.cutSets {
		in = strings.Trim(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimBytesReplacer is the Trim type replacer.
// This implements Replacer[[]byte] interface.
type trimBytesReplacer struct {
	cutSets []string
}

func (r *trimBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	for _, s := range r.cutSets {
		in = bytes.Trim(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimLeftStringReplacer is the TrimLeft type replacer.
// This implements Replacer[string] interface.
type trimLeftStringReplacer struct {
	cutSets []string
}

func (r *trimLeftStringReplacer) Replace(in string) string {
	n := len(in)
	for _, s := range r.cutSets {
		in = strings.TrimLeft(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimLeftBytesReplacer is the TrimLeft type replacer.
// This implements Replacer[[]byte] interface.
type trimLeftBytesReplacer struct {
	cutSets []string
}

func (r *trimLeftBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	for _, s := range r.cutSets {
		in = bytes.TrimLeft(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimRightStringReplacer is the TrimRight type replacer.
// This implements Replacer[string] interface.
type trimRightStringReplacer struct {
	cutSets []string
}

func (r *trimRightStringReplacer) Replace(in string) string {
	n := len(in)
	for _, s := range r.cutSets {
		in = strings.TrimRight(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimRightBytesReplacer is the TrimRight type replacer.
// This implements Replacer[[]byte] interface.
type trimRightBytesReplacer struct {
	cutSets []string
}

func (r *trimRightBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	for _, s := range r.cutSets {
		in = bytes.TrimRight(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimPrefixStringReplacer is the TrimPrefix type replacer.
// This implements Replacer[string] interface.
type trimPrefixStringReplacer struct {
	prefixes []string
}

func (r *trimPrefixStringReplacer) Replace(in string) string {
	n := len(in)
	for _, s := range r.prefixes {
		in = strings.TrimPrefix(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimPrefixBytesReplacer is the TrimPrefix type replacer.
// This implements Replacer[[]byte] interface.
type trimPrefixBytesReplacer struct {
	prefixes [][]byte
}

func (r *trimPrefixBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	for _, s := range r.prefixes {
		in = bytes.TrimPrefix(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimSuffixStringReplacer is the TrimSuffix type replacer.
// This implements Replacer[string] interface.
type trimSuffixStringReplacer struct {
	suffixes []string
}

func (r *trimSuffixStringReplacer) Replace(in string) string {
	n := len(in)
	for _, s := range r.suffixes {
		in = strings.TrimSuffix(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// trimSuffixBytesReplacer is the TrimSuffix type replacer.
// This implements Replacer[[]byte] interface.
type trimSuffixBytesReplacer struct {
	suffixes [][]byte
}

func (r *trimSuffixBytesReplacer) Replace(in []byte) []byte {
	n := len(in)
	for _, s := range r.suffixes {
		in = bytes.TrimSuffix(in, s)
		if n != len(in) {
			break
		}
	}
	return in
}

// encodeStringReplacer is the Encode type replacer.
// This implements Replacer[string] interface.
type encodeStringReplacer struct {
	pattern    *regexp.Regexp
	encodeFunc encoder.EncodeToStringFunc
}

func (r *encodeStringReplacer) Replace(in string) string {
	if r.pattern == nil {
		return r.encodeFunc([]byte(in))
	}
	return r.pattern.ReplaceAllStringFunc(in, func(b string) string {
		return r.encodeFunc([]byte(b))
	})
}

// encodeBytesReplacer is the Encode type replacer.
// This implements Replacer[[]byte] interface.
type encodeBytesReplacer struct {
	pattern    *regexp.Regexp
	encodeFunc encoder.EncodeToStringFunc
}

func (r *encodeBytesReplacer) Replace(in []byte) []byte {
	if r.pattern == nil {
		return []byte(r.encodeFunc(in))
	}
	return r.pattern.ReplaceAllFunc(in, func(b []byte) []byte {
		return []byte(r.encodeFunc(b))
	})
}

// hashStringReplacer is the Hash type replacer.
// This implements Replacer[string] interface.
type hashStringReplacer struct {
	pattern    *regexp.Regexp
	hashFunc   hash.HashFunc
	encodeFunc encoder.EncodeToStringFunc
}

func (r *hashStringReplacer) Replace(in string) string {
	if r.pattern == nil {
		h := r.hashFunc([]byte(in))
		return r.encodeFunc(h)
	}
	return r.pattern.ReplaceAllStringFunc(in, func(b string) string {
		return r.encodeFunc(r.hashFunc([]byte(b)))
	})
}

// hashBytesReplacer is the Hash type replacer.
// This implements Replacer[[]byte] interface.
type hashBytesReplacer struct {
	pattern    *regexp.Regexp
	hashFunc   hash.HashFunc
	encodeFunc encoder.EncodeToStringFunc
}

func (r *hashBytesReplacer) Replace(in []byte) []byte {
	if r.pattern == nil {
		h := r.hashFunc(in)
		return []byte(r.encodeFunc(h))
	}
	return r.pattern.ReplaceAllFunc(in, func(b []byte) []byte {
		return []byte(r.encodeFunc(r.hashFunc(b)))
	})
}

// regexStringReplacer is the Regexp type replacer.
// This implements Replacer[string] interface.
type regexStringReplacer struct {
	pattern *regexp.Regexp
	repl    []byte
	literal bool
}

func (r *regexStringReplacer) Replace(in string) string {
	if r.literal {
		return string(r.pattern.ReplaceAllLiteral([]byte(in), r.repl))
	}
	return string(r.pattern.ReplaceAll([]byte(in), r.repl))
}

// regexBytesReplacer is the Regexp type replacer.
// This implements Replacer[[]byte] interface.
type regexBytesReplacer struct {
	pattern *regexp.Regexp
	repl    []byte
	literal bool
}

func (r *regexBytesReplacer) Replace(in []byte) []byte {
	if r.literal {
		return r.pattern.ReplaceAllLiteral(in, r.repl)
	}
	return r.pattern.ReplaceAll(in, r.repl)
}

// expandStringReplacer is the Expand type replacer.
// This implements Replacer[string] interface.
type expandStringReplacer struct {
	pattern  *regexp.Regexp
	template []byte
}

func (r *expandStringReplacer) Replace(in string) string {
	content := []byte(in)
	result := []byte{}
	for _, submatches := range r.pattern.FindAllSubmatchIndex(content, -1) {
		result = r.pattern.Expand(result, r.template, content, submatches)
	}
	return string(result)
}

// expandBytesReplacer is the Expand type replacer.
// This implements Replacer[[]byte] interface.
type expandBytesReplacer struct {
	pattern  *regexp.Regexp
	template []byte
}

func (r *expandBytesReplacer) Replace(in []byte) []byte {
	result := []byte{}
	for _, submatches := range r.pattern.FindAllSubmatchIndex(in, -1) {
		result = r.pattern.Expand(result, r.template, in, submatches)
	}
	return result
}

func regexpPattern(pattern string, posix bool, allowEmpty bool) (*regexp.Regexp, error) {
	if pattern == "" {
		if allowEmpty {
			return nil, nil
		}
		return nil, &er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeReplacer,
			Description: ErrDscPattern,
			Detail:      "empty regular expression",
		}
	}
	if posix {
		exp, err := regexp.CompilePOSIX(pattern)
		if err != nil {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeReplacer,
				Description: ErrDscPattern,
				Detail:      "regular expression POSIX `" + pattern + "`",
			}
		}
		return exp, nil
	} else {
		exp, err := regexp.Compile(pattern)
		if err != nil {
			return nil, &er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeReplacer,
				Description: ErrDscPattern,
				Detail:      "regular expression `" + pattern + "`",
			}
		}
		return exp, nil
	}
}

// encryptStringReplacer is the encryption type replacer.
// This implements Replacer[string] interface.
type encryptStringReplacer struct {
	pattern     *regexp.Regexp
	encryptFunc encrypt.EncryptFunc
	encodeFunc  encoder.EncodeToStringFunc
	password    []byte
}

func (r *encryptStringReplacer) Replace(in string) (s string) {
	if r.pattern == nil {
		ciphertext, err := r.encryptFunc(r.password, []byte(in))
		if err != nil {
			return "!ERROR[" + err.Error() + "]"
		}
		return r.encodeFunc(ciphertext)
	}
	return r.pattern.ReplaceAllStringFunc(in, func(b string) string {
		ciphertext, err := r.encryptFunc(r.password, []byte(b))
		if err != nil {
			return "!ERROR[" + err.Error() + "]"
		}
		return r.encodeFunc(ciphertext)
	})
}

// encryptBytesReplacer is the encryption type replacer.
// This implements Replacer[[]byte] interface.
type encryptBytesReplacer struct {
	pattern     *regexp.Regexp
	encryptFunc encrypt.EncryptFunc
	encodeFunc  encoder.EncodeToStringFunc
	password    []byte
}

func (r *encryptBytesReplacer) Replace(in []byte) (b []byte) {
	if r.pattern == nil {
		ciphertext, err := r.encryptFunc(r.password, in)
		if err != nil {
			return []byte("!ERROR[" + err.Error() + "]")
		}
		return []byte(r.encodeFunc(ciphertext))
	}
	return r.pattern.ReplaceAllFunc(in, func(b []byte) []byte {
		// Here we need to copy the input bytes b
		// because the block cipher appends padding bytes to
		// the input data in the encryptFunc.
		// Padding bytes are unexpectedly inserted into the input data `in`.
		bb := make([]byte, len(b))
		copy(bb, b)
		ciphertext, err := r.encryptFunc(r.password, bb)
		if err != nil {
			return []byte("!ERROR[" + err.Error() + "]")
		}
		return []byte(r.encodeFunc(ciphertext))
	})
}

// hmacStringReplacer is the HMAC type replacer.
// This implements Replacer[string] interface.
type hmacStringReplacer struct {
	pattern    *regexp.Regexp
	hmacFunc   mac.HMACFunc
	encodeFunc encoder.EncodeToStringFunc
	key        []byte
}

func (r *hmacStringReplacer) Replace(in string) string {
	if r.pattern == nil {
		h := r.hmacFunc([]byte(in), r.key)
		return r.encodeFunc(h)
	}
	return r.pattern.ReplaceAllStringFunc(in, func(b string) string {
		return r.encodeFunc(r.hmacFunc([]byte(b), r.key))
	})
}

// hmacBytesReplacer is the HMAC type replacer.
// This implements Replacer[[]byte] interface.
type hmacBytesReplacer struct {
	pattern    *regexp.Regexp
	hmacFunc   mac.HMACFunc
	encodeFunc encoder.EncodeToStringFunc
	key        []byte
}

func (r *hmacBytesReplacer) Replace(in []byte) []byte {
	if r.pattern == nil {
		h := r.hmacFunc(in, r.key)
		return []byte(r.encodeFunc(h))
	}
	return r.pattern.ReplaceAllFunc(in, func(b []byte) []byte {
		return []byte(r.encodeFunc(r.hmacFunc(b, r.key)))
	})
}
