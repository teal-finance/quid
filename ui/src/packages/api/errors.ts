import { ResponseErrorType } from "./interfaces";

/** Response error class */
class ResponseError extends Error {
  type: ResponseErrorType;
  response: Response;
  status: number;
  statusText: string;

  constructor(response: Response, type: ResponseErrorType, errMsg?: string) {
    const _resp = response.clone()
    super(_resp.statusText);
    this.name = `ResponseError (${type})`;
    this.stack = super.stack;
    this.response = _resp;
    this.status = _resp.status;
    this.statusText = _resp.statusText;
    this.type = type;
    let _message = "";
    if (errMsg) {
      _message = errMsg
    }
    if (type == "request") {
      _message = `${_resp.status.toString()} ${_resp.statusText}`
    }
    this.message = _message;
  }
}

export { ResponseError };