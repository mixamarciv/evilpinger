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
	"time"

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
	//"time"

	"os"

	"regexp"

	"github.com/gosuri/uilive"
	"github.com/jroimartin/gocui"
	"github.com/qiniu/iconv"
)

var inifile string

func tr(s string, from string, to string) string {
	cd, err := iconv.Open(to, from)
	if err != nil {
		ret := "ERROR gofncstd3000.tr(): iconv.Open(" + to + "," + from + ") failed!"
		return ret
	}
	defer cd.Close()

	ret := cd.ConvString(s)
	return ret
}

func init() {
	inifile = os.Args[0] + ".ini"
}

func main() {
	//log.Println("start app " + mf.CurTimeStrShort())
	WriteErrorLog("start app")
	start_app2()
}

func start_app2() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	s, _ := readIniFile(inifile)

	h := new(Hosts_info)
	h.Init(len(s))

	line := make(chan string, len(s))
	for i := 0; i < len(s); i++ {
		//fmt.Printf("%+v\n", s[i])

		t := strings.SplitN(s[i], ":", 2)
		servername := strings.Trim(t[0], "\t\r\n ")
		if servername != "" {
			h.AddServerName(servername)

			go start_exec(line, s[i])
		}
	}

	go func() {
		for msg := range line {
			ok := h.Update(msg)
			if ok == 0 {
				log.Panicln("\n--> ERROR msg: \n%+v\n\n" + msg)
			}
			s := h.GetMsg()

			g.Execute(func(g *gocui.Gui) error {
				v, err := g.View("ctr")
				if err != nil {
					log.Panicln(err)
					return err
				}
				v.Clear()
				fmt.Fprintf(v, s)
				return nil
			})
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ctr", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Printf("\n--> ERROR002 err: \n%+v\n\n", err)
			log.Panicln(err)
			return err
		}
		fmt.Fprintln(v, "prepare..\nread file: "+inifile)
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func update_console(g *gocui.Gui, s string) {
	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("ctr")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintf(v, s)
		return nil
	})
}

//==============================================================================
func start_app() {
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

	init_cnt int

	fmt_str   string
	msfmt_str string
	fmt_space string

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
		ih, ok := h.host[servername]

		if h.servername_len <= len(servername) {
			h.servername_len = len(servername) + 1
		}

		if ok == false {
			continue
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

	ms_len := 5
	h.fmt_str = "%" + strconv.Itoa(h.servername_len) + "s" +
		"%" + strconv.Itoa(h.host_len) + "s" +
		"%" + strconv.Itoa(h.ttl_len) + "s" +
		"%" + strconv.Itoa(h.ip_len) + "s" +
		"%" + strconv.Itoa(ms_len+1) + "s" //ms
	h.msfmt_str = " %5s %5s %5s %5s %5s %5s %5s %5s"
	h.fmt_space = fmt.Sprintf("%"+strconv.Itoa(h.servername_len+h.host_len+h.ttl_len+h.ip_len+ms_len)+"s", "")
}

func (h *Hosts_info) AddServerName(servername string) (ok int) {
	h.servername_list = append(h.servername_list, servername)
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
		//h.servername_list = append(h.servername_list, servername)
		h.init_cnt = 0

		i := new(Host_info)
		i.servername = servername
		i.host = a[1]
		i.ip = a[2]
		i.ttl = a[3]
		i.ms = make([]string, 50)
		i.ms[0] = ms
		for j := 1; j < 20; j++ {
			i.ms[j] = "-"
		}

		h.host[servername] = i
		h.Init2()
	}

	return 1
}

func (h *Hosts_info) GetMsg() string {
	h.Init2()

	//lfmt := "%70s\n"

	s := h.fmt_space + " *out - time out\n"
	s += h.fmt_space + " *non - not found\n"
	s += h.fmt_space + " *fail - general failure\n"
	s += h.fmt_space + " *err - app error parse data\n"

	s += "\n"

	s += fmt.Sprintf(h.fmt_str, "name", "host", "ttl", "ip", "ms") + "\n"

	for _, servername := range h.servername_list {
		ih, ok := h.host[servername]
		//s += servername + "\n"
		//continue
		if ok == false {
			l := fmt.Sprintf(h.fmt_str, servername, "-", "-", "-", "-")
			l += fmt.Sprintf(h.msfmt_str, "-", "-", "-", "-", "-", "-", "-", "-")
			s += l + "\n"
		} else {
			l := fmt.Sprintf(h.fmt_str, ih.servername, ih.host, ih.ttl, ih.ip, ih.ms[0])
			l += fmt.Sprintf(h.msfmt_str, ih.ms[1], ih.ms[2], ih.ms[3], ih.ms[4], ih.ms[5], ih.ms[6], ih.ms[7], ih.ms[8])
			s += l + "\n"
		}
	}

	return s
}

