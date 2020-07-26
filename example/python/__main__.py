import sys
import json
import jwt
from urllib.parse import urlencode
from urllib.request import Request, urlopen
from prompt_toolkit import prompt

SERVER_URI = "http://localhost:8082"
TOKEN_TTL = "1h"


def main():
    global SERVER_URI, TOKEN_TTL
    args = sys.argv[1:]
    if (len(args) != 2):
        print("Pass the namespace name as first the argument "
              "and the namespace decoding key as the second")
        return
    namespace = args[0]
    key = args[1]
    # prompt for username and passord
    print("Asking a token to server "+SERVER_URI)
    username = prompt('Username: ')
    pwd = prompt('Password: ')
    # request the token
    url = SERVER_URI+"/request_token/"+TOKEN_TTL
    post_fields = {'username': username, 'password':
                   pwd, 'namespace': namespace}
    request = Request(url, urlencode(post_fields).encode())
    res = urlopen(request).read().decode()
    data = json.loads(res)
    token = data["token"]
    print("Received token from server: "+token)
    print("Decoding the token")
    # read the token
    pl = jwt.decode(token, key, algorithms=['HS256'])
    print("Payload", pl)


if __name__ == "__main__":
    main()
