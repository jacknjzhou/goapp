package main

import (
    "net"
    "flag"
    "log"
    "fmt"
    "io/ioutil"
    "net/http"
    "html/template"
    "regexp"
    )

type Page struct {
    Title string
    Body []byte
}
var (
    addr = flag.Bool("addr", false, "Find open address and print to final-port.txt")
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename,p.Body,0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
//     m := validPath.FindStringSubmatch(r.URL.Path)
//     if m == nil {
//         http.NotFound(w, r)
//         return "", error.New("Invalid Page Title")
//     }
//     return m[2], nil
// }

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    // t, err := template.ParseFiles(tmpl + ".html")
    // if err != nil {
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    //     return
    // }
    // err = t.Execute(w, p)
    // if err != nil {
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    // }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    // title := r.URL.Path[len("/view/"):]
    // title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
    // t, _ := template.ParseFiles("view.html")
    // t.Execute(w, p)
    // fmt.Fprintf(w,"<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    // title := r.URL.Path[len("/edit/"):]
    // title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
    // t, _ := template.ParseFiles("edit.html")
    // t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    // title := r.URL.Path[len("/save/"):]
    // title, err := getTitle(w, r)
    // if err != nil {
    //     return
    // }
    body := r.FormValue("body")

    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+ title, http.StatusFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi here, I Love %s!",r.URL.Path[1:])
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    // log.Println("makeHandler")
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func main() {
    flag.Parse()
    // http.HandleFunc("/", handler)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    if *addr {
        L, err := net.Listen("tcp", "127.0.0.1:0")
        if err != nil {
            log.Fatal(err)
        }

        err = ioutil.WriteFile("final-port.txt", []byte(L.Addr().String()), 0644)
        if err != nil {
            log.Fatal(err)
        }

        s := &http.Server{}
        s.Serve(L)
        return
    }
    http.ListenAndServe(":9999", nil)
    // p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    // p1.save()
    // p2,_ := loadPage("TestPage")
    // fmt.Println(string(p2.Body))
}
