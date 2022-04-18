class QuidRequestError extends Error {
  hasToLogin: boolean;
  response: Response = new Response(null);

  constructor(message: string, response?: Response, hasToLogin?: boolean) {
    super(message);
    this.name = "QuidError";
    this.stack = (new Error() as any).stack;
    this.hasToLogin = hasToLogin ?? false;
    if (response) {
      this.response = response;
    }
  }
}

export { QuidRequestError };