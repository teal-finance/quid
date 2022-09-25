## Python

### Decode tokens

```python
import jwt

try:
    payload = jwt.decode(token, key, algorithms=['HS256'])
except jwt.ExpiredSignatureError:
    # ...
```

Payload example:

```json
{
  "usr": "jane",
  "grp": ["group1", "group2"],
  "org": ["organization1", "organization2"],
  "exp": 1595950745
}
```

Note: `"exp"` is the expiration timestamp in [Unix time](https://en.wikipedia.org/wiki/Unix_time) format (seconds since 1970).

### Flask server example

```python
from datetime import datetime
from flask import Flask, abort
from flask_cors import CORS, cross_origin
from flask import request
import jwt
import os

"""
To run:
export FLASK_APP=server.py
flask run
"""

app = Flask(__name__)
cors = CORS(app)
app.config['CORS_HEADERS'] = 'Content-Type'

# the key for the "demo" namespace
key = os.environ["QUID_DEMO_KEY"]


def verify_token(token):
    try:
        payload = jwt.decode(token, key, algorithms=['HS256'])
    except jwt.ExpiredSignatureError:
        return False
    print("Token payload", payload)
    date = datetime.fromtimestamp(payload["exp"])
    print("Token expiration date", date)
    return True


@app.route('/')
@cross_origin()
def main_route():
    token = request.headers["Authorization"].split(" ")[1]
    is_valid = verify_token(token)
    if (is_valid is True):
        return {"response": "ok"}
    else:
        abort(401)
```