//----------------------------------------------------------------
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
		printerr("ERROR901: ошибка1 pipeOut внешнего процесса", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		printerr("ERROR901: ошибка1 pipeErr внешнего процесса", err)
		return
	}
	if err := cmd.Start(); err != nil {
		printerr("ERROR902: ошибка2 запуска внешнего процесса", err)
		return
	}

	r := bufio.NewReader(stdout)
	out := ""
	go func() {
		i := 0
		for {
			str, err := r.ReadString('\n')
			if err != nil {
				str = "ERROR readStdout Error!"
				//log.Panicln(str, err)
				//line <- out + "err err"
				return
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

	r2 := bufio.NewReader(stderr)
	go func() {
		i := 0
		for {
			str, err := r2.ReadString('\n')
			if err != nil {
				str = "ERROR readStdErr Error!"
				//log.Panicln(str, err)
				//line <- out + "err err"
				return
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
		//out := servername + " host ipnotfound and non"
		//line <- out
		time.Sleep(time.Second * 3)
		go start_exec(line, command)
		return
	}

}

func regexp_match(re, s string) bool {
	r, err := regexp.Compile(re)
	if err != nil {
		printerr("regexp_match Compile error", err)
		panic(1)
	}
	return r.MatchString(s)
}
func regexp_replace(text string, regx_from string, to string) string {
	reg, err := regexp.Compile(regx_from)
	if err != nil {
		printerr("regexp_replace Compile error", err)
		panic(1)
	}
	text = reg.ReplaceAllString(text, to)
	return text
}

func parseStr_getHostIp(s string) (string, string, int) {
	s = tr(s, "cp866", "UTF-8")
	s = mf.StrTrim(s)
	if s == "" {
		return "err", "err", 0
	}

	if regexp_match("При проверке связи не удалось обнаружить узел", s) {
		host := regexp_replace(s, "При проверке связи не удалось обнаружить узел", "")
		host = regexp_replace(host, "Проверьте имя узла и.*$", "")
		host = regexp_replace(host, "\\.$", "")
		host = mf.StrTrim(host)
		return host, "notfound", 1
	}
	if regexp_match("Проверьте имя узла и повторите попытку", s) {
		return "non", "notfound", 0
	}

	if regexp_match("Обмен пакетами с ", s) {
		//Обмен пакетами с 192.168.1.1 по с 32 байтами данных:
		//Обмен пакетами с ya.ru [213.180.204.3] с 32 байтами данных:
		if regexp_match("\\[[\\.0-9]+\\]", s) {
			//Обмен пакетами с ya.ru [213.180.204.3] с 32 байтами данных:
			ip := regexp_replace(s, "^.*\\[([0-9\\.]+)\\].*$", "$1")
			host := regexp_replace(s, "^.* ([a-zA-Z0-9\\.\\-]+) \\["+ip+"\\].*$", "$1")
			return host, ip, 1
		} else {
			//Обмен пакетами с 192.168.1.1 по с 32 байтами данных:
			ip := regexp_replace(s, "^.* ([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}) .*$", "$1")
			return ip, ip, 1
		}
	}

	WriteErrorLog("parseip error: " + s)
	return "err", "err", 1
}

func parseStr(s string) string {
	s = tr(s, "cp866", "UTF-8")
	s = mf.StrTrim(s)
	if s == "" {
		return "- -"
	}
	if regexp_match("Ответ от [0-9\\.]+: число байт=[\\d]+ время", s) {
		s = mf.StrRegexpReplace(s, "Ответ от [0-9\\.]+: число байт=[\\d]+ время([<=][0-9]+)мс TTL=([0-9]+)", "$2 $1")
		s = strings.Replace(s, "=", "", -1)
	} else if regexp_match("Превышен интервал ожидания для запроса.", s) {
		s = "- out"
	} else if regexp_match("Заданн.+ (узел|сеть) недоступ", s) || regexp_match("Проверьте имя узла и повторите попытку.", s) {
		s = "- non"
	} else if regexp_match("General failure.", s) {
		s = "- fail"
	} else if regexp_match("Статистика Ping для ", s) ||
		regexp_match("Пакетов.*отправлено.*получено.*потер", s) || regexp_match("Приблизительное время приема-передачи", s) ||
		regexp_match("Минимальное.*Максимальное.*Среднее", s) || regexp_match("\\(\\d+\\% потерь\\)", s) {
		s = "- -"
	} else {
		WriteErrorLog("parse error: " + s)
		s = "err err"
	}
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

func WriteErrorLog(s string) {
	mf.FileAppendStr(os.Args[0]+".log", "\n"+mf.CurTimeStrShort()+": "+s)
}
