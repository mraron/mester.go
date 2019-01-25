package main

import (
	"log"
	"os"
	"crypto/tls"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
)

type PointHistoryElem struct {
	Point int
	Time time.Time
}


type Identifier struct {
	Topic     string
	Problem   string
	Name      string
}

type Solution struct {
	Statement string
	Identifier
	Point     int
	PointHistory []PointHistoryElem
}



type Crawler struct {
	client *http.Client
	sols []Solution
	lookup map[Identifier]int
	getSolutions bool
	exportProblems bool
}

func NewCrawler(getSolutions, exportProblems bool) *Crawler {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	
	sols := make([]Solution,0)
	lookup := make(map[Identifier]int)
	
	latest, err := os.Open("data.json")
	if err == nil {
		defer latest.Close()
		
		dec := json.NewDecoder(latest)
		err = dec.Decode(&sols)
		if err != nil {
			panic(err)
		}
		
		for ind, sol := range sols {
			lookup[sol.Identifier]=ind
		}
	}
	
	
	return &Crawler{client, sols, lookup, getSolutions,exportProblems}
}

func (c* Crawler) Crawl(id, L, R int) (error){
	log.Print("Crawling ",id," from ",L," to ",R)
	ChooseLevelAndTopic(c.client, id, L, GetViewState(c.client, "https://mester.inf.elte.hu/faces/tema.xhtml"), true)
	for i := L; i <= R; i++ {
		log.Print("Topic #",i)
		val := GetViewState(c.client, "https://mester.inf.elte.hu/faces/tema.xhtml")
		ChooseLevelAndTopic(c.client, id, i, val, false)
		problemCount := GetProblemNumber(c.client)

		log.Print("Found ",problemCount, " problems")

		topic := GetTopicName(c.client)
		log.Print("Name of topic is ", topic);
		for j := 1; j <= problemCount; j++ {
			log.Print("Getting problem ",j)
			ChooseProblem(c.client, j)
			problem := GetProblemName(c.client)
			statement := strconv.Itoa(id) + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j) + ".pdf"
			if(c.exportProblems) {
				_ = os.Mkdir("statements/", 0777)
				o, err := os.Create("statements/"+statement)
				if err != nil {
					return err
				}

				o.Write(GetStatement(c.client))
				o.Close()
			}

			if(c.getSolutions) {
				for _, val := range GetSolvers(c.client) {
					ident := Identifier{topic, problem, val.Name}
					if _, ok := c.lookup[ident]; !ok {
						c.sols = append(c.sols, Solution{statement, ident, val.Point, make([]PointHistoryElem, 1)})
						c.sols[len(c.sols)-1].PointHistory[0]=PointHistoryElem{val.Point, time.Now()};
						c.lookup[ident]=len(c.sols)-1
					}else {
						if c.sols[c.lookup[ident]].Point != val.Point {
							c.sols[c.lookup[ident]].Point = val.Point
							c.sols[c.lookup[ident]].PointHistory = append(c.sols[len(c.sols)-1].PointHistory, PointHistoryElem{val.Point, time.Now()});
						}
					}
				}
			}
		}
	}

	return nil
}

func (c *Crawler) Export() (error) {
	w, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer w.Close()

	enc := json.NewEncoder(w)
	return enc.Encode(&c.sols)
}

