// Web server to be used for serving a single mix update folder.
// Will automatically get a lets encrypt certificate, and use
// that to provide https transport for the domain.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caddyserver/certmagic"
)

func main() {
	dir := flag.String("path", "./", "enter the full path of the directory to serve via https")
	prod := flag.Bool("prod", false, "set to true if you want a real and signed certificate, and are done with testing")
	url := flag.String("url", "", "Enter the url for your domain")
	mail := flag.String("mail", "", "The mail address to use when registering domain for letsencrypt cert")

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

	err := certmagic.HTTPS([]string{*url}, nil)
	if err != nil {
		log.Println("--- error: certmagic.HTTPS failed: ", err)
		return
	}
}
