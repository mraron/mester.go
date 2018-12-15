package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"fmt"
	"net/http"
	"strconv"
	"net/url"
	"io/ioutil"
)

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
	form.Add("j_idt10", "j_idt10")
	form.Add("j_idt10:j_idt11", "megjelenit")
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

func GetTopicName(c *http.Client) string {
	req, err := http.NewRequest("GET", "https://mester.inf.elte.hu/faces/megoldasaim.xhtml", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://mester.inf.elte.hu/faces/megoldasaim.xhtml")
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
	fmt.Println(d.Text(), "!")
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

