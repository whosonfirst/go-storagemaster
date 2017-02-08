package main

import (
	"flag"
	_ "fmt"
	"github.com/whosonfirst/go-storagemaster"
	"github.com/whosonfirst/go-storagemaster/provider"
	"io/ioutil"
	"log"
	"os"		
	"strings"
)

type Params []string

func (p *Params) String() string {
	return strings.Join(*p, "\n")
}

func (p *Params) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func (p *Params) ToExtras() (*storagemaster.StoragemasterExtras, error) {

	extras, err := storagemaster.NewStoragemasterExtras()

	if err != nil {
	   return nil, err
	}
	
	for _, str_pair := range *p {
		pair := strings.Split(str_pair, "=")
		extras.Set(pair[0], pair[1])
	}

	return extras, nil
}

func main() {

	var params Params

	flag.Var(&params, "param", "Zero or more query=value parameters.")

	var sm_provider = flag.String("provider", "s3", "...")

	var s3_credentials = flag.String("s3-credentials", "", "...")
	var s3_bucket = flag.String("s3-bucket", "", "...")
	var s3_prefix = flag.String("s3-prefix", "", "...")
	var s3_region = flag.String("s3-region", "", "...")

	flag.Parse()

	var sm storagemaster.Provider
	var err error
	
	if *sm_provider == "s3" {

		cfg := provider.S3Config{
			Bucket:      *s3_bucket,
			Prefix:      *s3_prefix,
			Region:      *s3_region,
			Credentials: *s3_credentials,
		}

		sm, err = provider.NewS3Provider(cfg)

		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal("Unknown provider")
	}

	args := flag.Args()
	
	if len(args) < 2 {
		log.Fatal("Insufficient arguments")
	}

	verb := strings.ToUpper(args[0])
	
	if verb == "GET" {

		for _, key := range args[1:] {

			bytes, err := sm.Get(key)

			if err != nil {
				log.Fatal(err)
			}

			log.Println(string(bytes))
		}

	} else if verb == "EXISTS" {

		for _, key := range args[1:] {

			exists, err := sm.Exists(key)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Does %s exist: %t\n", key, exists)
		}

	} else if verb == "PUT" {

		if len(args) != 3 {
		   log.Fatal("Invalid PUT args")
		}

		extras, err := params.ToExtras()

		if err != nil {
		   log.Fatal(err)
		}

		src := args[1]
		dest := args[2]

		fh, err := os.Open(src)

		if err != nil {
		   log.Fatal(err)
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
		   log.Fatal(err)
		}

		err = sm.Put(dest, body, extras)

		if err != nil {
		   log.Fatal(err)
		}

	/*
	} else if verb == "DELETE" {
	*/
	
	} else {
		log.Fatal("Unsupported verb")
	}
	
}
