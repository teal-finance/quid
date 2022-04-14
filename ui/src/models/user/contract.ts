import NamespaceContract from "../namespace/contract";

interface UserContract {
  id: number;
  username: string;
  namespace: NamespaceContract;
}

export default UserContract;