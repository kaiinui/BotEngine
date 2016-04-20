package botengine
import (
	"io"
	"encoding/json"
	"net/http"
	"bytes"
)

func DecodeWebhookPayload(r io.Reader) (error, FBWebhookPayload) {
	var payload FBWebhookPayload
	err := json.NewDecoder(r).Decode(&payload)

	return err, payload
}

type FBWebhookPayload struct {
	Object string `json:"object"`
	Entry []FBWebhookEntry `json:"entry"`
}

type FBWebhookEntry struct {
	Time int64 `json:"time"`
	Messaging []FBWebhookMessaging `json:"messaging"`
}

type FBWebhookMessaging struct {
	Sender FBWebhookMessagingSender `json:"sender"`
	Timestamp int64 `json:"timestamp"`
	Message *FBWebhookMessagingMessage `json:"message"` // Message-Received Callback
	Delivery *FBWebhookMessagingDelivery `json:"delivery"` // Message-Delivered Callback
	Optin *FBWebhookMessagingOptin `json:"optin"` // Authentication Callback
	Postback *FBWebhookMessagingPostback `json:"postback"` // Postback Callback
}

type FBWebhookMessagingSender struct {
	Id int64 `json:"id"`
}

// Postback Callback
type FBWebhookMessagingPostback struct {
	Payload string `json:"payload"`
}

// Authentication Callback
type FBWebhookMessagingOptin struct {
	Ref string `json:"ref"`
}

// Message-Delivered Callback
type FBWebhookMessagingDelivery struct {
	Mids []string `json:"mids"`
	Watermark int64 `json:"watermark"`
	Seq int `json:"seq"`
}

// Message-Received Callback
type FBWebhookMessagingMessage struct {
	Mid string `json:"mid"`
	Seq int `json:"seq"`
	Text string `json:"text"`
	Attachments []FBWebhookMessagingMessageAttachment `json:"attachments"`
}

type FBWebhookMessagingMessageAttachment struct {
	Type string `json:"type"` // `image`, `video`, `audio`
	Payload FBWebhookMessagingMessageAttachmentPayload `json:"payload"`
}

type FBWebhookMessagingMessageAttachmentPayload struct {
	 URL string `json:"url"`
}

var fbSendMessageApiEndpoint = "https://graph.facebook.com/v2.6/me/messages?access_token="

func DispatchFBSendMessageRequest(cli *http.Client, accessToken string, req FBSendMessageRequest) (*http.Response, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}

	return cli.Post(fbSendMessageApiEndpoint + accessToken, "application/json", buf)
}

func NewFBMessageWithTextRequest(recipient int64, text string) FBSendMessageRequest {
	return FBSendMessageRequest{
		Recipient: FBSendMessageRecipient{recipient},
		Message: FBSendMessageMessage{Text: text},
	}
}

func NewFBMessageWithTemplateRequest(recipient int64, elements []FBSendMessageTemplateElement) FBSendMessageRequest {
	interfaceElements := make([]interface{}, len(elements))
	for i, v := range elements {
		interfaceElements[i] = v
	}

	return FBSendMessageRequest{
		Recipient: FBSendMessageRecipient{recipient},
		Message: FBSendMessageMessage{
			Attachment: &FBSendMessageAttachment{
				Type: "template",
				Payload: FBSendMessageAttachmentPayload{
					TemplateType: "generic",
					Elements: interfaceElements,
				},
			},
		},
	}
}

func NewFBMessageWithImageURL(recipient int64, url string) FBSendMessageRequest {
	return FBSendMessageRequest{
		Recipient: FBSendMessageRecipient{recipient},
		Message: FBSendMessageMessage{
			Attachment: &FBSendMessageAttachment{
				Type: "image",
				Payload: FBSendMessageAttachmentPayload{
					URL: url,
				},
			},
		},
	}
}

func NewFBMessageWithButtons(recipient int64, text string, buttons []FBSendMessageTemplateButton) FBSendMessageRequest {
	return FBSendMessageRequest{
		Recipient: FBSendMessageRecipient{recipient},
		Message: FBSendMessageMessage{
			Attachment: &FBSendMessageAttachment{
				Type: "template",
				Payload: FBSendMessageAttachmentPayload{
					TemplateType: "button",
					Text: text,
					Buttons: buttons,
				},
			},
		},
	}
}

func NewFBMessageWithReceipt(recipient int64, name string, orderNumber string, currency string, paymentMethod string, totalCost int, elements []FBSendMessageTemplateReceiptElement) FBSendMessageRequest {
	interfaceElements := make([]interface{}, len(elements))
	for i, v := range elements {
		interfaceElements[i] = v
	}

	return FBSendMessageRequest{
		Recipient: FBSendMessageRecipient{recipient},
		Message: FBSendMessageMessage{
			Attachment: &FBSendMessageAttachment{
				Type: "template",
				Payload: FBSendMessageAttachmentPayload{
					TemplateType: "receipt",
					RecipientName: name,
					OrderNumber: orderNumber,
					Currency: currency,
					PaymentMethod: paymentMethod,
					Summary: &FBSendMessageTemplateSummary{
						TotalCost: totalCost,
					},
					Elements: interfaceElements,
				},
			},
		},
	}
}

type FBSendMessageRequest struct {
	Recipient FBSendMessageRecipient `json:"recipient"`
	Message FBSendMessageMessage `json:"message"`
}

type FBSendMessageRecipient struct {
	Id int64 `json:"id"`
}

type FBSendMessageMessage struct {
	Text string `json:"text"`
	Attachment *FBSendMessageAttachment `json:"attachment,omitempty"`
}

type FBSendMessageAttachment struct {
	Type string `json:"type"` // `template`, `image`
	Payload FBSendMessageAttachmentPayload `json:"payload"`
}

// Templates
type FBSendMessageAttachmentPayload struct {
	URL string `json:"url,omitempty"`
	TemplateType string `json:"template_type"` // `generic`, `button`
	// FBSendMessageTemplateReceiptElement || FBSendMessageTemplateElement
	Elements []interface{} `json:"elements,omitempty"`
	Text string `json:"text,omitempty"`
	Buttons []FBSendMessageTemplateButton `json:"buttons,omitempty"`

	// Receipt
	RecipientName string `json:"recipient_name,omitempty"`
	OrderNumber string `json:"order_number,omitempty"`
	Currency string `json:"currency,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Summary *FBSendMessageTemplateSummary `json:"summary,omitempty"`
}

type FBSendMessageTemplateElement struct {
	Title string `json:"title"`
	ImageURL string `json:"image_url"`
	Subtitle string `json:"subtitle"`
	Buttons []FBSendMessageTemplateButton `json:"buttons"`
}

type FBSendMessageTemplateButton struct {
	Type string `json:"type"` // `web_url`, `postback`
	URL string `json:"url,omitempty"`
	Title string `json:"title"`
	Payload string `json:"payload,omitempty"`
}

type FBSendMessageTemplateReceiptElement struct {
	Title string `json:"title"`
	Price int `json:"price"`
	Quantity int `json:"quantity"`
	Currency string `json:"currency"`
	ImageURL string `json:"image_url"`
}

type FBSendMessageTemplateSummary struct {
	TotalCost int `json:"total_cost"`
}