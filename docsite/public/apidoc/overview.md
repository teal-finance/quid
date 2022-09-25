## Overview

**Quid** is a [JWT][jwt] server (frontend + backend + client libraries)
to manage Administrators, Users, **Refresh Tokens** and **Access Tokens**
in independent **Namespaces** providing signature verification for the following algorithms:

- HS256 = HMAC using SHA-256
- HS384 = HMAC using SHA-384
- HS512 = HMAC using SHA-512
- RS256 = RSASSA-PKCS1-v1_5 using 2048-bits RSA key and SHA-256
- RS384 = RSASSA-PKCS1-v1_5 using 2048-bits RSA key and SHA-384
- RS512 = RSASSA-PKCS1-v1_5 using 2048-bits RSA key and SHA-512
- ES256 = ECDSA using P-256 and SHA-256
- ES384 = ECDSA using P-384 and SHA-384
- ES512 = ECDSA using P-521 and SHA-512
- EdDSA = Ed25519

[jwt]: https://wikiless.org/wiki/JSON_Web_Token "JSON Web Token"

![Authentication flow chart](/img/authentication-flow.png)

1. First, the user logs in with **Namespace** + **Username** + **Password**.
   The **Namespace** is usually the final application name,
   represented by _Application API_ at the bottom of the previous diagram.

2. Then, the client (e.g. JS code) receives a **Refresh Token**
   that is usually valid for a few hours
   to avoid to log again during the working session.

3. The client sends this **Refresh Token** to get an **Access Token**
   that is valid for a short time,
   usually a few minutes, say 10 minutes.
   So the client must _refresh_ its **Access Token** every 10 minutes.

4. During these 10 minutes,
   the client can request the _Application API_
   with the same **Access Token**.

5. When the _Application API_ receives a request from the client,
   it checks the [JWT][jwt] signature and expiration time.
   The **Access Token** is stateless:
   the _Application API_ does not need to store any information
   about the user (the **Access Token** content is enough).