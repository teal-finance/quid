## Nodejs

Decode an access token with a Nodejs server

```typescript
import express, { Express, NextFunction, Request, Response } from 'express';
import bodyParser from 'body-parser';
import helmet from 'helmet';
import dotenv from 'dotenv';
import cors from 'cors';
import * as jwt from 'jsonwebtoken';
import { conf } from "./conf";

dotenv.config();

const PORT = process.env.PORT || 5714;
const app: Express = express();
const keyBin = Buffer.from(conf.namespaceKey, 'hex');

app.use(cors({ origin: ["http://localhost:3000", "http://localhost:5173"], credentials: true }))
app.use(helmet());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
// jwt middleware
app.use(verifyjwt);

function verifyjwt(req: Request, res: Response, next: NextFunction) {
  const bearer = req.headers['authorization']
  if (!bearer) {
    console.log("User has no token")
    return res.status(401).json('Unauthorized user')
  }
  const token = bearer.split(" ")[1]
  console.log("Verifying token", token)
  console.log("with hexadecimal key", conf.namespaceKey)
  try {
    jwt.verify(token, keyBin, function (err, decoded) {
      if (err) {
        console.log("Token unverified", err);
        return res.status(401).json('Wrong token')
      } else {
        console.log("Decoded token", decoded)
      }
    });
  } catch (e) {
    return res.status(401).json('Token not valid')
  }
  next()
}

app.get('/', (req: Request, res: Response) => {
  //console.log("Request with auth header:", req.header("authorization"));
  res.send({ "response": "ok" })
});

app.listen(PORT, () => console.log(`Running on ${PORT} ⚡`));
```
