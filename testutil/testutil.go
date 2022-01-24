package testutil

import (
	"fmt"
	"testing"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"
const succeed = "\u2713"
const failed = "\u2717"

func Given(t *testing.T, txt string) {
	t.Logf("%s %s%s", colorCyan, txt, colorReset)
}

func GivenF(t *testing.T, txt string, a ...interface{}) {
	t.Logf(colorCyan + fmt.Sprintf(txt, a) + colorReset)
}

func WhenF(t *testing.T, txt string, a ...interface{}) {
	t.Logf(fmt.Sprintf(colorCyan+"  "+txt+colorReset, a...))
}

func Successf(t *testing.T, txt string, a ...interface{}) {
	t.Logf(fmt.Sprintf(colorGreen+"  "+succeed+"  "+txt+colorReset, a...))
}

func Errorf(t *testing.T, txt string, a ...interface{}) {
	t.Errorf(fmt.Sprintf(colorRed+"  "+failed+"  "+txt+colorReset, a...))
}
