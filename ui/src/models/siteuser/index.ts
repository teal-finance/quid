import { User as SwUser } from "@snowind/state";
import { ref } from "@vue/reactivity";
import { useStorage } from "@vueuse/core";
import Namespace from "../namespace";
import NamespaceTable from "../namespace/interface";
import { UserType } from "./types";

export default class SiteUser extends SwUser {
  devRefreshToken: string | null = null;
  type = ref<UserType>("nsAdmin");
  namespace = useStorage("namespace", Namespace.empty().toTableRow());
  adminUrl = "/ns";

  get mustSelectNamespace(): boolean {
    return this.namespace.value.id == 0
    //&& user.type == "serverAdmin";
  }

  changeNs(nst: NamespaceTable) {
    this.namespace.value = nst;
  }

  resetNs() {
    this.namespace.value = Namespace.empty().toTableRow();
  }
}
