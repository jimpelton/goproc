package trfunc

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Lines in a transfer function consist of space separated (either space or tab) double-prec
// values. splitReg is the regex that is used to split the lines.
const splitReg = "[\\s\\t]+"
var split *regexp.Regexp

func init() {
	var err error
	if split, err = regexp.Compile(splitReg); err != nil {
		log.Fatal(err.Error())
	}
}

type TFOpacity struct {
	knots []TFOpacityKnot
}

type TFOpacityKnot struct {
	S     float64 // between 0.0 .. 1.0
	Alpha float64 // between 0.0 .. 1.0
}

func OpenTFOpacityFile(fileName string) (*TFOpacity, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	 tf := &TFOpacity{}
	if err := tf.readLines(f); err != nil {
		return nil, err
	}

	return tf, nil
}

func OpenTFOpacityString(function string) (*TFOpacity, error) {
	tf := &TFOpacity{}
	r := strings.NewReader(function)
	if err := tf.readLines(r); err != nil {
		return nil, err
	}
	return tf, nil
}

func (tf *TFOpacity) Interpolate(v float64) float64 {
	if len(tf.knots) == 0 {
		return 0.0
	}

	if v <= tf.knots[0].Alpha {
		return tf.knots[0].Alpha
	}

	maxIdx := len(tf.knots) - 1

	if v >= tf.knots[maxIdx].Alpha {
		return tf.knots[maxIdx].Alpha
	}

	var (
		k0, k1 TFOpacityKnot
	)

	idx := int((v * float64(maxIdx)) + 0.5)
	if idx > maxIdx {
		k0 = tf.knots[maxIdx - 1]
		k1 = tf.knots[maxIdx]
	} else if idx == 0 {
		k0 = tf.knots[0]
		k1 = tf.knots[1]
	} else {
		k0 = tf.knots[idx - 1]
		k1 = tf.knots[idx]
	}

	d := (v - k0.S) / (k1.S - k0.S)
	return k0.Alpha * (1.0 - d) + k1.Alpha* d
}

func (tf *TFOpacity) readLines(r io.Reader) (err error) {
	const maxLines = 32768

	scan := bufio.NewScanner(r)
	for i:=0; scan.Scan() && i < maxLines; i++ {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts  := split.Split(line, 2)
		k := TFOpacityKnot{}

		if k.S, err = strconv.ParseFloat(parts[0], 64); err != nil {
			return err
		}

		if k.Alpha, err = strconv.ParseFloat(parts[1], 64); err != nil {
			return err
		}

		tf.knots = append(tf.knots, k)
	}

	return nil
}

type TFColor struct {
	knots []TFColorKnot
}

type TFColorKnot struct {
	s       float64
	r, g, b, a float64

}
func OpenTFColorFile(file string) (*TFColor, error) {
	return nil, nil
}

func OpenTFColorString(function string) (*TFColor, error) {
	return nil, nil
}
