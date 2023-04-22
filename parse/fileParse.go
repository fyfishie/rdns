/*
 * @Author: fyfishie
 * @Date: 2023-04-22:11
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-22:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package parse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/fyfishie/ipop"
	"github.com/fyfishie/rdns/lib"
)

type Parser struct {
	ipList         []string
	rdpath         string
	wtpath         string
	ipInputChan    chan string
	rdnsResultChan chan lib.RDNSResItem
	doneChan       chan struct{}
	doneNum        int
	runtineNum     int
	writeDone      chan struct{}
}

func NewParser(rdpath string, wtpath string, mode string, runtineNum int, ipList []string) *Parser {
	return &Parser{
		ipList:         ipList,
		rdpath:         rdpath,
		wtpath:         wtpath,
		ipInputChan:    make(chan string, 1000),
		rdnsResultChan: make(chan lib.RDNSResItem, 1000),
		doneChan:       make(chan struct{}, 100),
		writeDone:      make(chan struct{}),
		runtineNum:     runtineNum,
		doneNum:        runtineNum,
	}
}

func (p *Parser) Run() {
	p.writeRun(p.wtpath)
	for i := 0; i < p.runtineNum; i++ {
		p.rdnsRutine()
	}
	if p.rdpath != "" {
		rfi, err := os.OpenFile(p.rdpath, os.O_RDONLY, 0000)
		if err != nil {
			panic(err.Error())
		}
		defer rfi.Close()
		rdr := bufio.NewReader(rfi)
		count := 0
		for {
			count++
			if count%10000 == 0 {
				fmt.Printf("count: %v\n", count)
			}
			line, _, err := rdr.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err.Error())
			}
			lineStr := string(line)
			if !ipop.IsIPReg(lineStr) {
				fmt.Println("invalid ip format, ip = " + lineStr + "err = " + err.Error())
				continue
			}
			p.ipInputChan <- lineStr
		}
	}
	for _, ip := range p.ipList {
		if !ipop.IsIPReg(ip) {
			fmt.Println("invalid ip format, ip = " + ip)
			continue
		}
		p.ipInputChan <- ip
	}
	close(p.ipInputChan)
	for i := 0; i < p.doneNum; i++ {
		<-p.doneChan
	}
	close(p.doneChan)
	fmt.Println("close done")
	close(p.rdnsResultChan)
	fmt.Println("rdns closed")
	<-p.writeDone
}

// one.one.one.one.
// lookup 81.70.76.237: Name or service not known
func (p *Parser) rdnsRutine() {
	go func() {
		for ip := range p.ipInputChan {
			if !ipop.IsIPReg(ip) {
				fmt.Printf("invalid ip input: ip = %v\n", ip)
				continue
			}
			ptr, err := net.LookupAddr(ip)
			if err != nil {
				continue
			}
			domains := []string{}
			for _, rdnsRes := range ptr {
				if strings.HasPrefix(rdnsRes, "lookup") {
					continue
				}
				domains = append(domains, strings.TrimSuffix(rdnsRes, "."))
			}
			p.rdnsResultChan <- lib.RDNSResItem{
				IP:      ip,
				Domains: domains,
			}
		}
		p.doneChan <- struct{}{}
	}()
}

func (p *Parser) writeRun(wtpath string) {
	go func() {
		wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err.Error())
		}
		defer wfi.Close()
		wtr := bufio.NewWriter(wfi)
		for item := range p.rdnsResultChan {
			bs, err := json.Marshal(item)
			if err != nil {
				fmt.Println("error in parse rdns result into bytes, err = " + err.Error())
				continue
			}
			wtr.Write(bs)
			wtr.WriteString("\n")
		}
		fmt.Println("write done")
		wtr.Flush()
		p.writeDone <- struct{}{}
	}()
}
