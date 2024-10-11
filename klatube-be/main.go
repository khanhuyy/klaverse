// package main

// import (
// 	_ "embed"
// 	"encoding/json"
// 	"fmt"
// 	"html/template"
// 	"net/http"
// 	"os"
// )

// type Todo struct {
// 	Title string
// 	Done  bool
// }

// type TodoPageData struct {
// 	PageTitle string
// 	Todos     []Todo
// }

// type User struct {
//     Firstname string `json:"firstname"`
//     Lastname  string `json:"lastname"`
//     Age       int    `json:"age"`
// }

// //go:embed childe_tighnari_4-4.mp4
// var s1 string

// //go:embed navia_hutao_4-4.mp4
// var s2 string

// var videos = map[int64]string{
// 	1: s1,
// 	2: s2,
// }

// func main() {
// 	tmpl := template.Must(template.ParseFiles("layout.html"))
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		data := TodoPageData{
// 			PageTitle: "My TODO list",
// 			Todos: []Todo{
// 				{Title: "Task 1", Done: false},
// 				{Title: "Task 2", Done: true},
// 				{Title: "Task 3", Done: true},
// 			},
// 		}
// 		tmpl.Execute(w, data)
// 	})
// 	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
//         peter := User{
//             Firstname: "John",
//             Lastname:  "Doe",
//             Age:       25,
//         }

//         json.NewEncoder(w).Encode(peter)
//     })
// 	http.HandleFunc("/videos/1", func(w http.ResponseWriter, r *http.Request) {
// 		filePath, err := os.ReadFile(s1)
// 		if err != nil {
// 			fmt.Print(s1)
// 			fmt.Print(err)
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/octet-stream")
// 		w.Write(filePath)
// 		return
// 	})
// 	http.ListenAndServe(":8000", nil)
// }

// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"net/http"
// 	"os"
// 	"time"
// )

// func run() error {
// 	f, err := os.Open("classic.mp3")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	d, err := mp3.NewDecoder(f)
// 	if err != nil {
// 		return err
// 	}

// 	c, ready, err := oto.NewContext(d.SampleRate(), 2, 2)
// 	if err != nil {
// 		return err
// 	}
// 	<-ready

// 	p := c.NewPlayer(d)
// 	defer p.Close()
// 	p.Play()

// 	fmt.Printf("Length: %d[bytes]\n", d.Length())
// 	for {
// 		time.Sleep(time.Second)
// 		if !p.IsPlaying() {
// 			break
// 		}
// 	}

// 	return nil
// }

// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List([]string{part})
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func initBackend() {
	tmpl := template.Must(template.ParseFiles("layout.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}
		tmpl.Execute(w, data)
	})
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var res interface{}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(res)
		http.Redirect(w, r, "/", 200)
	})
	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		b, err := ioutil.ReadFile("client_secret.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		// If modifying these scopes, delete your previously saved credentials
		// at ~/.credentials/youtube-go-quickstart.json
		config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
		client := getClient(ctx, config)
		fmt.Println(client)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}

		service, err := youtube.New(client)

		handleError(err, "Error creating YouTube client")

		channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
	})
	http.ListenAndServe(":3000", nil)
}

func main() {
	initBackend()
}
func ToWords(n int64) string {
	exceptionOne := map[int]string{
		1: "mốt",
	}

	exception := map[int64]string{
		0: "lẻ",   // ngoai le hang chuc x
		1: "mười", // ngoai le hang chuc x
		4: "tư",   // ngoai le hang don vi x
		5: "năm",  // ngoai le hang don vi x
	}

	group := map[int64]string{
		1: "mươi",
		2: "trăm",
		3: "nghìn",
		6: "triệu",
		9: "tỷ",
	}

	num := map[string]string{
		"0": "không",
		"1": "một",
		"2": "hai",
		"3": "ba",
		"4": "bốn",
		"5": "năm",
		"6": "sáu",
		"7": "bảy",
		"8": "tám",
		"9": "chín",
	}
	s := fmt.Sprintf("%d", n)
	var lenString = int64(len(s))
	check9 := lenString - (lenString % 9)
	check3 := lenString - (lenString % 3) // position multi tripple (kiem tra bo 3)
	offset := lenString - check3 + 1      // phan le
	oneFlag := false                      // flag to check exceptionOne
	if lenString < 2 {
		result := num[s]
		charResult := []byte(result)
		charResult[0] = byte(unicode.ToUpper(rune(charResult[0])))
		result = string(charResult)
		result = fmt.Sprintf(strings.TrimRight(result, " "))
		result += " đồng"
		return result
	} else {
		result := ""
		var flag int64 // kiem tra so khong
		for pos, char := range s {
			if offset == 1 {
				offset = 3
				check3 -= 3
			} else {
				offset -= 1
			}
			unit := len(s) - pos - 1
			triple := unit % 3
			switch triple {
			case 1:
				if fmt.Sprintf("%c", char) == "0" {
					flag += 1
				} else if fmt.Sprintf("%c", char) == "1" {
					if flag%3 == 1 {
						result += num["0"] + " " + group[2] + " " + exception[1] + " "
					} else {
						result += exception[1] + " "
					}
					flag = 0
				} else {
					if flag%3 == 1 {
						result += num["0"] + " " + group[2] + " " + num[fmt.Sprintf("%c", char)] + " " + group[1] + " "
					} else {
						result += num[fmt.Sprintf("%c", char)] + " " + group[1] + " "
					}
					oneFlag = true
					flag = 0
				}
			case 2:
				if fmt.Sprintf("%c", char) == "0" {
					flag += 1
				} else {
					flag = 0
					result += num[fmt.Sprintf("%c", char)] + " " + group[2] + " "
				}
			default:
				if fmt.Sprintf("%c", char) == "0" {
					flag += 1
				} else if fmt.Sprintf("%c", char) == "1" {
					if flag > 1 {
						result += exception[0] + " " + num["1"] + " "
						oneFlag = false
					} else {
						if oneFlag {
							result += exceptionOne[1] + " "
							oneFlag = true
						} else {
							result += num["1"] + " "
						}
					}
					flag = 0
				} else if fmt.Sprintf("%c", char) == "4" {
					if flag > 1 {
						result += exception[0] + " " + exception[4] + " "
					} else {
						if pos != 0 {
							result += exception[4] + " "
						} else {
							result += num["4"] + " "
						}
					}
					flag = 0
				} else if fmt.Sprintf("%c", char) == "5" {
					if flag > 1 {
						result += exception[0] + " " + num["5"] + " "
					} else {
						result += exception[5] + " "
					}
					flag = 0
				} else {
					if flag > 1 {
						result += exception[0] + " " + num[fmt.Sprintf("%c", char)] + " "
					} else {
						result += num[fmt.Sprintf("%c", char)] + " "
					}
					flag = 0
				}
				if check3%9 == 0 && check3 > 0 {
					check9 -= 9
					result += group[9] + " "
				} else {
					if flag < 1 {
						result += group[check3%9] + " "
						flag = 0
					} else if flag/3 < 1 {
						result += group[check3%9] + " "
						flag = 0
					}
				}
			}
		}
		charResult := []byte(result)
		charResult[0] = byte(unicode.ToUpper(rune(charResult[0])))
		result = string(charResult)
		result = fmt.Sprintf(strings.TrimRight(result, " "))
		result += " đồng"
		return result
	}
}