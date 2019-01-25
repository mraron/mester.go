package main

import (
	"time"
	"github.com/unrolled/render"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"log"
	"github.com/fsnotify/fsnotify"
	"os"
	_ "html/template"
	"encoding/json"
	"sort"
	"html/template"
	"sync"

	"math"
)


const dataDir = "../data.json"

type PointHistoryElem struct {
    Time time.Time
    Point int
}

type Solution struct {
	Statement string
	Topic     string
	Problem   string
	Name      string
	Point     int
	PointHistory []PointHistoryElem
}

type Submission struct {
	Topic string
	Problem string
	Name string
	Point int
	Time time.Time
}

var Solutions []Solution

type Id struct {
	Topic string
	Problem string
}

var ProblemPage map[Id]*SolutionList
var UserPage map[string]*SolutionList

type Comparison struct {
	Topic string
	Problem string
	
	Tried1 bool
	Point1 int
		
	Tried2 bool
	Point2 int
	
	Verdict int
}

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

var Submissions []Submission

const DynamicRating = true

var RatingFunction func(Solution) float64

func CalculateSumRating(val Solution) float64 {
	return float64(val.Point)
}

func CalculateDynamicRating(val Solution) float64 {
	pointsum := float64(0)
	solvercount := 0
	for _, val2 := range ProblemPage[Id{val.Topic,val.Problem}].Solutions {
		if val2.Point>=val.Point {
			pointsum += float64(val2.Point)
			solvercount ++
		}
	}
	
	hossz := len(ProblemPage[Id{val.Topic, val.Problem}].Solutions)
	sign := float64(1.0)
	if ProblemPage[Id{val.Topic, val.Problem}].Solutions[hossz/2].Point>val.Point {
		sign = float64(-1)
	}

	return sign*math.Sqrt(float64(pointsum*float64(val.Point))/float64(solvercount))
}

type SolutionList struct {
	Solutions []Solution
	RelativeDistribution []float64
	MaximumElement float64
	Statement string
	Submissions []Submission
}


func LoadAndParseData() {
	log.Print("begin parsing: ", dataDir)
	dataLock.Lock()
	f, err := os.Open(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&Solutions)
	if err != nil {
		log.Fatal(err)
	}

	ProblemPage = make(map[Id]*SolutionList)
	for _, val := range Solutions {
		ProblemPage[Id{val.Topic, val.Problem}] = &SolutionList{}
		ProblemPage[Id{val.Topic, val.Problem}].RelativeDistribution = make([]float64, 0)
		ProblemPage[Id{val.Topic, val.Problem}].Solutions = make([]Solution, 0)
	}
	for _, val := range Solutions {
		ProblemPage[Id{val.Topic, val.Problem}].Statement = val.Statement
		if ProblemPage[Id{val.Topic, val.Problem}].MaximumElement < float64(val.Point) + 1 {
			ProblemPage[Id{val.Topic, val.Problem}].MaximumElement = float64(val.Point) + 1
		}
	}
	
	for _, val := range Solutions {
		ProblemPage[Id{val.Topic, val.Problem}].RelativeDistribution = append(ProblemPage[Id{val.Topic, val.Problem}].RelativeDistribution, float64(val.Point)/ProblemPage[Id{val.Topic, val.Problem}].MaximumElement)
		ProblemPage[Id{val.Topic, val.Problem}].Solutions = append(ProblemPage[Id{val.Topic, val.Problem}].Solutions, val) 	
	}

	UserPage = make(map[string]*SolutionList)
	for _, val := range Solutions {
		UserPage[val.Name] = &SolutionList{}
		UserPage[val.Name].RelativeDistribution = make([]float64, 0)
		UserPage[val.Name].Solutions = make([]Solution, 0)
		UserPage[val.Name].Submissions = make([]Submission, 0)
	}
	for _, val := range Solutions {
		if UserPage[val.Name].MaximumElement < float64(val.Point) + 1 {
			UserPage[val.Name].MaximumElement = float64(val.Point) + 1
		}
	}

	for _, val := range Solutions {
		UserPage[val.Name].RelativeDistribution = append(UserPage[val.Name].RelativeDistribution, float64(val.Point)/ProblemPage[Id{val.Topic, val.Problem}].MaximumElement)
		UserPage[val.Name].Solutions = append(UserPage[val.Name].Solutions, val)
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

	BigRanking = make(RankList, 0)

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
	
	Submissions = make([]Submission, 0)
	for _, sol := range Solutions {
		for _, helem := range sol.PointHistory {
			Submissions = append(Submissions, Submission{sol.Topic, sol.Problem, sol.Name, helem.Point, helem.Time})
			UserPage[sol.Name].Submissions = append(UserPage[sol.Name].Submissions, Submission{sol.Topic, sol.Problem, sol.Name, helem.Point, helem.Time})
		}
	}
	
	sort.Slice(Submissions, func(i, j int) bool {
		return Submissions[j].Time.Before(Submissions[i].Time)
	})
	
	for key, _ := range UserPage {
		sort.Slice(UserPage[key].Submissions, func(i, j int) bool {
			return UserPage[key].Submissions[j].Time.Before(UserPage[key].Submissions[i].Time)
		})
	}
	
	dataLock.Unlock()
	log.Print("parsed: ", dataDir)
}

func init() {
	if DynamicRating {
		RatingFunction = CalculateDynamicRating
	}else {
		RatingFunction = CalculateSumRating
	}
	
	LoadAndParseData()	
}

var dataLock sync.Mutex

func LockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("lock")
		dataLock.Lock()
		next.ServeHTTP(w, r)
	})
}

func UnlockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		println("unlock")
		dataLock.Unlock()
	})
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				time.Sleep(5*time.Second)
				LoadAndParseData()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	
	err = watcher.Add(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	
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
	router.HandleFunc("/problem/", func(w http.ResponseWriter, r *http.Request) {
		topics := r.URL.Query()["topic"]
		problems := r.URL.Query()["problem"]
		renderer.HTML(w, http.StatusOK, "problem", ProblemPage[Id{topics[0], problems[0]}])
	})

	router.HandleFunc("/problem/ranking/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "ranking", ProblemRanking)
	})

	router.HandleFunc("/user/{name}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		renderer.HTML(w, http.StatusOK, "user", UserPage[vars["name"]])
	})

	router.HandleFunc("/topic_ranking/", func(w http.ResponseWriter, r *http.Request) {
		topics := r.URL.Query()["topic"]
		renderer.HTML(w, http.StatusOK, "ranking", TopicRankList[topics[0]])
	})

	router.HandleFunc("/ranking/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "ranking", BigRanking)
	})
	
	router.HandleFunc("/compare/{you}/{other}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clist := make([]*Comparison, 0)
		
		for _, sol := range UserPage[vars["you"]].Solutions {
			clist = append(clist, &Comparison{sol.Topic, sol.Problem, true, sol.Point, false, -1, -1})
		}
		
		for _, sol := range UserPage[vars["other"]].Solutions {
			found := false
			for id, comp := range clist {
				if comp.Topic == sol.Topic && comp.Problem == sol.Problem {
					found = true
					clist[id].Point2 = sol.Point
					clist[id].Tried2 = true
				}
			}
			
			if !found {
				clist = append(clist, &Comparison{sol.Topic, sol.Problem, false, -1, true, sol.Point, -1})
			}
		}
		
				
		for i, _ := range clist {
			if clist[i].Tried1 && (!clist[i].Tried2 || (clist[i].Tried2 && clist[i].Point1 > clist[i].Point2)) {
				clist[i].Verdict = -1
			} else if clist[i].Tried1 && clist[i].Tried2 && clist[i].Point1 == clist[i].Point2 {
				clist[i].Verdict = 0
			}else {
				clist[i].Verdict = 1
			}
		}
		
		
		sort.SliceStable(clist, func(i, j int) bool {
			if clist[i].Verdict!=clist[j].Verdict {
				return clist[i].Verdict>clist[j].Verdict
			}
			
			vali := 0
			if clist[i].Tried1 {
				vali += 10
			}
			
			if !clist[i].Tried2 {
				vali += 1
			}
			
			valj := 0
			if clist[j].Tried1 {
				valj += 10
			}
			
			if !clist[j].Tried2 {
				valj += 1
			}
			
			return vali<valj
		})

		renderer.HTML(w, http.StatusOK, "compare", clist)
	})
	
	router.HandleFunc("/submissions/", func(w http.ResponseWriter, r *http.Request) {
		renderer.HTML(w, http.StatusOK, "submissions", Submissions)
	});

	http.Handle("/", handlers.LoggingHandler(os.Stdout, UnlockMiddleware(LockMiddleware(router))))
	http.Handle("/statements/", http.StripPrefix("/statements/", http.FileServer(http.Dir("../statements/"))))
	
	http.ListenAndServe(":8080", nil)
}
