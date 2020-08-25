# mixwebserver

Simple web server to serve mixer folder repository on Clear Linux, and also 3rd-part repository over  http/https. Uses Letsencrypt for cert.

The binary "mixwebserver" are compiled for use on Linux amd64 platform.

```text
  -dir3rdParty string
        specify the full path to the folder to serve as 3rd-party repository (default "./3rd-party")
  -mail string
        The mail address to use when registering domain for letsencrypt cert
  -path string
        enter the full path of the directory to serve via https (default "./")
  -prod
        set to true if you want a real and signed certificate, and are done with testing
  -readTimeout int
        the number of minutes for the http request timeout (default 120)
  -url string
        Enter the url for your domain
```
