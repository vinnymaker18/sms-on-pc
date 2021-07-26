import json
import random
import requests
import sys

SERVER_URL = 'https://minibox.vinayemani.xyz/sms'
USER_ID = 1
ORIGIN_ADDR = '+918825400845'

def random_text_msg():
    msg_len = random.randint(5, 20)
    chars = 'abcdefghi '
    return ''.join(random.choice(chars) for _ in range(msg_len))


def get_sms():
    response = requests.get(SERVER_URL, params=dict(userid=USER_ID))
    return json.loads(response.content)


def post_new_sms():
    response = requests.post(SERVER_URL, data=dict(userid=USER_ID, origin=ORIGIN_ADDR, msgbody=random_text_msg()))

def main():
    random.seed()
    args = sys.argv[1:]
    if args and args[0] == 'list':
        for msg in get_sms():
            print(msg)
    else:
        post_new_sms()


if __name__ == '__main__':
    main()
