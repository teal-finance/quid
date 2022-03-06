
export default class Category {
  title: string;
  icon: string;
  url: string | null;

  constructor({ title, icon, url = null }: { title: string, icon: string, url?: string | null }) {
    this.title = title;
    this.icon = icon;
    this.url = url;
  }
}
