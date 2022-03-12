import { GroupContract } from "./contract";

export default class Group {
  id: number;
  name: string;

  constructor({ id, name }: { id: number, name: string }) {
    this.id = id;
    this.name = name;
  }

  static fromContract(data: GroupContract): Group {
    return new Group({ id: data.id, name: data.name })
  }
}