package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config") // name of config file (without extension)
	//viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".") // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	url := "https://api.ipify.org?format=text"
	fmt.Printf("Getting IP address from  ipify\n")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("IP is:%s\n", ip)

	// Construct a new API object
	api, err := cloudflare.New(viper.GetString("token"), viper.GetString("username"))
	if err != nil {
		log.Fatal(err)
	}

	cur, err := api.DNSRecord(viper.GetString("zone_id"), viper.GetString("dns_id"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cur.Content)

	if cur.Content == string(ip) {
		fmt.Println("record and current IP are identical")
	} else {
		fmt.Println("record and current IP are different")

		res := api.UpdateDNSRecord(viper.GetString("zone_id"), viper.GetString("dns_id"), cloudflare.DNSRecord{Content: string(ip)})
		if res == nil {
			fmt.Println("Record Updated")
		} else {
			fmt.Println("FAILURE")
			fmt.Println(res)
		}
	}

}
