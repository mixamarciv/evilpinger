package main

import (
	"fmt"
	//"io/ioutil"
	//"io"
	"log"
	"strconv"

	//"crypto/md5"
	//"regexp"
	"strings"
	//"time"

	//"github.com/satori/go.uuid"

	"os/exec"
	//"github.com/parnurzeal/gorequest"

	//"text/template"

	mf "gofncstd3000"

	//"github.com/palantir/stacktrace"
	//"runtime/debug"

	//"encoding/json"
	"bufio"
	//"bytes"
	"time"

	"os"

	"regexp"

	"github.com/gosuri/uilive"
)

func init() {

}

func main() {
	log.Println("start app " + mf.CurTimeStrShort())

	/*****
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	for i := 0; i <= 100; i++ {
		fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\ntest 1\ntest 2\ntest 3\n", i, 100)
		time.Sleep(time.Millisecond * 1)
	}

	fmt.Fprintln(writer, "Finished: Downloaded 100GB")
	writer.Stop()
	****/

	line := make(chan string)
	go start_app(line)
	time.Sleep(time.Millisecond * 1)
	/********
		writer := uilive.New()
		writer.Start()
		for {
			msg := <-line
			fmt.Fprintf(writer, msg+"\n")
			time.Sleep(time.Millisecond * 1)
		}
		writer.Stop()
	********/
	var input string
	fmt.Scanln(&input)
}

func start_app(ml chan string) {
	s, _ := readIniFile(os.Args[0] + ".ini")

	line := make(chan string, len(s))
	for i := 0; i < len(s); i++ {
		fmt.Printf("%+v\n", s[i])
		go start_exec(line, s[i])
	}

	h := new(Hosts_info)
	h.Init(len(s))

	writer := uilive.New()
	writer.Start()

	for msg := range line {
		ok := h.Update(msg)
		if ok == 0 {
			fmt.Fprintf(writer, "\n--> ERROR msg: \n%+v\n\n", msg)
		}
		s := h.GetMsg()
		fmt.Fprintf(writer, s)
	}
	writer.Stop()

	var input string
	fmt.Scanln(&input)
}

/********************
server1: ping 192.168.1.1 /t
server2: ping 192.168.1.2 /t
vlad: ping 192.168.1.105 /t
vpn1: ping 10.8.0.1 /t
**********************/

//-------------------------------------------------------------------
type Host_info struct {
	servername string
	host       string
	ip         string
	ttl        string
	ms         []string
}

type Hosts_info struct {
	host            map[string]*Host_info
	servername_list []string

	init_cnt       int
	fmt_str        string
	servername_len int
	host_len       int
	ip_len         int
	ttl_len        int
	ms_len         int
}

func (h *Hosts_info) Init(size int) {
	h.host = make(map[string]*Host_info, size)

	h.servername_len = 4
	h.host_len = 4
	h.ip_len = 4
	h.ttl_len = 4
}

func (h *Hosts_info) Init2() {
	if h.init_cnt > 1 {
		return
	}

	h.init_cnt++
	for _, servername := range h.servername_list {
		ih := h.host[servername]

		if h.servername_len <= len(ih.servername) {
			h.servername_len = len(ih.servername) + 1
		}
		if h.host_len <= len(ih.host) {
			h.host_len = len(ih.host) + 1
		}
		if h.ip_len <= len(ih.ip) {
			h.ip_len = len(ih.ip) + 1
		}
		if h.ttl_len <= len(ih.ttl) {
			h.ttl_len = len(ih.ttl) + 1
		}
	}
	h.fmt_str = "%" + strconv.Itoa(h.servername_len) + "s" +
		"%" + strconv.Itoa(h.host_len) + "s" +
		"%" + strconv.Itoa(h.ttl_len) + "s" +
		"%" + strconv.Itoa(h.ip_len) + "s" +
		"%4s" //ms
}

func (h *Hosts_info) AddServerName(s string) (ok int) {
	//h.servername_list = append(h.servername_list, servername)
	return 1
}

func (h *Hosts_info) Update(s string) (ok int) {

	a := strings.Split(s, " ")
	if len(a) < 5 {
		fmt.Printf("\n--> ERROR: len(a)<5 \n%+v\n\n", a)
		return 0
	}

	servername := a[0]

	ms := a[4]

	//h.host[]
	if i, ok := h.host[servername]; ok {
		i.servername = servername
		i.host = a[1]
		i.ip = a[2]
		i.ttl = a[3]
		i.ms = append([]string{ms}, i.ms...)
		i.ms = i.ms[0:50]
		//ms++
	} else {
		h.servername_list = append(h.servername_list, servername)
		h.init_cnt = 0

		i := new(Host_info)
		i.servername = servername
		i.host = a[1]
		i.ip = a[2]
		i.ttl = a[3]
		i.ms = make([]string, 50)
		i.ms[0] = ms
		i.ms[1] = "-"
		i.ms[2] = "-"
		i.ms[3] = "-"
		i.ms[4] = "-"
		h.host[servername] = i
	}

	return 1
}

