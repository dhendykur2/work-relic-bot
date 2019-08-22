# Work Relic Bot
[WorkRelicBot](https://t.me/WorkRelicBot) is Telegram bot to monitor your daily task.

# Run in your local
```console
$ git clone https://github.com/dhendykur2/work-relic-bot
```
Make sure you already install [dep](https://github.com/golang/dep)

```console
$ cd work-relic-bot
```

```console
$ dep ensure
```

 ***Create the new Bot? check https://core.telegram.org/bots***

```console
 $ cp credentials_example.json credentials.json
 $ nano credentials.json
```

if you have python in your local machine you can change your telegram webhook url by using changeWebhookTelegram.py script. Or you can just type in your browser with this format https://api.telegram.org/botBOT_TOKEN/setWebhook?url=WEBHOOK_URL

```python
import requests

def setWebhook(token, webhook):
    if token == "" and webhook == "":
        print("Fill token and webhook")
        return
    res = requests.get("https://api.telegram.org/bot{}/setWebhook?url={}".format(token, url))
    print(res)
    print(res.content)
    return

bot_token = "" # your bot Token
webhook_url = "" # your webhook_url
setWebhook(bot_token, webhook_url)
```

Thanks to [ngrok](https://ngrok.com/) for port forwarding while development. So the **webhook_url** will be your ngrok url that port-forward to your localhost. NOTE: while changing ngrok url as a webhook please select https.

# Ngrok Example
```console
$ ./ngrok http 3000
```
Output
```console
ngrok by @inconshreveable                                                                                                                                        (Ctrl+C to quit)

Session Status                online
Account                        (Plan: Free)
Version                       2.3.34
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    http://SAMPLE.ngrok.io -> http://localhost:3000
Forwarding                    https://SAMPLE.ngrok.io -> http://localhost:3000

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

run the python script then fill the **webhook_url** with https://SAMPLE.ngrok.io.

# Run
```console
$ go run main.go
```

# Test Case
under development.

