package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
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
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		viper.Unmarshal(&config)
	})
	viper.Unmarshal(&config)
	fmt.Println("Config: read")
}

func startServer() {
	port := viper.GetInt32("config.http.port")
	fmt.Println("Server: booting on port", port)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{slug}", handleMap).Methods("GET")

	http.Handle("/", rtr)

	portString := ":" + viper.GetString("config.http.port")

	if err := http.ListenAndServe(portString, nil); err != nil {
		panic(err)
	}
}

func handleMap(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["slug"]

	fmt.Println("Request: handling: ", name)
	m := findMap(name, config.Maps)
	if len(m.ID) > 0 {
		token, ok := r.URL.Query()["token"]
		if ok && m.isValid(token[0]) {
			go func() {
				// call binaries
				for _, bin := range m.Binaries {
					fmt.Println(bin.Name)
					exe := exec.Command(bin.Handler)
					err := exe.Run()
					fmt.Println(err)
				}
			}()
		} else {
			fmt.Println("Request: rejected ", name, "(403)")
			w.WriteHeader(http.StatusForbidden)
		}

	} else {
		fmt.Println("Request: rejected: ", name, "(404)")
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
