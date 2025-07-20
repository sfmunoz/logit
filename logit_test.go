//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//
//

package logit_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/sfmunoz/logit"
)

// 2025-07-20T16:41:50.744Z
const timePat = `2[0-9]{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]\.[0-9]{3}Z`

// 0d00h00m00.000s
const durPat = `[0-9]d([0-1][0-9]|2[0-3])h[0-5][0-9]m[0-5][0-9]\.[0-9]{3}s`

func assert(t *testing.T, out *bytes.Buffer, re string) {
	want, err := regexp.Compile(re)
	if err != nil {
		t.Fatalf("regexp.Compile(%s) failed: %s", re, err)
	}
	got := out.String()
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("assert(): got='%s' doesn't have '\\n' suffix", got)
	}
	got = strings.TrimRight(got, "\n")
	if !want.MatchString(got) {
		t.Fatalf("assert(): got='%s' doesn't match want='%s'", got, want)
	}
}

func TestLogit1(t *testing.T) {
	var out bytes.Buffer
	l := logit.Logit().
		WithWriter(&out).
		WithColor(false)
	l.Info("hello")
	re := `^` + timePat + ` ` + durPat + ` \[I] hello$`
	assert(t, &out, re)
}

func TestLogit2(t *testing.T) {
	var out bytes.Buffer
	l := logit.Logit().
		WithWriter(&out).
		WithSymbolSet(logit.SymbolUnicodeUp).
		WithColor(false)
	l.Info("hello")
	re := `^` + timePat + ` ` + durPat + ` Ⓘ hello$`
	assert(t, &out, re)
}

func TestLogit3(t *testing.T) {
	var out bytes.Buffer
	l := logit.Logit().
		WithWriter(&out).
		WithColor(false).
		With("k1", "v1")
	l.Info("hello")
	re := `^` + timePat + ` ` + durPat + ` \[I] hello k1=v1$`
	assert(t, &out, re)
}

func TestLogit4(t *testing.T) {
	var out bytes.Buffer
	l := logit.Logit().
		WithWriter(&out).
		WithColor(false).
		With("k1", "v1").
		WithGroup("g1").
		With("k2", "v2")
	l.Info("hello")
	re := `^` + timePat + ` ` + durPat + ` \[I] hello k1=v1 g1.k2=v2$`
	assert(t, &out, re)
}
