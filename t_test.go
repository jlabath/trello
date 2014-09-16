package trello

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

var testCfg struct {
	Key    string //trello api key
	Token  string //trello auth token
	ToList string //name of the board to move card to
	Card   string //id or last part of short link
}

func getTestAuth() *Auth {
	return &Auth{
		Key:   testCfg.Key,
		Token: testCfg.Token}
}

func getTestCardId() string {
	return testCfg.Card
}

func getTestToList() string {
	return testCfg.ToList
}

func getTestClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	return client
}

func TestComment1(t *testing.T) {
	cl := getTestClient()
	err := CommentPost(getTestAuth(), cl, getTestCardId(), "test 1\nhello world")
	ok(t, err)
	comments, err := CommentList(getTestAuth(), cl, getTestCardId())
	ok(t, err)
	found := false
	for _, v := range comments {
		if v.Text() == "test 1\nhello world" {
			found = true
			break
		}
	}
	equals(t, found, true)
}

func TestCardGet(t *testing.T) {
	cl := getTestClient()
	card, err := CardGet(getTestAuth(), cl, getTestCardId())
	ok(t, err)
	equals(t, "Test", card.Name())
	equals(t, "FOOBAR", card.Desc())
	assert(t, "" != card.IdList(), "List ID should not be empty", card)
	assert(t, "" != card.IdBoard(), "Board ID should not be empty", card)
	assert(t, "" != card.Id(), "ID should not be empty", card)
	assert(t, "" != card.ShortUrl(), "ShortUrl should not be empty", card)
	assert(t, "" != card.DateLastActivity(), "DateLastActivity should not be empty", card.DateLastActivity())
}

func TestMoveCard(t *testing.T) {
	cl := getTestClient()
	a := getTestAuth()
	card, err := CardGet(a, cl, getTestCardId())
	ok(t, err)
	equals(t, "Test", card.Name())
	equals(t, "FOOBAR", card.Desc())
	assert(t, "" != card.IdBoard(), "Board ID should not be empty", card)
	board, err := BoardGet(a, cl, card.IdBoard())
	ok(t, err)
	assert(t, "" != board.Id(), "ID should not be empty", board.Id())
	assert(t, "" != board.Name(), "Name should not be empty", board.Name())
	lists, err := BoardListsGet(a, cl, board.Id())
	ok(t, err)
	var destList List
	for _, l := range lists {
		assert(t, "" != l.Id(), "ID should not be empty", l.Id())
		assert(t, "" != l.Name(), "Name should not be empty", l.Name())
		if l.Name() == getTestToList() {
			destList = l
		}
	}
	updCard, err := CardListPut(a, cl, card.Id(), destList.Id())
	ok(t, err)
	equals(t, updCard.Id(), card.Id())
	assert(t, updCard.IdList() != card.IdList(), "the two lists should be different old:%s new:%s", card.IdList(), updCard.IdList())
	//move it back
	_, err = CardListPut(a, cl, card.Id(), card.IdList())
	ok(t, err)
}

func init() {
	f, err := os.Open("test_config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for {
		if derr := dec.Decode(&testCfg); derr == io.EOF {
			break
		} else if derr != nil {
			log.Fatal(derr)
			break
		}
		//fmt.Printf("%#v\n", testCfg)
	}
}
