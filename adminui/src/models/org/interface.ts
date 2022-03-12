import OrgContract from "./contract";

interface OrgTable extends OrgContract {
  actions: Array<string>;
}

export { OrgTable }