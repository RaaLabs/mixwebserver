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
	dir := flag.String("path", "./", "enter the full path of the directory to serve via https")
	prod := flag.Bool("prod", false, "set to true if you want a real and signed certificate, and are done with testing")
	url := flag.String("url", "", "Enter the url for your domain")
	mail := flag.String("mail", "", "The mail address to use when registering domain for letsencrypt cert")
	readTimeout := flag.Int("readTimeout", 120, "the number of minutes for the http request timeout")

	flag.Parse()

	// Create a file server, and serve the files given by the -dir flag
	fd := http.FileServer(http.Dir(*dir))
	http.Handle("/", fd)

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

	// Get certificates and start https
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

	httpsLn, err := certmagic.Listen([]string{*url})
	if err != nil {
		log.Printf("error: certmagic.Listen https failed: %v ", err)
	}

	httpsConf := http.Server{
		ReadTimeout: 2 * time.Minute,
		Handler:     nil,
		Addr:        ":443",
	}

	if err := httpsConf.Serve(httpsLn); err != nil {
		log.Printf("error: httpsConf.Serve failed: %v ", err)
	}

}
