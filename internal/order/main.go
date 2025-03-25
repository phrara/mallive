package main

import (
	"log"
	_ `net/http`

	"github.com/phrara/mallive/common/config"
	"github.com/spf13/viper"
)


func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}


func main() {

	log.Println(viper.GetString("order.server.address"))

	
	// mux := http.NewServeMux()
	// mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, world!"))
	// })
	
	// log.Println("Listening on 8089!")
	// if err := http.ListenAndServe(":8089", mux); err != nil {
	// 	log.Fatal(err)
	// }

}