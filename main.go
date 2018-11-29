package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var port *int
var outputFile *string
var configFile *string
var file *os.File

func prepareFlags() {
	port = flag.Int("port", 8000, "port to run http service on")
	outputFile = flag.String("output", "output.log", "file to log output to")
	configFile = flag.String("cfg", "mappings.yaml", "config file")
	flag.Parse()
}

func prepareLogger(output string) {
	file, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// send log to file and stdout
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}

func main() {

	// prepare app for booting
	prepareFlags()
	prepareLogger(*outputFile)

	defer log.Println("Service: Shutting Down")
	defer log.Println("------")
	defer file.Close()
	log.Println("Service: Starting")

	loadConfig()
	startServer()
}

// Config
type Config struct {
	Maps []Map `mapstructure:"maps"`
}

// Map
type Map struct {
	ID       string `mapstructure:"id"`
	Token    string
	Binaries []Bin
}

type Bin struct {
	Name    string
	Handler string
	Request map[string]string
}

func (m Map) isValid(token string) bool {
	return m.Token == token
}

func (b Bin) buildCommand(request *http.Request) (string, error) {
	base := b.Handler
	data := extractDataFromRequest(b.Request, request)

	fmt.Println("data", data)
	// search within the base string, find any markers
	return base, nil
}

func extractDataFromRequest(mapping map[string]string, r *http.Request) map[string]string {
	var data = make(map[string]string)
	for key, value := range mapping {
		s := strings.Split(key, ":")
		rType, keyName := strings.ToLower(s[0]), s[1]
		if rType == "get" {
			val := r.URL.Query()[keyName]
			if len(val) > 0 {
				data[value] = val[0]
			}
		}
	}
	return data
}

var config Config

func loadConfig() {
	viper.SetConfigFile(*configFile)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		viper.Unmarshal(&config)
	})
	viper.Unmarshal(&config)
	log.Println("Config: read")
}

func startServer() {
	log.Println("Server: booting on port", strconv.Itoa(*port))

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{slug}", handleMap).Methods("GET")

	http.Handle("/", rtr)

	portString := ":" + strconv.Itoa(*port)

	if err := http.ListenAndServe(portString, nil); err != nil {
		panic(err)
	}
}

func handleMap(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["slug"]

	log.Printf("Request: handling: /%s\n", name)
	m := findMap(name, config.Maps)
	if len(m.ID) > 0 {
		token, ok := r.URL.Query()["token"]
		if ok && m.isValid(token[0]) {
			go func() {
				// call binaries
				for _, bin := range m.Binaries {
					log.Println(bin.Name)
					cmd, err := bin.buildCommand(r)
					if err != nil {
						log.Println("cant build command")
					}
					exe := exec.Command(cmd)
					err = exe.Run()
					if err != nil {
						log.Println(err)
					}
				}
			}()
		} else {
			log.Printf("Request: rejected /%s (403)", name)
			w.WriteHeader(http.StatusForbidden)
		}

	} else {
		log.Println("Request: rejected: /%s (404)", name)
		w.WriteHeader(http.StatusNotFound)
	}
}

func findMap(name string, maps []Map) Map {
	for _, element := range maps {
		if element.ID == name {
			return element
		}
	}
	return Map{}
}
