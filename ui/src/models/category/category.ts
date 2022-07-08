import { UserType } from "../siteuser/types";

export default class Category {
  title: string;
  icon: string;
  url: string | null;
  type: UserType;

  constructor({ title, icon, url = null, type = "serverAdmin" }: {
    title: string, icon: string, url?: string | null, type: UserType
  }) {
    this.title = title;
    this.icon = icon;
    this.url = url;
    this.type = type
  }
}
