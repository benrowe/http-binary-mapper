package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var port *int

func main() {
	port = flag.Int("port", 8000, "port to run http service on")
	outputFile := flag.String("output", "output.log", "file to log output to")
	file, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	log.SetOutput(file)

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
}

func (m Map) isValid(token string) bool {
	return m.Token == token
}

var config Config

func loadConfig() {
	viper.SetConfigFile(*flag.String("cfg", "config.yaml", "config file name"))
	// viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// viper.SetConfigType("yaml")

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
	log.Println("Server: booting on port", port)

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

	log.Println("Request: handling: ", name)
	m := findMap(name, config.Maps)
	if len(m.ID) > 0 {
		token, ok := r.URL.Query()["token"]
		if ok && m.isValid(token[0]) {
			go func() {
				// call binaries
				for _, bin := range m.Binaries {
					log.Println(bin.Name)
					exe := exec.Command(bin.Handler)
					err := exe.Run()
					log.Println(err)
				}
			}()
		} else {
			log.Println("Request: rejected ", name, "(403)")
			w.WriteHeader(http.StatusForbidden)
		}

	} else {
		log.Println("Request: rejected: ", name, "(404)")
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
