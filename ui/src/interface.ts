import { ConfirmationOptions } from "primevue/confirmationoptions";

interface ConfirmOptions {
  require: (option: ConfirmationOptions) => void;
  close: () => void;
}

interface UserStatusContract {
  admin_type: "NsAdmin" | "QuidAdmin";
  username: string;
  ns: {
    id: number;
    name: string;
  }
}

interface NotifyService {
  error: (content: string) => void;
  warning: (title: string, content: string, timeOnScreen?: number) => void;
  success: (title: string, content: string, timeOnScreen?: number) => void;
  done: (content: string) => void;
  confirmDelete: (msg: string, onConfirm: CallableFunction, onReject?: CallableFunction, title?: string) => void;
}

type AlgoType = "HMAC" | "HS256" | "HS384" | "HS512" | "RS256" | "RS384" | "RS512" | "PS256" | "PS384" | "PS512" | "ES256" | "ES384" | "ES512" | "EDDSA";

export { ConfirmOptions, NotifyService, UserStatusContract, AlgoType }