func (h *Hosts_info) GetMsg() string {
	h.Init2()

	lfmt := "%70s\n"
	msfmt := "%4s %4s %4s %4s"

	l := fmt.Sprintf(h.fmt_str, "name", "host", "ttl", "ip", "ms")
	l += fmt.Sprintf(msfmt, "", "", "", "")

	s := fmt.Sprintf(lfmt, l)

	for _, servername := range h.servername_list {
		ih := h.host[servername]

		l = fmt.Sprintf(h.fmt_str, ih.servername, ih.host, ih.ttl, ih.ip, ih.ms[0])
		l += fmt.Sprintf(msfmt, ih.ms[1], ih.ms[2], ih.ms[3], ih.ms[4])
		s += fmt.Sprintf(lfmt, l)
	}

	return s
}

//----------------------------------------------------------------

func start_app0() {
	log.Println(mf.CurTimeStrShort() + " start app")

	out, err := exec.Command("date").Output()
	if err != nil {
		printerr("какаято абшибка", err)
	}
	fmt.Printf("The date is %s\n", out)

	line := make(chan string)

	go start_exec(line, "ping")
	go start_exec(line, "ping")
	go start_exec(line, "ping")

	var input string
	fmt.Scanln(&input)
}

func readIniFile(fileread string) ([]string, error) {
	var ret []string
	file, err := os.Open(fileread)
	if err != nil {
		printerr("Error: can't open file: "+fileread, err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		printerr("Error: can't read file: "+fileread, err)
		return nil, err
	}
	return ret, nil
}

func start_cmd(command string, args []string) *exec.Cmd {
	switch len(args) {
	case 0:
		return exec.Command(command)
	case 1:
		return exec.Command(command, args[0])
	case 2:
		return exec.Command(command, args[0], args[1])
	case 3:
		return exec.Command(command, args[0], args[1], args[2])
	case 4:
		return exec.Command(command, args[0], args[1], args[2], args[3])
	case 5:
		return exec.Command(command, args[0], args[1], args[2], args[3], args[4])
	case 6:
		return exec.Command(command, args[0], args[1], args[2], args[3], args[4], args[5])
	}
	return exec.Command(command)
}

func start_exec(line chan string, command string) {
	//cmd := exec.Command("ping", "-t", `ya.ru`)

	var cmd *exec.Cmd
	cmd = exec.Command(command)

	arr1 := strings.SplitN(command, ":", 2)
	servername := strings.Trim(arr1[0], "\r\n\t ")

	arr2 := strings.Split(strings.Trim(arr1[1], "\r\n\t "), " ")
	appname := arr2[0]
	appargs := arr2[1:]

	cmd = start_cmd(appname, appargs)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		printerr("ERROR901: ошибка1 запуска внешнего процесса", err)
		return
	}
	if err := cmd.Start(); err != nil {
		printerr("ERROR902: ошибка2 запуска внешнего процесса", err)
		return
	}

	r := bufio.NewReader(stdout)
	go func() {
		i := 0
		out := ""
		for {
			str, err := r.ReadString('\n')
			if err != nil {
				str = "ERROR readStdout Error!"
			}
			i++
			if out == "" {
				host, ip, ok := parseStr_getHostIp(str)
				if ok == 1 {
					out = servername + " " + host + " " + ip + " "
					continue
				}
				continue
			}
			str = out + parseStr(str)
			//fmt.Println(str)
			line <- str
		}
	}()

	if err := cmd.Wait(); err != nil {
		printerr("ERROR903: ошибка3 ожидания внешнего процесса", err)
		return
	}
}

func regexp_match(re, s string) bool {
	r, err := regexp.Compile(re)
	if err != nil {
		printerr("regexp_match error", err)
		panic(1)
	}
	return r.MatchString(s)
}
func regexp_replace(text string, regx_from string, to string) string {
	reg, err := regexp.Compile(regx_from)
	if err != nil {
		printerr("regexp_replace error", err)
		panic(1)
	}
	text = reg.ReplaceAllString(text, to)
	return text
}

func parseStr_getHostIp(s string) (string, string, int) {
	s = mf.Tr(s, "cp866", "UTF-8")
	s = mf.StrTrim(s)
	if s == "" {
		return "err01", "err01", 0
	}

	//Обмен пакетами с 192.168.1.1 по с 32 байтами данных:
	//Обмен пакетами с ya.ru [213.180.204.3] с 32 байтами данных:
	if regexp_match("\\[[\\.0-9]+\\]", s) {
		//Обмен пакетами с ya.ru [213.180.204.3] с 32 байтами данных:
		ip := regexp_replace(s, "^.*\\[([0-9\\.]+)\\].*$", "$1")
		host := regexp_replace(s, "^.* ([a-zA-Z\\.\\-]+) \\["+ip+"\\].*$", "$1")
		return host, ip, 1
	} else {
		//Обмен пакетами с 192.168.1.1 по с 32 байтами данных:
		ip := regexp_replace(s, "^.* ([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}) .*$", "$1")
		return ip, ip, 1
	}

	return "err02", "err02", 0
}

func parseStr(s string) string {
	s = mf.Tr(s, "cp866", "UTF-8")
	s = mf.StrTrim(s)
	s = mf.StrReplaceRegexp2(s, "Ответ от [0-9\\.]*: число байт=32 время([<=][0-9]+)мс TTL=([0-9]+)", "$2 $1")
	s = strings.Replace(s, "=", "", -1)
	return s
}

func printerr(title string, err error) bool {
	if err != nil {
		serr := "\n\n== ERROR: ======================================\n"
		serr += title + "\n"
		serr += mf.ErrStr(err)
		serr += "\n\n== /ERROR ======================================\n"
		log.Println(serr)
	}
	return false
}
