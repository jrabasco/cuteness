package main

import (
	"cuteness/files"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Category struct {
	Ref  string
	Name string
}

type Cuteness struct {
	Categories []Category
	Name       string
	Image      string
	Text       string
}

func contains(categories []string, category string) bool {
	for _, a := range categories {
		if a == category {
			return true
		}
	}
	return false
}

func handleCuteness(token string, tmplOK *template.Template, tmplNOK *template.Template, imgRoot string, txtRoot string, prefix string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cToken, err := r.Cookie("cuteness-token")
		if err != nil || (*cToken).Value != token {
			tmplNOK.Execute(w, prefix)
			return
		}
		categoryRefs := files.ListDirectories(imgRoot)

		categories := []Category{}
		for _, ref := range categoryRefs {
			linkref := fmt.Sprintf("/%s/show/%s", prefix, ref)
			categories = append(categories, Category{linkref, strings.Title(ref)})
		}

		basePath := path.Base(r.URL.Path)

		name := ""
		image := ""
		if basePath != "/" {
			name = basePath
			if contains(categoryRefs, name) {
				folder := fmt.Sprintf("%s/%s", imgRoot, name)
				imageName := files.RandomFile(folder)
				image = fmt.Sprintf("%s/%s/%s", imgRoot, name, imageName)
				if string(image[0]) != "/" {
					image = "/" + image
				}
			}
		}

		cute := Cuteness{
			categories,
			name,
			image,
			files.RandomFileContents(txtRoot),
		}
		tmplOK.Execute(w, cute)
	}
}

func handleAuth(password string, token string, prefix string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cPassword := path.Base(r.URL.Path)
		if cPassword == password {
			expiration := time.Now().Add(14 * 24 * time.Hour)
			cookie := http.Cookie{
				Name:     "cuteness-token",
				Value:    token,
				Expires:  expiration,
				Path:     "/" + prefix,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
		}
		http.Redirect(w, r, "/"+prefix+"/show", 303)
	}
}

func main() {
	imgRootP := flag.String("images", "imgs", "Images folder")
	txtRootP := flag.String("texts", "txts", "Texts folder")
	portP := flag.Int("port", 1234, "Port to listen on")
	flag.Parse()

	password := os.Getenv("CUTENESS_PW")
	if password == "" {
		password = "password"
	}

	token := os.Getenv("CUTENESS_TOKEN")
	if token == "" {
		token = "token"
	}

	cutenessPrefix := os.Getenv("CUTENESS_PREFIX")
	if cutenessPrefix == "" {
		cutenessPrefix = "cuteness"
	}

	tmplOK := template.Must(template.ParseFiles("cuteness.html"))
	tmplNOK := template.Must(template.ParseFiles("authenticate.html"))
	http.HandleFunc("/show/", handleCuteness(token, tmplOK, tmplNOK, *imgRootP, *txtRootP, cutenessPrefix))
	http.HandleFunc("/show", handleCuteness(token, tmplOK, tmplNOK, *imgRootP, *txtRootP, cutenessPrefix))
	http.HandleFunc("/auth/", handleAuth(password, token, cutenessPrefix))

	fmt.Printf("Listening on %d\n", *portP)
	http.ListenAndServe(fmt.Sprintf(":%d", *portP), nil)
}
