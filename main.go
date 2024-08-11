package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*
var templatesFS embed.FS

const CONTENT string = "novels"

var templates *template.Template

func init() {
	t := template.New("")

	// List of template files to parse
	templateFiles := []string{"templates/index.gohtml", "templates/novel.gohtml"}

	for _, tmpl := range templateFiles {
		content, err := templatesFS.ReadFile(tmpl)
		if err != nil {
			log.Fatalf("Failed to read template file %s: %v", tmpl, err)
		}
		t = template.Must(t.New(tmpl).Parse(string(content)))
	}

	templates = t
}

type FrontPage struct {
	Novels []Novel
}
type Novel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Chapters    []string `json:"chapters"`
}
type NovelPage struct {
	Novel      Novel
	PrefixText string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/novel/", novelHandler)
	http.HandleFunc("/novel/select/", selectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func selectHandler(writer http.ResponseWriter, request *http.Request) {
	novelTitle := strings.TrimPrefix(request.URL.Path, "/novel/select/")
	switch request.Method {
	case http.MethodGet:
		chapter, err := GetSelectedChapter(novelTitle)
		if err != nil {
			files, _ := GetFiles(CONTENT + "/" + novelTitle)
			_, _ = fmt.Fprintf(writer, "%s", files[0])
			return
		}
		_, _ = fmt.Fprintf(writer, "%s", chapter)
		return
	case http.MethodPost:
		if err := request.ParseForm(); err != nil {
			http.Error(writer, "Error parsing form data", http.StatusBadRequest)
			return
		}
		chapter := request.FormValue("chapter")
		err := SaveSelectedChapter(novelTitle, chapter)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, _ = fmt.Fprintf(writer, "Chapter %s saved", chapter)
		return
	default:
		return
	}
}

func novelHandler(writer http.ResponseWriter, request *http.Request) {
	novelTitle := strings.TrimPrefix(request.URL.Path, "/novel/")
	switch request.Method {
	case http.MethodGet:
		requestedChapter := request.URL.Query().Get("chapter")
		if requestedChapter == "" {
			writer.Header().Set("Content-Type", "text/html")
			chapters, _ := GetFiles("novels/" + novelTitle)
			novel := Novel{Title: novelTitle, Chapters: chapters, Description: ""}
			prefixText, _ := GetPrependText(novelTitle)
			novelPage := NovelPage{Novel: novel, PrefixText: prefixText}
			err := templates.ExecuteTemplate(writer, "templates/novel.gohtml", novelPage)
			if err != nil {
				log.Fatal(err)
				return
			}
		} else {
			prefixText, _ := GetPrependText(novelTitle)
			chapter, err := GetChapter(novelTitle, requestedChapter)
			if err != nil {
				fmt.Fprintf(writer, "Invalid novel title [%s] or chapter [%s]", prefixText, chapter)
				return
			}
			fmt.Fprintf(writer, "%s\n%s", prefixText, chapter)
		}
	case http.MethodPost:
		if err := request.ParseForm(); err != nil {
			http.Error(writer, "Error parsing form data", http.StatusBadRequest)
			return
		}
		prefixText := request.FormValue("prefixText")
		err := SavePrependText(novelTitle, prefixText)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, _ = fmt.Fprintf(writer, "PrefixText saved")
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	novelFolders, _ := GetSubfolders(CONTENT)
	novels := make([]Novel, len(novelFolders))
	for i, novelFolder := range novelFolders {
		chapters, _ := GetFiles(CONTENT + "/" + novelFolder)
		novels[i] = Novel{Title: novelFolder, Chapters: chapters, Description: ""}
		log.Printf("Novel %d: %s", len(chapters), novelFolder)
	}

	data := FrontPage{
		Novels: novels,
	}

	for i, novel := range novels {
		log.Printf("Novel %d: %s [%d]", i, novel.Title, len(novel.Chapters))
	}

	writer.Header().Set("Content-Type", "text/html")
	err := templates.ExecuteTemplate(writer, "templates/index.gohtml", data)
	if err != nil {
		log.Fatal(err)
		return
	}
}
