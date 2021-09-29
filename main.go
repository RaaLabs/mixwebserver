// Web server to be used for serving a single mix update folder.
// Will automatically get a lets encrypt certificate, and use
// that to provide https transport for the domain.

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caddyserver/certmagic"
)

func main() {
	dir := flag.String("dirMain", "./", "enter the full path of the main mix  update directory to serve via https")
	dir3rdParty := flag.String("dir3rdParty", "./3rd-party", "specify the full path to the update folder to serve as 3rd-party repository")
	prod := flag.Bool("prod", false, "set to true if you want a real and signed certificate, and are done with testing")
	url := flag.String("url", "", "Enter the url for your domain")
	mail := flag.String("mail", "", "The mail address to use when registering domain for letsencrypt cert")
	readTimeout := flag.Int("readTimeout", 120, "the number of minutes for the http request timeout")

	flag.Parse()

	// Create a file server, and serve the main mix files given by the -dir flag
	fd := http.FileServer(http.Dir(*dir))
	http.Handle("/", http.StripPrefix("/", fd))

	// Create a file server for the 3rd-party repository files
	fd3rdParty := http.FileServer(http.Dir(*dir3rdParty))
	http.Handle("/3rd-party/", http.StripPrefix("/3rd-party/", fd3rdParty))

	switch *url {
	case "":
		log.Println("You need to specify a domain using the -url flag")
		return
	}

	switch *mail {
	case "":
		log.Println("You need to specify a mail address using the -mail flag")
		return
	}

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = *mail

	switch *prod {
	case true:
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	case false:
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	default:
		log.Println("Info: you need to set the prod flag to true or false")
		return
	}

	if *prod {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	} else {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	}

	// Serve HTTP

	httpLn, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Printf("error: certmagic.Listen http failed: %v ", err)
	}

	httpConf := http.Server{
		ReadTimeout: time.Duration(*readTimeout) * time.Minute,
		Handler:     nil,
		Addr:        ":80",
	}

	go httpConf.Serve(httpLn)

	// Serve HTTPS

	httpsLn, err := certmagic.Listen([]string{*url})
	if err != nil {
		log.Printf("error: certmagic.Listen https failed: %v ", err)
	}

	httpsConf := http.Server{
		// ReadTimeout: 2 * time.Minute,
		ReadTimeout: time.Duration(*readTimeout) * time.Minute,
		Handler:     nil,
		Addr:        ":443",
	}

	if err := httpsConf.Serve(httpsLn); err != nil {
		log.Printf("error: httpsConf.Serve failed: %v ", err)
	}

}
