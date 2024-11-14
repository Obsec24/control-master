package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
    	"syscall"
)

import (
	"github.com/dan2ysgl/logrus"
  	"os"
)

//Configuring vars for logging
var StandardFields logrus.Fields
var Logger = logrus.New()
const defaultFailedCode = 10

func InitLogger(name string, ip string, log_file string, testing_label string, version string) {
   Logger.SetFormatter(&logrus.JSONFormatter{})
   Logger.SetReportCaller(true)
   Logger.SetLevel(logrus.DebugLevel)
   file, err := os.OpenFile(log_file, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
   //defer file.Close()
   if err != nil {
        Logger.Error("Unable to open log file", err)
   }
   StandardFields = logrus.Fields{
        "apk" : name,
        "version" : version,
        "device" : ip,
	    "testing_label" : testing_label,
  }
   Logger.SetOutput(file)
}

func Log(msg string) {
	fmt.Println("[*]", msg)
}

func Command(out io.Writer, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = out
	cmd.Run()
}

func RunCommandOld(name string, arg ...string) (exitCode int){
	cmd := exec.Command(name, arg...)
		
	if err := cmd.Run(); err != nil {
		//try to get exit code
		if exitError, ok := err.(*exec.ExitError); ok{
			ws := exitError.Sys().(syscall.WaitStatus)
            		exitCode = ws.ExitStatus()
		}else{
			exitCode = defaultFailedCode
		}
	} else {
		//success, exitCode 0
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return
}

func RunCommand(name string, arg ...string) int {
	cmd := exec.Command(name, arg...)
    exitCode := 0
	if err := cmd.Start(); err != nil {
	    fmt.Println("[*] Error: %v ", err)
	}

	if err := cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            // The program has exited with an exit code != 0
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                exitCode = status.ExitStatus()
                fmt.Println("ExitCode received: %d", exitCode)
            }
        }
    }
    return exitCode
 }

func UploadAPK(req *http.Request) error {
	req.ParseMultipartForm(10 << 20) // 10 MB
	if file, _, err := req.FormFile("apk"); err != nil {
		return errors.New("Control-Upload: Error parsing multipart form")
	}else{
		defer file.Close()
		if content, err := ioutil.ReadAll(file); err != nil {
			return errors.New("Control-Upload: Error reading apk")
		}else if err := ioutil.WriteFile("base.apk", content, 0644); err != nil {
				return errors.New("Control-Upload: Error building apk")
		}else{
			return nil
		}
	}
}
