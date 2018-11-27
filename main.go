package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/fsnotify/fsnotify"
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
	ID string `mapstructure:"id"`
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
	fmt.Println("config read")
}

func startServer() {
	port := viper.GetInt32("config.http.port")
	fmt.Println("booting server on port", port)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{slug}", handleMap).Methods("GET")

	http.Handle("/", rtr)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleMap(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["slug"]
	fmt.Println("incoming request: ", name)
	w.Write([]byte("test"))
}
