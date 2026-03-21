export abstract class BaseViewModel<T> {
  constructor(data: T) {
    Object.assign(this, data);
  }

  toJSON(): T {
    return JSON.parse(JSON.stringify(this));
  }
}
