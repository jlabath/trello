package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//CommentPost posts a comment on the trello card specified by cardId
func CommentPost(a *Auth, client *http.Client, cardId string, comment string) error {
	postData := url.Values{}
	postData.Set("key", a.Key)
	postData.Set("token", a.Token)
	postData.Set("text", comment)
	apiUrl := fmt.Sprintf("%s/actions/comments", cardUrl(cardId))
	resp, err := client.PostForm(apiUrl, postData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return apiError(resp)
	}
	return nil
}

//CommentList returns a list of comments from a Card.
//This has a max limit or results returned,
//by whatever the cap set by trello (50 as of today).
func CommentList(a *Auth, client *http.Client, cardId string) ([]Comment, error) {
	var comments []Comment
	args := url.Values{}
	args.Set("key", a.Key)
	args.Set("token", a.Token)
	args.Set("filter", "commentCard")
	args.Set("fields", "data")
	args.Set("memberCreator_fields", "fullName,username")
	apiUrl := fmt.Sprintf("%s/actions?%s", cardUrl(cardId), args.Encode())
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, apiError(resp)
	}
	//json decode
	dec := json.NewDecoder(resp.Body)
	for {
		var c []trelloAction
		if derr := dec.Decode(&c); derr == io.EOF {
			break
		} else if derr != nil {
			return nil, derr
		}
		for _, v := range c {
			comments = append(comments, v)
		}
	}
	return comments, nil

}

//CardGet returns a Card interface representing the trello card
func CardGet(a *Auth, client *http.Client, cardId string) (Card, error) {
	var card Card
	args := url.Values{}
	args.Set("key", a.Key)
	args.Set("token", a.Token)
	apiUrl := fmt.Sprintf("%s?%s", cardUrl(cardId), args.Encode())
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, apiError(resp)
	}
	//json decode
	dec := json.NewDecoder(resp.Body)
	for {
		var c trelloCard
		if derr := dec.Decode(&c); derr == io.EOF {
			break
		} else if derr != nil {
			return nil, derr
		}
		card = c
		//fmt.Printf("%#v\n", c)
	}
	return card, nil
}

//CardListPut moves the Card to the desired List
func CardListPut(a *Auth, client *http.Client, cardId string, listId string) (Card, error) {
	var card Card
	args := url.Values{}
	args.Set("key", a.Key)
	args.Set("token", a.Token)
	args.Set("value", listId)
	apiUrl := fmt.Sprintf("%s/idList", cardUrl(cardId))
	putReq, err := http.NewRequest("PUT", apiUrl, strings.NewReader(args.Encode()))
	if err != nil {
		return nil, err
	}
	putReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(putReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, apiError(resp)
	}
	//json decode
	dec := json.NewDecoder(resp.Body)
	for {
		var c trelloCard
		if derr := dec.Decode(&c); derr == io.EOF {
			break
		} else if derr != nil {
			return nil, derr
		}
		card = c
		//fmt.Printf("%#v\n", c)
	}
	return card, nil
}

//BoardGet returns a Board interface representing the trello Board
func BoardGet(a *Auth, client *http.Client, boardId string) (Board, error) {
	var board Board
	args := url.Values{}
	args.Set("key", a.Key)
	args.Set("token", a.Token)
	args.Set("fields", "name")
	apiUrl := fmt.Sprintf("%s?%s", boardUrl(boardId), args.Encode())
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, apiError(resp)
	}
	//json decode
	dec := json.NewDecoder(resp.Body)
	for {
		var c trelloBoard
		if derr := dec.Decode(&c); derr == io.EOF {
			break
		} else if derr != nil {
			return nil, derr
		}
		board = c
		//fmt.Printf("%#v\n", c)
	}
	return board, nil
}

//BoardListsGet returns a slice of List interfaces representing open lists on the board
func BoardListsGet(a *Auth, client *http.Client, boardId string) ([]List, error) {
	var lists []List
	args := url.Values{}
	args.Set("key", a.Key)
	args.Set("token", a.Token)
	apiUrl := fmt.Sprintf("%s/lists?%s", boardUrl(boardId), args.Encode())
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, apiError(resp)
	}
	//json decode
	dec := json.NewDecoder(resp.Body)
	for {
		var c []trelloList
		if derr := dec.Decode(&c); derr == io.EOF {
			break
		} else if derr != nil {
			return nil, derr
		}
		//fmt.Printf("%#v\n", c[0])
		for _, x := range c {
			lists = append(lists, x)
		}
	}
	return lists, nil
}

//apiError returns error containing HTTP status and response body
//if there was failure reading response IO error could be returned as well
func apiError(resp *http.Response) error {
	buf, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		return ioerr
	}
	return fmt.Errorf("%s:%s", resp.Status, string(buf))
}

//cardUrl returns the REST API url of the trello card
func cardUrl(cardId string) string {
	url := fmt.Sprintf("%s/cards/%s", apiUrl(), cardId)
	return url
}

func boardUrl(id string) string {
	url := fmt.Sprintf("%s/boards/%s", apiUrl(), id)
	return url
}

func apiUrl() string {
	return "https://api.trello.com/1"
}
