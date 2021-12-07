/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var xorCmd = &cobra.Command{
	Use:   "xor",
	Short: "XOR encryption",
	Long:  `Encrypt stuff using XOR, is this 1882?`,
	Run: func(cmd *cobra.Command, args []string) {
		formatIn, _ := cmd.Flags().GetString("format-in")
		formatOut, _ := cmd.Flags().GetString("format-out")
		file, _ := cmd.Flags().GetString("file")
		str, _ := cmd.Flags().GetString("str")
		key, _ := cmd.Flags().GetString("key")
		keyType, _ := cmd.Flags().GetString("key-type")
		nowrap, _ := cmd.Flags().GetBool("nowrap")
		fmt.Println(xor(formatIn, formatOut, file, str, key, keyType, nowrap))
	},
}

func init() {
	rootCmd.AddCommand(xorCmd)

	xorCmd.Flags().StringP(
		"format-in", "", "hex", "Output format, choose between csharp, vba, raw, hex",
	)
	xorCmd.Flags().StringP(
		"format-out", "", "csharp", "Output format, choose between csharp, vba, raw, hex",
	)
	xorCmd.Flags().StringP(
		"file", "f", "", "Input file",
	)
	xorCmd.Flags().StringP(
		"str", "s", "", "Input string",
	)
	xorCmd.Flags().StringP(
		"key", "k", "13", "Key to use with the XOR algorithm",
	)
	xorCmd.Flags().StringP(
		"key-type", "", "string", "Key type, choose between decimal, hex, string",
	)
	xorCmd.Flags().BoolP(
		"nowrap", "", false, "Disables wrapping for csharp and vba formats",
	)

}

func encrypt(in []byte, key string, keyType string) []byte {

	var out []byte

	if keyType == "string" {
		k := []byte(key)

		for i, _ := range in {
			out = append(out, in[i]^k[i%len(k)])
		}
	} else if keyType == "hex" {
		k, err := hex.DecodeString(key)

		if err != nil {
			log.Fatal(err.Error())
		}

		for i, _ := range in {
			out = append(out, in[i]^k[i%len(k)])
		}
	} else if keyType == "decimal" {
		decKey, err := strconv.Atoi(key)

		if err != nil {
			log.Fatal(err.Error())
		}

		k := byte(decKey)

		for i, _ := range in {
			out = append(out, in[i]^k)
		}
	}

	return out
}

func xor(formatIn string, formatOut string, file string, str string, key string, keyType string, nowrap bool) string {

	var in string

	if file == "" && str == "" {
		log.Fatal("Either provide a file or a string")
	} else if file != "" && str != "" {
		log.Fatal("Provide either a file or a string, not both")
	} else if file != "" {
		inBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err.Error())
		}
		if formatIn == "raw" {
			in = base64.StdEncoding.EncodeToString(inBytes)
		} else {
			//remove the newline from here otherwise it throws a fit
			in = strings.Replace(string(inBytes), "\n", "", -1)
		}
	} else if str != "" {
		if formatIn == "raw" {
			in = base64.StdEncoding.EncodeToString([]byte(str))
		} else {
			in = str
		}
	} else {
		fmt.Println("Something unexpected happened")
	}

	plaintext := sDecode(in, formatIn)
	ciphertext := encrypt(plaintext, key, keyType)
	fin := sEncode(ciphertext, formatOut, nowrap)

	return fin
}
