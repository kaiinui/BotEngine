# BotEngine
Only a function you need to make a bot!

Supported platforms: Facebook Messenger

At first BotEngine supports `python`, while it will support various languages such as `ruby`, `php`, `node.js`, `Java`, `Go`.

```python
def on_message_received(context, message, action):
    if context.state == u"initial"
        return BotResponse.elements_and_state(recommended_places, u"search_place")
        
    elif context.state == u"search_place" and action == u"book(place:2)"
        context.set_meta_data(u"place_to_book", 2)
        
        return BotResponse.text_and_state(u"Are you sure to book the place 2?", u"confirm_booking")
        
    elif context.state == u"confirm_booking"
        place = context.get_meta_data(u"place_to_book")
        context.set_meta_data(u"place_booked", place)
        
        return BotResponse.text_and_state(u"You booked place " + place + u".", u"booked")
        
    elif context.state == u"booked"
        place = context.get_meta_data(u"place_booked")
    
        return BotResponse.text_and_state(u"Forget where you booked? You booked place " + place u"!")
```
