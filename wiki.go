package main

import (
	"fmt"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// html 파일 미리 parse 해 놓음.
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// invalid path 막기 위한 validation
var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	/* : xxx.html 파일 캐싱 (미리 parse)
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		for _, e := range m {
			fmt.Printf("m's memeber : %s\n", e)
		}

		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	/* : makeHandler
	title, err := getTitle(w, r)
	if err != nil {
		return
	}*/
	/* url validation 추가
	title := r.URL.Path[len("/view/"):]*/
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}

	/* html hard-coded : using template
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)*/
	/* html template parse : common func
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)*/
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	/* url validation 추가
	title := r.URL.Path[len("/edit/"):]*/
	/* : makeHandler
	title, err := getTitle(w, r)
	if err != nil {
		return
	}*/

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	/* html hard-coded : using template
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>", p.Title, p.Title, p.Body)*/
	/* html template parse : common func
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)*/
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	/* : makeHandler
	title, err := getTitle(w, r)
	if err != nil {
		return
	}*/

	/* url validation 추가
	title := r.URL.Path[len("/save/"):]*/
	body := r.FormValue("body")
	p := &Page{
		Title: title,
		Body:  []byte(body),
	}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	/* : makeHandler
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)*/

	log.Fatal(http.ListenAndServe(":8080", nil))

	/* 연습 페이지 : crud
	p1 := &Page{Title: "TestPage", Body:  []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))*/
}

