package botengine

import (
	"os"
	"msg_resolver"
"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"strconv"
	"math/rand"
	"strings"
	"google.golang.org/appengine"
)

func init() {
	s := botengine.NewServer(os.Getenv("FB_PAGE_ACCESS_TOKEN"))

	s.Run(onThink)
}

func onThink(ctx botengine.Context, msg botengine.Message, act botengine.Action, user botengine.User) botengine.ThinkResponse {
	actx := appengine.NewContext(ctx.Req)
	tk := os.Getenv("FB_PAGE_ACCESS_TOKEN")

	if strings.Contains(msg.Text, "城崎") {
		sendKinosakiRecommends(actx, user.ID, tk)
	} else if act.Kind == "action/book(1)" {
		sendConfirmButton(actx, user.ID, tk)
	} else if act.Kind == "action/confirm_book(1)" {
		sendReceipt(actx, user.ID, tk)
		return botengine.ThinkResponse{
			Message: botengine.Message{
				Text: "ご予約ありがとうございました。",
			},
		}
	} else {
		return botengine.ThinkResponse{
			Message: botengine.Message{
				Text: "whoa!",
			},
		}
	}

	return botengine.ThinkResponse{}
}

func sendKinosakiRecommends(ctx context.Context, sender int64, accessToken string) {
	els := []botengine.FBSendMessageTemplateElement {
		botengine.FBSendMessageTemplateElement{
			Title: "城崎温泉 かがり火の宿「大西屋」",
			ImageURL: "http://www.suisyou.com/img/main.jpg",
			Subtitle: "新和風数奇屋造り「全館たたみ敷き」のお宿",
			Buttons: []botengine.FBSendMessageTemplateButton {
				botengine.FBSendMessageTemplateButton{
					Type: "web_url",
					URL: "http://www.suisyou.com/index.html",
					Title: "詳細をみる",
				},
				botengine.FBSendMessageTemplateButton{
					Type: "postback",
					Title: "予約する",
					Payload: "action/book(1)",
				},
			},
		},
		botengine.FBSendMessageTemplateElement{
			Title: "城崎温泉 温もりの宿 おけ庄",
			ImageURL: "http://img.travel.rakuten.co.jp/share/image_up/7963/LARGE/D4Uu22.jpeg",
			Subtitle: "川沿いの灯る灯籠など城崎の景色を楽しめます",
			Buttons: []botengine.FBSendMessageTemplateButton {
				botengine.FBSendMessageTemplateButton{
					Type: "web_url",
					URL: "http://travel.rakuten.co.jp/HOTEL/7963/7963.html",
					Title: "詳細をみる",
				},
				botengine.FBSendMessageTemplateButton{
					Type: "postback",
					Title: "予約する",
					Payload: "action/book(1)",
				},
			},
		},
		botengine.FBSendMessageTemplateElement{
			Title: "城崎温泉　湯楽",
			ImageURL: "https://www.tocoo.jp/cnvimages/720/407/HotelImage/4001769/Sisetu_Touroku_Henkou/Kihon/main.jpg",
			Subtitle: "「お客様の五感に訴える…」これが湯楽の接客テーマ。",
			Buttons: []botengine.FBSendMessageTemplateButton {
				botengine.FBSendMessageTemplateButton{
					Type: "web_url",
					URL: "https://www.tocoo.jp/detail/4001769",
					Title: "詳細をみる",
				},
				botengine.FBSendMessageTemplateButton{
					Type: "postback",
					Payload: "action/book(1)",
					Title: "予約する",
				},
			},
		},
	}

	fbr := botengine.NewFBMessageWithTemplateRequest(sender, els)

	cli := urlfetch.Client(ctx)
	botengine.DispatchFBSendMessageRequest(cli, accessToken, fbr)
}

func sendConfirmButton(ctx context.Context, sender int64, accessToken string) {
	fbr := botengine.NewFBMessageWithButtons(sender, "城崎温泉 かがり火の宿「大西屋」を予約します。よろしいですか？", []botengine.FBSendMessageTemplateButton{
		botengine.FBSendMessageTemplateButton{
			Type: "postback",
			Payload: "action/confirm_book(1)",
			Title: "はい",
		},
		botengine.FBSendMessageTemplateButton{
			Type: "postback",
			Payload: "action/reject_book(1)",
			Title: "いいえ",
		},
	})
	cli := urlfetch.Client(ctx)
	botengine.DispatchFBSendMessageRequest(cli, accessToken, fbr)
}

func sendReceipt(ctx context.Context, sender int64, accessToken string) {
	fbr := botengine.NewFBMessageWithReceipt(sender, "乾 夏衣", "BK-" + strconv.Itoa(rand.Intn(1000000)), "JPY", "VISA", 32400, []botengine.FBSendMessageTemplateReceiptElement{
		botengine.FBSendMessageTemplateReceiptElement{
			Title: "城崎温泉 かがり火の宿「大西屋」",
			Quantity: 1,
			Price: 32400,
			Currency: "JPY",
			ImageURL: "http://www.suisyou.com/img/main.jpg",
		},
	})
	cli := urlfetch.Client(ctx)
	botengine.DispatchFBSendMessageRequest(cli, accessToken, fbr)
}