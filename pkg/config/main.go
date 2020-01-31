// Copyright (C) 2019 LEAP
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"flag"
	"log"
	"os"
)

const DefaultAuthenticationModule string = "anon"

type Opts struct {
	Tls        bool
	CaCrt      string
	CaKey      string
	TlsCrt     string
	TlsKey     string
	Port       string
	Auth       string
	AuthSecret string
}

func FallbackToEnv(variable *string, envVar, defaultVar string) {

	if *variable == "" {
		val, exists := os.LookupEnv(envVar)
		if exists && val != "" {
			*variable = val
		} else {
			*variable = defaultVar
		}
	}
}

func doCaFilesSanityCheck(caCrt string, caKey string) {
	if _, err := os.Stat(caCrt); os.IsNotExist(err) {
		log.Fatal("cannot find caCrt file")
	}
	if _, err := os.Stat(caKey); os.IsNotExist(err) {
		log.Fatal("cannot find caKey file")
	}
}

func doTlsFilesSanityCheck(tlsCrt string, tlsKey string) {
	if _, err := os.Stat(tlsCrt); os.IsNotExist(err) {
		log.Fatal("cannot find tlsCrt file")
	}
	if _, err := os.Stat(tlsKey); os.IsNotExist(err) {
		log.Fatal("cannot find tlsKey file")
	}
}

func NewOpts() *Opts {
	opts := new(Opts)
	initializeFlags(opts)
	checkConfigurationOptions(opts)
	return opts
}

func initializeFlags(opts *Opts) {
	flag.StringVar(&opts.CaCrt, "caCrt", "", "Path to the CA public key used for VPN certificates")
	flag.StringVar(&opts.CaKey, "caKey", "", "Path to the CA private key used for VPN certificates")
	flag.BoolVar(&opts.Tls, "tls", false, "Enable TLS on the service")
	flag.StringVar(&opts.TlsCrt, "tlsCrt", "", "Path to the cert file for TLS")
	flag.StringVar(&opts.TlsKey, "tlsKey", "", "Path to the key file for TLS")
	flag.StringVar(&opts.Port, "port", "", "Port where the server will listen (default: 8000)")
	flag.StringVar(&opts.Auth, "auth", "", "Authentication module (anonymous, sip)")
	flag.StringVar(&opts.AuthSecret, "authSecret", "", "Authentication secret (optional)")
	flag.Parse()

	FallbackToEnv(&opts.CaCrt, "VPNWEB_CACRT", "")
	FallbackToEnv(&opts.CaKey, "VPNWEB_CAKEY", "")
	FallbackToEnv(&opts.TlsCrt, "VPNWEB_TLSCRT", "")
	FallbackToEnv(&opts.TlsKey, "VPNWEB_TLSKEY", "")
	FallbackToEnv(&opts.Port, "VPNWEB_PORT", "8000")
	FallbackToEnv(&opts.Auth, "VPNWEB_AUTH", DefaultAuthenticationModule)
	FallbackToEnv(&opts.AuthSecret, "VPNWEB_AUTH_SECRET", "")
}

func checkConfigurationOptions(opts *Opts) {
	if opts.CaCrt == "" {
		log.Fatal("missing caCrt parameter")
	}
	if opts.CaKey == "" {
		log.Fatal("missing caKey parameter")
	}

	if opts.Tls == true {
		if opts.TlsCrt == "" {
			log.Fatal("missing tls_crt parameter")
		}
		if opts.TlsKey == "" {
			log.Fatal("missing tls_key parameter")
		}
	}

	doCaFilesSanityCheck(opts.CaCrt, opts.CaKey)
	if opts.Tls == true {
		doTlsFilesSanityCheck(opts.TlsCrt, opts.TlsKey)
	}

	log.Println("Authentication module:", opts.Auth)
}
