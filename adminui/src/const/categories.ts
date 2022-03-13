import Category from "@/models/category";

const categories = new Set<Category>([
  new Category({ title: "Namespaces", icon: "eos-icons:namespace", url: "/namespaces" }),
  new Category({ title: "Organisations", icon: "fluent:organization-12-regular", url: "/org" }),
  new Category({ title: "Admins", icon: "eos-icons:admin-outlined", url: "/admins" }),
  new Category({ title: "Groups", icon: "clarity:group-solid", url: "/group" }),
  new Category({ title: "Users", icon: "clarity:user-solid", url: "/user" }),
]);

export default categories;