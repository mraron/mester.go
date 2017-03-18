package main

import (
	"log"
	"os"
	"time"
	"crypto/tls"
	"net/http"
	"strconv"
	"encoding/json"
)

type Solution struct {
	Statement string
	Topic     string
	Problem   string
	Name      string
	Point     int
}

type Crawler struct {
	client *http.Client
	sols []Solution
	getSolutions bool
	exportProblems bool
}

func NewCrawler(getSolutions, exportProblems bool) *Crawler {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	return &Crawler{client, make([]Solution,0),getSolutions,exportProblems}
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
		for j := 1; j <= problemCount; j++ {
			log.Print("Getting problem ",j)
			ChooseProblem(c.client, j)
			problem := GetProblemName(c.client)
			statement := strconv.Itoa(id) + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j) + ".pdf"
			if(c.exportProblems) {
				o, err := os.Create(statement)
				if err != nil {
					return err
				}

				o.Write(GetStatement(c.client))
				o.Close()
			}

			if(c.getSolutions) {
				for _, val := range GetSolvers(c.client) {
					c.sols = append(c.sols, Solution{statement, topic, problem, val.Name, val.Point})
				}
			}

			time.Sleep(time.Millisecond * 50)
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

