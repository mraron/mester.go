package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var auth = &http.Cookie{Name: "JSESSIONID", Value: "8b429976f02fc7808e8c399255a7"}

func ChooseLevelAndTopic(c *http.Client, id int, topic int, str string, isreq bool) (string, error) {
	form := make(url.Values)
	form.Add("form", "form")
	form.Add("form:name", strconv.Itoa(id))
	form.Add("form:temalist", strconv.Itoa(topic))
	form.Add("form:j_idt16", "választom")
	form.Add("javax.faces.ViewState", str)
	form.Add("javax.faces.source", "form:name")
	form.Add("javax.faces.partial.event", "change")
	form.Add("javax.faces.partial.execute", "form:name form:name")
	form.Add("javax.faces.partial.render", "form:temalist")
	form.Add("javax.faces.behavior.event", "change")
	form.Add("javax.faces.partial.ajax", "ajax")

	req, err := http.NewRequest("POST", "https://mester.inf.elte.hu/faces/tema.xhtml", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	//fmt.Println(req.PostForm.Encode())
	if isreq {
		req.Header.Set("Faces-Request", "partial/ajax")
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/tema.xhtml")
	req.Header.Set("Host", "mester.inf.elte.hu")
	req.Header.Set("Origin", "https://mester.inf.elte.hu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.AddCookie(auth)

	resp, err := c.Do(req)
	ba, err := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(ba))
	return string(ba), err
}

func GetViewState(c *http.Client, url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	val := ""
	d.Find("input[name=\"javax.faces.ViewState\"]").Each(func(i int, s *goquery.Selection) {
		val, _ = s.Attr("value")
	})

	return val
}

func GetProblemNumber(c *http.Client) int {
	req, err := http.NewRequest("POST", "https://mester.inf.elte.hu/faces/feladat.xhtml", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	mx := 0
	d.Find("select").Children().Each(func(i int, s *goquery.Selection) {
		val_, _ := s.Attr("value")
		val, _ := strconv.Atoi(val_)
		if mx < val {
			mx = val
		}
	})
	return mx
}

func ChooseProblem(c *http.Client, i int) string {
	vs := GetViewState(c, "https://mester.inf.elte.hu/faces/feladat.xhtml")

	form := make(url.Values)
	form.Add("form", "form")
	form.Add("form:name", strconv.Itoa(i))
	form.Add("form:j_idt13", "választom")
	form.Add("javax.faces.ViewState", vs)

	req, err := http.NewRequest("POST", "https://mester.inf.elte.hu/faces/feladat.xhtml", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/feladat.xhtml")
	req.Header.Set("Host", "mester.inf.elte.hu")
	req.Header.Set("Origin", "https://mester.inf.elte.hu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	val := ""
	d.Find("input[name=\"javax.faces.ViewState\"]").Each(func(i int, s *goquery.Selection) {
		val, _ = s.Attr("value")
	})

	return val
}

func GetStatement(c *http.Client) []byte {
	vs := GetViewState(c, "https://mester.inf.elte.hu/faces/feladat.xhtml")

	form := make(url.Values)
	form.Add("j_idt9", "j_idt9")
	form.Add("j_idt9:j_idt10", "megjelenit")
	form.Add("javax.faces.ViewState", vs)
	req, err := http.NewRequest("POST", "https://mester.inf.elte.hu/faces/leiras.xhtml", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/feladat.xhtml")
	req.Header.Set("Host", "mester.inf.elte.hu")
	req.Header.Set("Origin", "https://mester.inf.elte.hu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

type Solver struct {
	Name  string
	Point int
}

func GetSolvers(c *http.Client) []Solver {
	req, err := http.NewRequest("POST", "https://mester.inf.elte.hu/faces/eredmenylista.xhtml", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/eredmenylista.xhtml")
	req.Header.Set("Host", "mester.inf.elte.hu")
	req.Header.Set("Origin", "https://mester.inf.elte.hu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ans := make([]Solver, 0)
	d.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		name, point := "", 0

		s.Find("td").Each(func(j int, s *goquery.Selection) {
			if name == "" {
				name = strings.TrimSpace(s.Text())
			} else {
				var err error
				point, err = strconv.Atoi(strings.TrimSpace(s.Text()))
				if err != nil {
					log.Fatal(err)
				}
			}
		})

		ans = append(ans, Solver{name, point})
	})

	return ans
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func GetTopicName(c *http.Client) string {
	req, err := http.NewRequest("GET", "https://mester.inf.elte.hu/faces/eredmenylista.xhtml", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/eredmenylista.xhtml")
	req.Header.Set("Host", "mester.inf.elte.hu")
	req.Header.Set("Origin", "https://mester.inf.elte.hu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(d.Text())
	ans := ""
	d.Find("h1").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			return
		}
		//topic, name := "", ""
		//fmt.Sscanf(s.Text(), "Téma: %s, Feladat: %s", &topic, &name)
		ans = strings.Split(strings.Split(s.Text(), ", Feladat:")[0], "Téma: ")[1]
	})

	return ans
}

func GetProblemName(c *http.Client) string {
	req, err := http.NewRequest("GET", "https://mester.inf.elte.hu/faces/eredmenylista.xhtml", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/eredmenylista.xhtml")
	req.Header.Set("Host", "mester.inf.elte.hu")
	req.Header.Set("Origin", "https://mester.inf.elte.hu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.AddCookie(auth)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	d, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(d.Text())
	ans := ""
	d.Find("h1").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			return
		}
		//topic, name := "", ""
		//fmt.Sscanf(s.Text(), "Téma: %s, Feladat: %s", &topic, &name)
		ans = strings.Split(s.Text(), ", Feladat: ")[1]
	})

	return ans
}

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

//JSESSIONID=5b5b92d950870f5d57fb596c8132
func main() {

	fmt.Println(time.Now())
	//	val := GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml")

	/*ChooseLevelAndTopic(client, 0, 1, val, false)
	ChooseProblem(client, 1)
	o, err := os.Create("teszt.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()
	o.Write(GetStatement(client))
	*/
	//sols := make([]Solution, 0)
	c := NewCrawler(true, false)
	err := c.Crawl(0,1,12)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Crawl(1,13,21)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Crawl(2,22,42)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Crawl(3,43,48)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Export()
	if err != nil {
		log.Fatal(err)
	}

	/*
	ChooseLevelAndTopic(client, 1, 13, GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml"), true)
	for i := 13; i <= 21; i++ {
		val := GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml")
		fmt.Println(ChooseLevelAndTopic(client, 1, i, val, false))
		problemnum := GetProblemNumber(client)
		//fmt.Println(problemnum)
		topic := GetTopicName(client)
		for j := 1; j <= problemnum; j++ {
			ChooseProblem(client, j)
			problem := GetProblemName(client)
			statement := "1" + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j) + ".pdf"
			o, err := os.Create(statement)
			if err != nil {
				log.Fatal(err)
			}
			
			o.Write(GetStatement(client))
			o.Close()
			for _, val := range GetSolvers(client) {
				sols = append(sols, Solution{statement, topic, problem, val.Name, val.Point})
			}
			time.Sleep(time.Millisecond * 50)
		}
	}
	ChooseLevelAndTopic(client, 2, 22, GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml"), true)
	for i := 22; i <= 42; i++ {
		val := GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml")
		fmt.Println(ChooseLevelAndTopic(client, 2, i, val, false))
		problemnum := GetProblemNumber(client)
		//fmt.Println(problemnum)
		topic := GetTopicName(client)
		for j := 1; j <= problemnum; j++ {
			ChooseProblem(client, j)
			problem := GetProblemName(client)
			statement := "2" + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j) + ".pdf"
			o, err := os.Create(statement)
			if err != nil {
				log.Fatal(err)
			}
			
			o.Write(GetStatement(client))
			o.Close()
			
			for _, val := range GetSolvers(client) {
				sols = append(sols, Solution{statement, topic, problem, val.Name, val.Point})
			}
			time.Sleep(time.Millisecond * 50)
		}
	}
	ChooseLevelAndTopic(client, 3, 43, GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml"), true)
	for i := 43; i <= 48; i++ {
		val := GetViewState(client, "https://mester.inf.elte.hu/faces/tema.xhtml")
		fmt.Println(ChooseLevelAndTopic(client, 3, i, val, false))
		problemnum := GetProblemNumber(client)
		//fmt.Println(problemnum)
		topic := GetTopicName(client)
		for j := 1; j <= problemnum; j++ {
			ChooseProblem(client, j)
			problem := GetProblemName(client)
			statement := "3" + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(j) + ".pdf"
			o, err := os.Create(statement)
			if err != nil {
				log.Fatal(err)
			}
			
			o.Write(GetStatement(client))
			o.Close()
			
			for _, val := range GetSolvers(client) {
				sols = append(sols, Solution{statement, topic, problem, val.Name, val.Point})
			}
			time.Sleep(time.Millisecond * 50)
		}
	}
	
	w, err := os.Create("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	enc := json.NewEncoder(w)
	err = enc.Encode(&sols)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("----------------------------------")
	//fmt.Println(GetProblemNumber(client))
	//ChooseProblem(client, 4)
	//for _, val := range GetSolvers(client) {
	//	fmt.Println(val)
	//}*/
}
