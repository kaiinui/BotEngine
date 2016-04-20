package botengine

import (
	"github.com/go-martini/martini"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"encoding/json"
	"bytes"
	"google.golang.org/appengine/urlfetch"
	"golang.org/x/net/context"
	"github.com/codegangsta/martini-contrib/render"
)

type BotServer struct {
	fbPageAccessToken string
}

type BotThinkFunction func(Context, Message, Action, User) ThinkResponse

func NewServer(fbPageAccessToken string) BotServer {
	return BotServer{
		fbPageAccessToken: fbPageAccessToken,
	}
}

type ThinkRequest struct {
	Context Context
	Message Message
	Action Action
	User User
}

type Context struct {
	State string
	Req *http.Request
}

type Message struct {
	Text string
}

type Action struct {
	Kind string
	Attributes map[string]string
}

type User struct {
	ID int64
}

type ThinkResponse struct {
	Message Message
	Context Context
}

func (s BotServer) Run(fn BotThinkFunction) {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/webhook", verifyWebhook)
	m.Post("/webhook", receiveMessage)
	m.Post("/message.resolve", func(req *http.Request) {
		s.resolveMessage(req)
	})
	m.Post("/message.think", func(req *http.Request, ren render.Render) {
		passToThinkMessageFunc(req, ren, fn)
	})

	http.Handle("/", m)
}

func passToThinkMessageFunc(req *http.Request, ren render.Render, fn BotThinkFunction) {
	ctx := appengine.NewContext(req)

	body := req.Body
	defer body.Close()
	var r ThinkRequest
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		panic(err)
	}

	log.Infof(ctx, "Incoming Message: %s", r)

	r.Context.Req = req
	resp := fn(r.Context, r.Message, r.Action, r.User)

	log.Infof(ctx, "Think Result: %s", resp)

	ren.JSON(200, resp)
}

func verifyWebhook(req *http.Request) string {
	// webhook endpoint verification
	// https://developers.facebook.com/docs/messenger-platform/quickstart
	if req.URL.Query().Get("hub.verify_token") == "botless" {
		return req.URL.Query().Get("hub.challenge")
	}

	return "failure"
}

func receiveMessage(req *http.Request) string {
	ctx := appengine.NewContext(req)

	body := req.Body
	defer body.Close()
	err, payload := DecodeWebhookPayload(body)
	if err != nil {
		log.Warningf(ctx, "Error while parsing payload: %s", err)
		panic(err)
	}

	log.Infof(ctx, "Payload: %s", payload)

	for _, entry := range payload.Entry {
		for _, messaging := range entry.Messaging {
			var pld []byte
			buf := bytes.NewBuffer(pld)
			err := json.NewEncoder(buf).Encode(messaging)
			if err != nil {
				log.Warningf(ctx, "Error while parsing messaging: %s", err)
				panic(err)
			}

			task := &taskqueue.Task{
				Path: "/message.resolve",
				Method: "POST",
				Payload: buf.Bytes(),
			}
			taskqueue.Add(ctx, task, "default")
		}
	}

	return "ok"
}

func (s BotServer) resolveMessage(req *http.Request) {
	ctx := appengine.NewContext(req)

	var msg FBWebhookMessaging
	if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
		log.Warningf(ctx, "Error while parsing message: %s", err)
		panic(err)
	}

	if thinkReq := makeThinkRequestWithFBMessaging(ctx, msg); thinkReq != nil {
		// 1. /message.think に処理を投げる
		var reqBodyBytes []byte
		reqBodyBuf := bytes.NewBuffer(reqBodyBytes)
		if err := json.NewEncoder(reqBodyBuf).Encode(thinkReq); err != nil {
			panic(err)
		}

		cli := urlfetch.Client(ctx)
		resp, err := cli.Post(makeThisUriToPath(ctx, "/message.think"), "application/json", reqBodyBuf) // TODO: Channel化
		if err != nil {
			log.Warningf(ctx, "Error while requesting /message.think: %s", err)
			panic(err)
		}

		// 2. /message.think の処理結果を FB bot に転送
		respBody := resp.Body
		defer respBody.Close()
		var thinkResp ThinkResponse
		if err := json.NewDecoder(respBody).Decode(&thinkResp); err != nil {
			panic(err)
		}

		if thinkResp.Message.Text != "" {
			fbReq := NewFBMessageWithTextRequest(msg.Sender.Id, thinkResp.Message.Text)
			DispatchFBSendMessageRequest(cli, s.fbPageAccessToken, fbReq)
		}
	}
}

func makeThinkRequestWithFBMessaging(ctx context.Context, msg FBWebhookMessaging) *ThinkRequest {
	if msg.Postback != nil {
		return &ThinkRequest{
			Action: Action{Kind: msg.Postback.Payload},
			User: User{ID: msg.Sender.Id},
		}
	} else if msg.Optin != nil {
		return &ThinkRequest{
			Action: Action{Kind: "$referer", Attributes: map[string]string {
				"$referer": msg.Optin.Ref,
			}},
			User: User{ID: msg.Sender.Id},
		}
	} else if msg.Delivery != nil {
		// No need
	} else if msg.Message != nil {
		return &ThinkRequest{
			Message: Message{
				Text: msg.Message.Text,
			},
			User: User{ID: msg.Sender.Id},
		}
	}

	return nil
}

func makeThisUriToPath(ctx context.Context, path string) string {
	if appengine.IsDevAppServer() {
		return "http://localhost:8080" + path
	} else {
		return "https://" + appengine.DefaultVersionHostname(ctx) + path
	}
}
