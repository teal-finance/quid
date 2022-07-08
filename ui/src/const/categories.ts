import Category from "@/models/category";

const serverAdminCategories = new Set<Category>([
  new Category({ title: "Namespaces", icon: "eos-icons:namespace", url: "/namespaces", type: "serverAdmin" }),
  new Category({ title: "Admins", icon: "eos-icons:admin-outlined", url: "/admins", type: "serverAdmin" }),
  new Category({ title: "Organisations", icon: "fluent:organization-12-regular", url: "/org", type: "serverAdmin" })
]);

const nsAdminCategories = new Set<Category>([
  new Category({ title: "Groups", icon: "clarity:group-solid", url: "/group", type: "nsAdmin" }),
  new Category({ title: "Users", icon: "clarity:user-solid", url: "/user", type: "nsAdmin" }),
]);

export { serverAdminCategories, nsAdminCategories };