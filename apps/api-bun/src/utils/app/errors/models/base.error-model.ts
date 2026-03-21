export interface ErrorModel {
  success: boolean;
  message: string;
}

export class BaseError extends Error implements ErrorModel {
  public readonly success = false;

  constructor(message: string) {
    super(message);

    Object.setPrototypeOf(this, new.target.prototype);

    this.name = new.target.name;

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    }
  }

  toJSON(): ErrorModel {
    return {
      success: this.success,
      message: this.message,
    };
  }
}
