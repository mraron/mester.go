package main

import (
	"log"
	"net/http"
	"flag"
)

var auth = &http.Cookie{Name: "JSESSIONID", Value: "SUTI"}

var getSolutions = flag.Bool("getSolutions", true, "megoldásokat letöltése")
var exportProblems = flag.Bool("exportProblems", false, "pdf fájlokat az aktuális mappában hozzon-e létre")
var topicKezdo = flag.Bool("topicKezdo", false, "kezdő témát töltsön-e")
var topicKozephalado = flag.Bool("topicKozephalado", false, "középhaladó témát töltsön-e")
var topicHalado = flag.Bool("topicHalado", true, "haladó témát töltsön-e")
var topicNT2017 = flag.Bool("topicNT2017", true, "NT/OKTV/válogató témát töltsön-e")
var topicKomal = flag.Bool("topicKomal", true, "kömal témát töltsön-e")
var topicAll = flag.Bool("topicAll", false, "minden témát töltsön-e")
var cookie = flag.String("cookie", "", "a belépéskor kapott JSESSIONID süti értéke [SZÜKSÉGES]")

func main() {
	flag.Parse()
	if len(*cookie)==0 {
		log.Fatal("-help a te barátod")
	}
	auth = &http.Cookie{Name: "JSESSIONID", Value: *cookie}

	c := NewCrawler(*getSolutions, *exportProblems)
	var err error

	//kezdő
	if *topicKezdo || *topicAll {
		err = c.Crawl(0,1,12)
		if err != nil {
			log.Fatal(err)
		}
	}

	//középhaladó
	if *topicKozephalado || *topicAll {
		err = c.Crawl(1,13,25)
		if err != nil {
			log.Fatal(err)
		}
	}

	//haladó
	if *topicHalado || *topicAll {
		err = c.Crawl(2, 26, 46)
		if err != nil {
			log.Fatal(err)
		}
	}

	//nt/oktv/válogató
	if *topicNT2017 || *topicAll {
		err = c.Crawl(3,47,60)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *topicKomal || *topicAll {
	    err = c.Crawl(4, 61, 61)
	    if err != nil {
			log.Fatal(err)
	    }
	}

	err = c.Export()
	if err != nil {
		log.Fatal(err)
	}
}
