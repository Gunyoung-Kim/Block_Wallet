package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Gunyoung-Kim/blockchain/blockchain"
)

var port string

const templateDir string = "explorer/templates/"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{PageTitle: "Home", Blocks: blockchain.Blocks(blockchain.BlockChain())}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.BlockChain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

// Start explorer
func Start(portNum int) {
	handler := http.NewServeMux()
	port = fmt.Sprintf(":%d", portNum)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "fragments/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("EXPLORER Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
