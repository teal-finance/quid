class QuidError extends Error {
  public hasToLogin: boolean;

  constructor(public message: string, public _hasToLogin: boolean = false) {
    super(message);
    this.name = "QuidError";
    this.stack = (new Error() as any).stack; // eslint-disable-line
    this.hasToLogin = _hasToLogin;
  }
}

export default QuidError;