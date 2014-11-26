package main
 
import (
    "appengine"
   "appengine/datastore"
   "appengine/user"
    "net/http"
    "html/template"
	"time"
)

type Greeting struct {
   Author  string
   Content string
   Date    time.Time
}


type User struct {
   Account   string   
   Password string   
   Name string      
} 


type Error struct {
	message string
} 


var guestbookTemplate = template.Must(template.ParseFiles("tmpl/chat.htm"))
var loginTemplate = template.Must(template.ParseFiles("tmpl/loginscreen.htm"))
var applyTemplate = template.Must(template.ParseFiles("tmpl/applyscreen.htm"))
var homeTemplate = template.Must(template.ParseFiles("tmpl/index.htm"))
 
func main() {
http.HandleFunc("/", home)
http.HandleFunc("/login",login)
http.HandleFunc("/chat", chat)
http.HandleFunc("/sign", sign)
http.HandleFunc("/apply", apply)
http.ListenAndServe(":8080", nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w,nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html")
	switch r.Method{
		case "GET":
			loginTemplate.Execute(w,r.FormValue("name"))
		case "POST":

		a :=r.FormValue("account")
		p :=r.FormValue("password")
		c := appengine.NewContext(r)
		que:=datastore.NewQuery("User").Filter("Account =",a).Filter("Password =",p).Limit(1)
		result:=make([]User,0,10)
		que.GetAll(c,&result)
			if len(result)>0 {
				//greet.Author=a
   				guestbookTemplate .Execute(w,result)
			}else {
   				loginTemplate.Execute(w,r.FormValue("account"))
			}
		}
}


func guestbookKey(c appengine.Context) *datastore.Key {
        // The string "default_guestbook" here could be varied to have multiple guestbooks.
        return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}



func chat(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(c, &greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := guestbookTemplate.Execute(w, greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func sign(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := Greeting{
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/chat", http.StatusFound)
}

func apply(w http.ResponseWriter, r *http.Request) {
   switch r.Method{
   case "GET":
      applyTemplate.Execute(w, r.FormValue("name"))
   case "POST":
   c := appengine.NewContext(r)
   a :=r.FormValue("account")
que := datastore.NewQuery("User").Filter("Account =", a).Limit(1)
	result:=make([]User,0,10)
	que.GetAll(c,&result)
	

	
	if len(result)>0 {
		error:=Error {
		message: "Account Already Exists",
		}
   		applyTemplate.Execute(w,error)
	}else {

	user := User {
   		Account:r.FormValue("account"),
   		Password:r.FormValue("password"),
   		Name:r.FormValue("name"),
	}
	datastore.Put(c, datastore.NewIncompleteKey(c, "User",nil), &user)
	http.Redirect(w, r, "/", http.StatusFound)
	homeTemplate.Execute(w,nil)
	
	}
	}
}
