package botengine
import (
	"testing"
	"bytes"
	"net/http"
	"os"
)

var _examplePayload1 = `{
  "object":"page",
  "entry":[
    {
      "id": "PAGE_ID",
      "time":1457764198246,
      "messaging":[
        {
          "sender":{
            "id": 123456
          },
          "recipient":{
            "id": "PAGE_ID"
          },
          "timestamp":1457764197627,
          "message":{
            "mid":"mid.1457764197618:41d102a3e1ae206a38",
            "seq":73,
            "text":"hello, world!"
          }
        }
      ]
    }
  ]
}`
var _example1 = []byte(_examplePayload1)
var example1Reader = bytes.NewReader(_example1)

func TestDecodingPayload(t *testing.T) {
	err, payload := DecodeWebhookPayload(example1Reader)

	if err != nil {
		t.Fatal("Err is not nil: %s", err.Error())
	}

	if len(payload.Entry) != 1 {
		t.Fatal("Entry size is not equal to 1.")
	}

	if payload.Entry[0].Time != 1457764198246 {
		t.Fatal("Time is not equal to 1457764198246.")
	}

	if len(payload.Entry[0].Messaging) != 1 {
		t.Fatal("Messaging size is not equal to 1.")
	}

	msg := payload.Entry[0].Messaging[0]

	if msg.Sender.Id != 123456 {
		t.Fatal("Sender id is malformed.")
	}

	if msg.Timestamp != 1457764197627 {
		t.Fatal("Timestamp is not equal to 1457764197627.")
	}

	m := msg.Message

	if m.Mid != "mid.1457764197618:41d102a3e1ae206a38" {
		t.Fatal("Mid is not equal to mid.1457764197618:41d102a3e1ae206a38.")
	}

	if m.Seq != 73 {
		t.Fatal("Seq is not equal to 73.")
	}

	if m.Text != "hello, world!" {
		t.Fatal("Text is not equal to 'hello, world!'")
	}
}