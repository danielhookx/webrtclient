package signal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Hub struct {
	Name string
	Room string
}

func Cors(w http.ResponseWriter, r *http.Request) error{
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Headers", "*") //Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,FZM-APP-ID
	h.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	h.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	h.Set("Access-Control-Allow-Credentials", "true")

	// 放行所有OPTIONS方法，因为有的模板是要请求两次的
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return errors.New("OPTIONS request")
	}
	return nil
}

func HTTPPubSubServer(port int) (chan *Hub, chan *Hub) {
	pubChan := make(chan *Hub)
	subChan := make(chan *Hub)
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		err := Cors(w,r)
		if err != nil {
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var hub Hub
		err = json.Unmarshal(body, &hub)
		if err != nil {
			log.Printf("publish body err: %v", err)
			return
		}
		pubChan <- &hub
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		err := Cors(w,r)
		if err != nil {
			return
		}
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var hub Hub
		log.Print(string(body))
		err = json.Unmarshal(body, &hub)
		if err != nil {
			log.Printf("subscribe body err: %v", err)
			return
		}
		subChan <- &hub
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return pubChan, subChan
}
