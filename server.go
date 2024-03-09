package main

import (
    "fmt"
    "strconv"
    "net/http"
    "os"
    "math"
    "time"
    "os/exec"
    "strings"
    "github.com/mackerelio/go-osstat/cpu"
    "github.com/mackerelio/go-osstat/loadavg"
)


func loadavg_pct(w http.ResponseWriter, req *http.Request) {
        pct := ""
 	load_pct, exists := cache.Get("load_pct")
 	if ! exists {
            after, err := cpu.Get()
       	    if err != nil {
                fmt.Fprintf(os.Stderr, "E2%s\n", err)
                return
            }
           load, err := loadavg.Get()
           if err != nil {
                fmt.Fprintf(os.Stderr, "E1%s\n", err)
                return
           }

           pct = strconv.Itoa(int(math.Round(100.0*load.Loadavg1/float64(after.CPUCount))))
           cache.Set("load_pct", pct, cache_timeout*time.Second)
           fmt.Fprintf(w, "%s%%\n", strings.TrimSpace(string(pct)))
	
       	} else {
         
          fmt.Fprintf(w, "%s%%\n", strings.TrimSpace(string(load_pct.(string))))
        }

}
func loadavg_free_pct(w http.ResponseWriter, req *http.Request) {
        load_pct_free, exists := cache.Get("load_pct_free")
        if ! exists {
            after, err := cpu.Get()
            if err != nil {
                fmt.Fprintf(os.Stderr, "E2%s\n", err)
                return
            }
           load, err := loadavg.Get()
           if err != nil {
                fmt.Fprintf(os.Stderr, "E1%s\n", err)
                return
           }
   	   load_pct_free := 100-int(math.Round(100.0*load.Loadavg1/float64(after.CPUCount)))
           load_pct_free_max := max(load_pct_free,1)
           pct_free_string := strconv.Itoa(load_pct_free_max)
           cache.Set("load_pct_free", pct_free_string, cache_timeout*time.Second)
           fmt.Fprintf(w, "%s%%\n", strings.TrimSpace(string(pct_free_string)))

       } else {

          fmt.Fprintf(w, "%s%%\n", strings.TrimSpace(string(load_pct_free.(string))))
       }

}


func cpu_idle(w http.ResponseWriter, req *http.Request)  {
    cpu_idle, exists := cache.Get("cpu_idle")
    if ! exists {
    	before, err := cpu.Get()
        if err != nil {
                fmt.Fprintf(os.Stderr, "%s\n", err)
                return
        }
        time.Sleep(time.Duration(250) * time.Millisecond)
        after, err := cpu.Get()
        if err != nil {
                fmt.Fprintf(os.Stderr, "%s\n", err)
                return
        }
        total := float64(after.Total - before.Total)
        cpu_idle := strconv.Itoa(int(math.Round(float64(after.Idle-before.Idle)/total*100)))
        cache.Set("cpu_idle", cpu_idle, cache_timeout*time.Second)
        fmt.Fprintf(w, "%s%%\n", strings.TrimSpace(string(cpu_idle)))
    } else {
	 fmt.Fprintf(w, "%s%%\n", strings.TrimSpace(string(cpu_idle.(string))))
    }

}
func pepe(w http.ResponseWriter, req *http.Request)  {
	
	 fmt.Fprintf(w,"hello\n")
         cmd := exec.Command("uptime")
	output, err := cmd.CombinedOutput()
        if err != nil {
		fmt.Printf("Error getting system load: %v\n", err)
	}
        fmt.Fprintf(w,  strings.TrimSpace(string(output)))

   }

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

var cache = NewCache() 
var cache_timeout time.Duration = 10
var default_port string = ":7001"
func main() {
    port := default_port
    if len(os.Args) > 1 {
     	inputNumber := os.Args[1]
    	number, err := strconv.Atoi(inputNumber)
    	
	if err != nil {
                fmt.Println("Not a valid port.Using Default "+default_port, err)
    	} else {
    		port = fmt.Sprintf(":%d", number)
    	}
    }
    http.HandleFunc("/loadavg_free_pct", loadavg_free_pct)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/loadavg_pct", loadavg_pct)
    http.HandleFunc("/cpu_idle", cpu_idle)
    fmt.Println("Start Listening on port "+port)
    http.ListenAndServe(port, nil)
    fmt.Println("Exiting... ")
}
