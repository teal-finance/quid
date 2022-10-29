import { ResponseErrorParams } from "./interfaces";

/** Response error class */
class ResponseError extends Error {
  status: number | null;
  statusText: string | null;
  content: Record<string, any> | Array<any> | null;
  text: string | null;

  constructor(params: ResponseErrorParams) {
    super(params.statusText);
    this.name = "ResponseError";
    this.status = params.status ?? null;
    this.statusText = params.statusText ?? null;
    this.content = params?.content ?? null;
    this.text = params?.text ?? null;
    let _message = "";
    if (params.errMsg) {
      _message = params.errMsg
    }
    _message = `${params.status.toString()} ${params.statusText}`
    this.message = _message;
  }

  static async create(response: Response): Promise<ResponseError> {
    //console.log("HEAD CT", response.headers.get("Content-Type"));
    let content: Record<string, any> | Array<any> | undefined = undefined;
    let text: string | undefined = undefined;
    if (response.headers.get("Content-Type")?.startsWith("application/json")) {
      content = await response.json()
    } else {
      text = await response.text()
    }
    return new ResponseError({
      status: response.status,
      statusText: response.statusText,
      content: content,
      text: text,
    })
  }
}

export { ResponseError };