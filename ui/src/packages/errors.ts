class ResponseError extends Error {
  response: Response;

  constructor(response: Response) {
    super(response.statusText);
    this.name = "ResponseError";
    this.stack = (new Error() as any).stack;
    this.response = response;
  }
}

export { ResponseError };