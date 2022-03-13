import { User as SwUser } from "@snowind/state";
import { UserType } from "./types";

export default class User extends SwUser {
  devRefreshToken: string | null = null;
  type: UserType = "nsAdmin";
}
