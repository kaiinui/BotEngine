# BotEngine
Only a function you need to make a Facebook Messenger bot!

```go
/**
 *  A simple hotel booking bot.
 */

func onThink(ctx botengine.Context, msg botengine.Message, act botengine.Action, user botengine.User) botengine.ThinkResponse {
	if ctx.State == botengine.InitialState && strings.Contains(msg.Text, "hotel") {
	    // Search for hotels and prepare `carouselElements` ...
	
	    return botengine.Carousel("search_for_hotel", carouselElements)
	} else if ctx.State == "search_for_hotel" {
	    ctx.PutValue("hotel_id", act.Attributes["hotel_id"])
	    
	    return botengine.Confirm("confirm_book", "Are you sure to book this hotel?")
	} else if ctx.State == "confirm_book" {
	    
	    if act.Kind == botengine.ConfirmYesActionKind {
	        return botengine.Text(botengine.InitialState, "Thank you. We booked your hotel!")
	    } else {
	        return botengine.Text(botengine.InitialState, "I'm glad to see you again!")
	    }
	}

	return botengine.Text(botengine.InitialState, "Not supported operation.")
}
```
