import NamespaceContract from "../namespace/contract";

interface UserContract {
  id: number;
  name: string;
  namespace: NamespaceContract;
}

export default UserContract;
