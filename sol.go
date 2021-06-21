package main

import (
	"fmt"
	"github.com/CodyGuo/win"
	"github.com/kardianos/service"
	"log"
	"net/http"
	"os"
)

var serviceConfig = &service.Config{
	Name:        "ShutdownOnLAN",
	DisplayName: "shutdown on lan",
	Description: "shutdown pc by lan request",
}

func main() {
	prog := &Program{}
	s, err := service.New(prog, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
		return
	}
	//install, uninstall, start, stop 的另一种实现方式
	err = service.Control(s, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	log.Println("开始服务")
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	log.Println("停止服务")
	return nil
}

func (p *Program) run() {
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		shutdown()
	})
	http.HandleFunc("/reboot", func(w http.ResponseWriter, r *http.Request) {
		reboot()
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Shutdown on lan is running!")
	})
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func reboot() {
	getPrivileges()
	win.ExitWindowsEx(win.EWX_REBOOT, 0)
}

func shutdown() {
	getPrivileges()
	win.ExitWindowsEx(win.EWX_SHUTDOWN, 0)
}

func getPrivileges() {
	var hToken win.HANDLE
	var tkp win.TOKEN_PRIVILEGES

	win.OpenProcessToken(win.GetCurrentProcess(), win.TOKEN_ADJUST_PRIVILEGES|win.TOKEN_QUERY, &hToken)
	win.LookupPrivilegeValueA(nil, win.StringToBytePtr(win.SE_SHUTDOWN_NAME), &tkp.Privileges[0].Luid)
	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Attributes = win.SE_PRIVILEGE_ENABLED
	win.AdjustTokenPrivileges(hToken, false, &tkp, 0, nil, nil)
}
