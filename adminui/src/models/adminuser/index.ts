import { User as SwUser } from "@snowind/state";
import { UserType } from "./types";

export default class AdminUser extends SwUser {
  devRefreshToken: string | null = null;
  type: UserType = "nsAdmin";
}
