import requests


def setWebhook(token, webhook):
    if token == "" and webhook == "":
        print("Fill token and webhook")
        return
    res = requests.get("https://api.telegram.org/bot{}/setWebhook?url={}".format(token, url))
    print(res)
    print(res.content)
    return

bot_token = ""
webhook_url = ""
setWebhook(bot_token, webhook_url)


