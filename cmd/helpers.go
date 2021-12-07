package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

func sDecode(in string, format string) []byte {

	var out []byte

	switch f := format; f {
	case "raw":
		out = rawDecode(in)
	case "hex":
		out = hexDecode(in)
	case "csharp":
		out = csharpDecode(in)
	case "vba":
		out = vbaDecode(in)
	}
	return out
}

func sEncode(in []byte, format string, nowrap bool) string {

	var out string

	switch f := format; f {
	case "raw":
		out = rawEncode(in)
	case "hex":
		out = hexEncode(in)
	case "csharp":
		out = csharpEncode(in, nowrap)
	case "vba":
		out = vbaEncode(in, nowrap)
	}
	return out
}

func rawDecode(in string) []byte {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		fmt.Println("could not base64 decode")
		log.Fatal(err.Error())
	}
	return out
}

func hexDecode(in string) []byte {
	if strings.HasPrefix(in, "\\x") {
		in = strings.Replace(in, "\\x", "", -1)
	} else if strings.HasPrefix(in, "0x") {
		in = strings.Replace(in, "0x", "", -1)
	}

	in = strings.Replace(in, " ", "", -1)
	in = strings.Replace(in, ",", "", -1)
	out, err := hex.DecodeString(in)
	if err != nil {
		log.Fatal(err.Error())
	}
	return out
}

func csharpDecode(in string) []byte {
	log.Fatal("Not supported yet")
	return []byte("Not supported yet")
}

func vbaDecode(in string) []byte {
	log.Fatal("Not supported yet")
	return []byte("Not supported yet")
}

func rawEncode(in []byte) string {
	return string(in)
}

func hexEncode(in []byte) string {
	return hex.EncodeToString(in)
}

func csharpEncode(in []byte, nowrap bool) string {
	var form []string
	var out string

	if !nowrap {
		j := 0
		size := 15
		for i := 0; i < len(in); i += size {
			j += size
			if j > len(in) {
				j = len(in)
			}
			//here I want to take the subslice, convert each byte to hex and join them
			var subform []string
			for _, b := range in[i:j] {
				subform = append(subform, fmt.Sprintf("0x%s", hex.EncodeToString([]byte{b})))
			}
			form = append(form, strings.Join(subform, ","))
		}
		out = fmt.Sprintf("byte[] buf = new byte[%d] {\n%s\n};", len(in), strings.Join(form, ",\n"))
	} else {
		for _, b := range in {
			form = append(form, fmt.Sprintf("0x%s", hex.EncodeToString([]byte{b})))
		}
		out = fmt.Sprintf("byte[] buf = new byte[%d] {%s};", len(in), strings.Join(form, ","))
	}

	return out
}

func vbaEncode(in []byte, nowrap bool) string {
	var form []string
	var out string

	if !nowrap {
		j := 0
		size := 50
		for i := 0; i < len(in); i += size {
			j += size
			if j > len(in) {
				j = len(in)
			}
			//here I want to take the subslice, convert each byte to hex and join them
			var subform []string
			for _, b := range in[i:j] {
				subform = append(subform, fmt.Sprintf("%d", b))
			}
			form = append(form, strings.Join(subform, ", "))
		}
		out = fmt.Sprintf("buf = Array(%s)", strings.Join(form, ", _\n"))
	} else {
		for _, b := range in {
			form = append(form, fmt.Sprintf("%d", b))
			out = fmt.Sprintf("buf = Array(%s)", strings.Join(form, ", "))
		}
	}

	return out
}
