package main

import (
	"github.com/unrolled/render"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"os"
	_ "html/template"
	"encoding/json"
	"sort"
	"html/template"
	_"math"

	"math"
)


const dataDir = "../data.json"

type Solution struct {
	Statement string
	Topic     string
	Problem   string
	Name      string
	Point     int
}

var Solutions []Solution

type Id struct {
	Topic string
	Problem string
}

var ProblemPage map[Id][]Solution
var UserPage map[string][]Solution

type RankRow struct {
	Name string
	Link string
	PointSum float64
}
type RankList []RankRow

func (r RankList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RankList) Less(i, j int) bool {
	return r[i].PointSum>r[j].PointSum
}

func (r RankList) Len() int {
	return len(r)
}

var TopicRankList map[string]RankList

var UserList map[string]string
var TopicList map[string]string
var ProblemList map[string]Id

var BigRanking RankList
var ProblemRanking RankList;

var SolversCount map[string]int
var PointSum map[string]int

const DynamicRating = true

var RatingFunction func(Solution) float64

func CalculateSumRating(val Solution) float64 {
	return float64(val.Point)
}

func CalculateDynamicRating(val Solution) float64 {
	pointsum := float64(0)
	solvercount := 0
	for _, val2 := range ProblemPage[Id{val.Topic,val.Problem}] {
		if val2.Point>=val.Point {
			pointsum += float64(val2.Point)
			solvercount ++
		}
	}

	if float64(PointSum[val.Problem])/float64(SolversCount[val.Problem])>float64(val.Point) {
		return 0
	}

	return math.Sqrt(float64(pointsum*float64(val.Point))/float64(solvercount))
}

func init() {
	if DynamicRating {
		RatingFunction = CalculateDynamicRating
	}else {
		RatingFunction = CalculateSumRating
	}

	f, err := os.Open(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&Solutions)
	if err != nil {
		log.Fatal(err)
	}

	ProblemPage = make(map[Id][]Solution)
	for _, val := range Solutions {
		ProblemPage[Id{val.Topic, val.Problem}] = make([]Solution, 0)
	}
	for _, val := range Solutions {
		ProblemPage[Id{val.Topic, val.Problem}] = append(ProblemPage[Id{val.Topic, val.Problem}], val)
	}

	UserPage = make(map[string][]Solution)
	for _, val := range Solutions {
		UserPage[val.Name] = make([]Solution, 0)
	}
	for _, val := range Solutions {
		UserPage[val.Name] = append(UserPage[val.Name], val)
	}

	SolversCount = make(map[string]int)
	PointSum = make(map[string]int)

	for _, val := range Solutions {
		SolversCount[val.Problem] ++
		PointSum[val.Problem] += val.Point
	}

	TopicRankList = make(map[string]RankList)

	for _, val := range Solutions {
		TopicRankList[val.Topic] = make(RankList, 0)
	}

	for _, val := range Solutions {
		found := false
		for ind, val2 := range TopicRankList[val.Topic] {
			if val2.Name == val.Name {
				found = true
				TopicRankList[val.Topic][ind].PointSum += RatingFunction(val)
			}
		}
		if !found {
			TopicRankList[val.Topic] = append(TopicRankList[val.Topic], RankRow{val.Name,"/user/"+val.Name+"/", RatingFunction(val)})
		}

		found = false
		for ind, val2 := range BigRanking {
			if val2.Name == val.Name {
				found = true
				BigRanking[ind].PointSum += RatingFunction(val)
			}
		}
		if !found {
			BigRanking = append(BigRanking, RankRow{val.Name,"/user/"+val.Name+"/", RatingFunction(val)})
		}

		found = false
		for ind, val2 := range ProblemRanking {
			if val2.Name == val.Topic + " / " + val.Problem {
				found = true
				ProblemRanking[ind].PointSum ++
			}
		}
		if !found {
			ProblemRanking = append(ProblemRanking, RankRow{val.Topic + " / " + val.Problem, "/problem/"+val.Topic+"/"+val.Problem+"/", 1})
		}
	}

	for ind, _ := range TopicRankList {
		sort.Sort(TopicRankList[ind])
	}

	sort.Sort(BigRanking)
	sort.Sort(ProblemRanking)

	UserList = make(map[string]string)
	ProblemList = make(map[string]Id)
	TopicList = make(map[string]string)
	for _, val := range Solutions {
		UserList[val.Name] = val.Name
		ProblemList[val.Problem] = Id{val.Topic, val.Problem}
		TopicList[val.Topic] = val.Topic
	}
}



func main() {

	router := mux.NewRouter()
	renderer := render.New(render.Options{
		Layout: "layout",
		Extensions: []string{".tmpl", ".html"},
		Funcs: []template.FuncMap{
			template.FuncMap{"add": func(a,b int) int {
				return a+b
			},},
		},
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "index", struct {
			UserList map[string]string
			ProblemList map[string]Id
			TopicList map[string]string
		}{UserList, ProblemList, TopicList})
	})
	router.HandleFunc("/problem/{topic}/{problem}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		renderer.HTML(w, http.StatusOK, "problem", ProblemPage[Id{vars["topic"], vars["problem"]}])
	})

	router.HandleFunc("/problem/ranking/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "ranking", ProblemRanking)
	})

	router.HandleFunc("/user/{name}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		renderer.HTML(w, http.StatusOK, "user", UserPage[vars["name"]])
	})

	router.HandleFunc("/topic/{topic}/ranking/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		renderer.HTML(w, http.StatusOK, "ranking", TopicRankList[vars["topic"]])
	})

	router.HandleFunc("/ranking/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "ranking", BigRanking)
	})

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}