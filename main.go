/*
 * @Author: fyfishie
 * @Date: 2023-04-22:11
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-22:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	"flag"
	"strings"

	"github.com/fyfishie/rdns/parse"
)

var inputFile = flag.String("if", "", "input file path")
var outputFile = flag.String("of", "output.txt", "output file path")

// var mode = flag.String("m", "mix", "input mode, mix or point (default \"mix\")")
var ips = flag.String("i", "", "input ips, join by \",\"")
var speed = flag.Int("n", 1000, "speed limit, (default \"1000\")")

func init() {
	flag.Parse()
	// if (*mode) != lib.MODE_MIX || (*mode) != lib.MODE_POINT {
	// 	panic("invalid input mode " + *mode + "\n")
	// }
}
func main() {
	ipList := []string{}
	if *ips != "" {
		ipList = strings.Split(*ips, ",")
	}

	p := parse.NewParser("./input.txt", "output.txt", "point", 10000, ipList)
	p.Run()
}
