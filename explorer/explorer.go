package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Gunyoung-Kim/blockchain/blockchain"
)

const (
	port        string = ":3000"
	templateDir string = "explorer/templates/"
)

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{PageTitle: "Home", Blocks: blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockdata")
		fmt.Println(data)
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

// Start explorer
func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "fragments/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	log.Fatal(http.ListenAndServe(port, nil))
}
