package traffic

import (
	. "github.com/Obsec24/control-master/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const  LOG_FILE = "/app/logging/log/operation.privapp.log"

var ip string
var name string
var testing_label string
var version string

type resp struct {
	Ok  bool
	Msg string
	Code int
}

func init() {
	ip, name, testing_label, version = "", "", "", ""
	InitLogger(name, ip, LOG_FILE, testing_label, version)
}

func config(res http.ResponseWriter, req *http.Request) {
	ip = req.FormValue("ip")
	testing_label = req.FormValue("testing_label")
	version = req.FormValue("version")

	if ip == "" || testing_label == "" {
                send, _ := json.Marshal(resp{false, "Mandatory target device IP or testing_label are missing", 1})
                fmt.Fprintf(res, string(send))
        } else {
		if tmp := req.FormValue("name"); tmp != "" {
			name = tmp
		}
		InitLogger(name, ip, LOG_FILE, testing_label, version)
		Log("Configured => IP: " + ip + " Name: " + name)
		send, _ := json.Marshal(resp{true, "Configured target device IP", 0})
		fmt.Fprintf(res, string(send))
		Logger.WithFields(StandardFields).Info("Target device IP configured")
	}
}

func cert(res http.ResponseWriter, req *http.Request) {
	Log("Request for CA certificate configuration received")
	 if ip == "" {
                send, _ := json.Marshal(resp{false, "Target device IP is not configured", 1})
                fmt.Fprintf(res, string(send))
        } else {
		//Logger.WithFields(StandardFields).Info("Configuring CA certificate")
		exitcode := RunCommand("scripts/cert.sh", ip)
		if exitcode == 0 {
			send, _ := json.Marshal(resp{true, "CA certificate configured, after the phone reboots the certificate should be installed", 0})
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Info("MITM CA certificate installed")
		} else {
			send, _ := json.Marshal(resp{false, "CA certificate configuration failed!", 1})
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Error("CA certificate configuration failed!")
		}
	}
}

func hooker(res http.ResponseWriter, req *http.Request) {
	Log("Request to install Frida on target device received")
	if ip == "" {
                send, _ := json.Marshal(resp{false, "Target device IP is not configured", 1})
                fmt.Fprintf(res, string(send))
        } else {
		//Logger.WithFields(StandardFields).Info("Installing frida on target device")
        	exitcode := RunCommand("scripts/hooker.sh", ip)
		if exitcode == 0 {
			send, _ := json.Marshal(resp{true, "Frida successfuly installed", 0})
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Info("Frida successfuly installed")
		}else{
			send, _ := json.Marshal(resp{false, "Frida install failed", 0})
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Error("Frida install failed!")
		}
	}
}

func upload(res http.ResponseWriter, req *http.Request) {
	Log("Receiving APK")
	//Logger.WithFields(StandardFields).Info("Uploading APK")
	if err := UploadAPK(req); err == nil {
		name = ApkName()
		Log("Configured => Name: " + name)
		InitLogger(name, ip, LOG_FILE, testing_label, version)
		send, _ := json.Marshal(resp{true, "APK uploaded successfuly", 0})
		fmt.Fprintf(res, string(send))
		Logger.WithFields(StandardFields).Info("APK uploaded successfuly")
	} else {
		send, _ := json.Marshal(resp{false, "Error uploading APK", 1})
		fmt.Fprintf(res, string(send))
		Logger.WithFields(StandardFields).Error("Error uploading APK")
	}
}

func phaseOne(res http.ResponseWriter, req *http.Request) {
	Log("Request to start the proxy and frida's bypass pinning received")
	timeout := req.FormValue("timeout")
	permissions := req.FormValue("permissions")
	reboot := req.FormValue("reboot")
	if tmp := req.FormValue("name"); tmp != "" {
		name = tmp
	}
	if ip == "" || timeout == "" {
		send, _ := json.Marshal(resp{false, "Timeout or target device is not configured", 1})
		fmt.Fprintf(res, string(send))
	} else if _, err := os.Stat("base.apk"); os.IsNotExist(err) && name == "" {
		send, _ := json.Marshal(resp{false, "Upload APK to be analyzed or configure name", 1})
		fmt.Fprintf(res, string(send))
	} else {		
		//Logger.WithFields(StandardFields).Info("Starting the proxy and frida's bypass pinning")
		//Add other params
		args := []string{"-t", timeout, "-d", ip, "-a", name, "-l", testing_label}
		if permissions == "True"{
			args = append(args, "-p")
		}
		if reboot == "True" {
			args = append(args, "-r")
		}
		exitcode := RunCommand("scripts/start.py", args...)
		if exitcode == 0 {
			send, _ := json.Marshal(resp{true, "Traffic capture in idle phase finished successfully", exitcode})
			Log("Successful traffic capture in idle phase")
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Info("Successful traffic capture in idle phase")
		} else {
			send, _ := json.Marshal(resp{false, "Traffic capture in idle phase failed", exitcode})
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Error("Traffic capture in idle phase failed")
		}
		
	}
}

func phaseTwo(res http.ResponseWriter, req *http.Request) {
	Log("Starting to capture traffic in phase II")
	Command(nil, "mv", "output.log", "first.privapp.log")
	args := []string{"-d", ip, "-a", name, "-l", testing_label}
	timeout := req.FormValue("timeout")
	if timeout == "" {
		send, _ := json.Marshal(resp{false, "Timeout is not configured", 1})
		fmt.Fprintf(res, string(send))
	} else {
		args = append(args, "-t")
		args = append(args, timeout)
		monkey := req.FormValue("monkey")
		if monkey == "False" {	
			Logger.WithFields(StandardFields).Info("Manual phase two configured")
		}else if monkey == "True" {
			args = append(args, "-m")
			Logger.WithFields(StandardFields).Info("Automated(monkey) phase two configured")
		}
		exitcode := RunCommand("scripts/monkey.py", args...)
		if exitcode == 0 {
			send, _ := json.Marshal(resp{true, "Successful traffic capture in phase two", 0})
			Log("Successful traffic capture in Phase II")
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Info("Successful traffic capture in phase two")
		} else {
			send, _ := json.Marshal(resp{false, "Traffic capture in phase two failed", exitcode})
			fmt.Fprintf(res, string(send))
			Logger.WithFields(StandardFields).Error("Traffic capture in phase two failed")
		}
		
	}
			
}

func analysis(res http.ResponseWriter, req *http.Request) {
	Log("Starting analysis of captured traffic")
	var ret bool
	ret = true
	//Logger.WithFields(StandardFields).Info("Starting analysis of captured traffic")
	Command(nil, "mv", "output.log", "second.privapp.log")
	csv, _ := os.OpenFile("logging/log/out.privapp.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	defer csv.Close()
	if _, err := os.Stat("first.privapp.log"); err== nil{
		Command(csv, "analyze/analyze.py", "first.privapp.log", "Phase 1", ip, name, testing_label, version)
	} else if os.IsNotExist(err){
		Logger.WithFields(StandardFields).Info("Traffic not generated in phase one.")
	}else{
		Logger.WithFields(StandardFields).Error("Error while opening phase-one traffic file!")
	}
	if _, err := os.Stat("second.privapp.log"); err== nil{
		Command(csv, "analyze/analyze.py", "second.privapp.log", "Phase 2", ip, name, testing_label, version)
	} else if os.IsNotExist(err){
		Logger.WithFields(StandardFields).Info("Traffic not generated in phase two.")
	}else{
		Logger.WithFields(StandardFields).Error("Error while opening phase-two traffic file!")
	} 

    send, _ := json.Marshal(resp{ret, "Analysis finished", 0})
    fmt.Fprintf(res, string(send))
    Logger.WithFields(StandardFields).Info("Analysis finished")
}

func result(res http.ResponseWriter, req *http.Request) {
	Log("Results reading started")
	send, err := ioutil.ReadFile("logging/log/out.privapp.log")
	if err != nil {
		send, _ = json.Marshal(resp{false, "Error reading results file", 1})
		fmt.Fprintf(res, string(send))
		Logger.WithFields(StandardFields).Error("Error reading results file")
	} else {
		fmt.Fprintf(res, string(send))
		Logger.WithFields(StandardFields).Info("Results sent")
	}
}

func screenshotPhaseOne(res http.ResponseWriter, req *http.Request) {
	Log("Phase-one screenshots reading started")
	res.Header().Set("Content-Encoding", "gzip")
        send, err := ioutil.ReadFile("fp.screenshot")
        if err != nil {
                send, _ = json.Marshal(resp{false, "Error reading screenshots phase-1", 1})
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Error("Error reading screenshots phase-1")
        } else {
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Info("Phase-one screenshot reading success")
        }
}

func screenshotPhaseTwo(res http.ResponseWriter, req *http.Request) {
	Log("Phase-two screenshots reading started")
        send, err := ioutil.ReadFile("sp.screenshoot")
        if err != nil {
                send, _ = json.Marshal(resp{false, "Error reading screenshots phase-2", 1})
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Error("Error reading screenshots phase-2")
        } else {
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Info("Phase-two screenshot reading success")
        }
}

func rawPhaseOne(res http.ResponseWriter, req *http.Request) {
	Log("Phase-one raw data reading started")
        send, err := ioutil.ReadFile("first.privapp.log")
        if err != nil {
                send, _ = json.Marshal(resp{false, "Error reading raw data of phase one", 1})
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Error("Error reading raw data of phase one")
        } else {
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Info("Phase-one raw data reading success")
        }
}

func rawPhaseTwo(res http.ResponseWriter, req *http.Request) {
	Log("Phase-two raw data reading started")
        send, err := ioutil.ReadFile("second.privapp.log")
        if err != nil {
                send, _ = json.Marshal(resp{false, "Error reading raw data of phase two", 1})
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Error("Error reading raw data of phase two")
        } else {
                fmt.Fprintf(res, string(send))
                Logger.WithFields(StandardFields).Info("Phase-two raw data reading success")
        }
}

func sanitize(res http.ResponseWriter, req *http.Request) {
	Log("Sanitization started")
        Command(nil, "scripts/kill.sh")
	Command(nil, "rm","first.privapp.log")
	//Command(nil, "rm","fp.tar.gz")
	//Command(nil, "rm","sp.tar.gz")
	Command(nil, "rm","second.privapp.log")
	Command(nil, "rm","logging/log/out.privapp.log")
	Command(nil, "scripts/uninstall.sh", ip, name)
	Logger.WithFields(StandardFields).Info("Sanitization done")
	ip, name, version = "", "", ""
	InitLogger(name, ip, LOG_FILE, testing_label, version)
	send, _ := json.Marshal(resp{true, "Sanitization done", 0})
        fmt.Fprintf(res, string(send))
}

func printState(res http.ResponseWriter, label string) {
	fmt.Fprintf(res, "\n[%s] Estado actual de variables globales:\n", label)
	fmt.Fprintf(res, "  ip            = %s\n", ip)
	fmt.Fprintf(res, "  name          = %s\n", name)
	fmt.Fprintf(res, "  testing_label = %s\n", testing_label)
	fmt.Fprintf(res, "  version       = %s\n", version)
}

