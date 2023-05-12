package main

import (
	"log"

	"github.com/gmalka/MyServer/internal/app"
	"github.com/spf13/viper"
	_ "github.com/onsi/ginkgo/v2/ginkgo"
)

type g struct {

}

func main() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}