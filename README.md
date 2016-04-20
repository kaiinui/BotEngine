# BotEngine

Only a function you need to make a Facebook Messenger bot!

## Introduction

BotEngine is a FB Messenger chat bot framework which works on Google App Engine. 

The design is deadly simple. The code you must write is only a function `onThink`. BotEngine treats everything without "thinking", such as managing states and contexts, requesting to FB Messenger API or handling webhooks.

Currently it supports only `Go`, however I plan to support `python`, `ruby`, `node.js`, `Java`.

```go
/**
 *  A simple hotel booking bot.
 */

func onThink(ctx botengine.Context, msg botengine.Message, act botengine.Action, user botengine.User) botengine.ThinkResponse {
	if ctx.State == botengine.InitialState && strings.Contains(msg.Text, "hotel") {
	    // Search for hotels and prepare `carouselElements` ...
	
	    return botengine.Carousel("search_for_hotel", carouselElements)
	} else if ctx.State == "search_for_hotel" {
	    ctx.PutSessionValue("hotel_id", act.Attributes["hotel_id"])
	    
	    return botengine.Confirm("confirm_book", "Are you sure to book this hotel?")
	} else if ctx.State == "confirm_book" {
	    id := ctx.GetSessionValue("hotel_id")
	    
	    if act.Kind == botengine.ConfirmYesActionKind {
	        return botengine.Text(botengine.InitialState, "Thank you. We booked your hotel. ID: " + id)
	    } else {
	        return botengine.Text(botengine.InitialState, "I'm glad to see you again!")
	    }
	}

	return botengine.Text(botengine.InitialState, "Not supported operation.")
}

func init() {
    s := botengine.NewServer(os.Getenv("FB_PAGE_ACCESS_TOKEN"))
    s.Run(onThink)
}
```
