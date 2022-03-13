import { ConfirmationOptions } from "primevue/confirmationoptions";

interface ConfirmOptions {
  require: (option: ConfirmationOptions) => void;
  close: () => void;
}

interface NotifyService {
  error: (content: string) => void;
  warning: (title: string, content: string, timeOnScreen?: number) => void;
  success: (title: string, content: string, timeOnScreen?: number) => void;
  done: (content: string) => void;
  confirmDelete: (msg: string, onConfirm: CallableFunction, onReject?: CallableFunction, title?: string) => void;
  toastSuccess: (msg: string, delay?: number | undefined) => void;
  toastInfo: (msg: string, delay?: number | undefined) => void;
}

export { ConfirmOptions, NotifyService }