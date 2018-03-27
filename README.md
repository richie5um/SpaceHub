# SpaceHub

## What is it?
A super-simple SSL-terminating reverse-proxy.

## What?
Sometimes you have services that expose HTTP endpoints, but want external clients to connect using HTTPS. This may be your own design, or, because your services are normally hosted behind an ApplicationGateway instance that does SSL-termination for you.

This app aims to provide a very basic implemtation of that pattern - there are *many* more alternatives to this but often they are _so_ flexible they actually make things complicated. 

"Eliminate flexibility to provide simplicity."

## How
* Download and build the source

    git clone ...
    go get
    go build

* Or download a release
* Run

    SpaceHub -port=443 -targetURL=http://localhost:80
  
## Certificates
SpaceHub requires a 'cert.crt' and a 'cert.key' file.

> Note: the 'cert.key' must be in unprotected format.

#### How to convert a PFX file to .crt and .key
Run:
* openssl pkcs12 -in cert.pfx -nocerts -out cert.encrypted.key 
* openssl pkcs12 -in cert.pfx -clcerts -nokeys -out cert.crt
* openssl rsa -in cert.encrypted.key -out cert.key

## Contributions
Feedback / PRs welcome.
