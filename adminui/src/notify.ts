import { ToastServiceMethods } from "primevue/toastservice";
import { ConfirmOptions, NotifyService } from "./interface";


const useNotify = function (toast: ToastServiceMethods, confirm: ConfirmOptions): NotifyService {
  return {
    error: (content: string) => {
      toast.add({ severity: 'error', summary: 'Error', detail: content, group: "main" });
    },
    warning: (title: string, content: string, timeOnScreen = 5000) => {
      toast.add({ severity: 'error', summary: title, detail: content, life: timeOnScreen, group: "main" });
    },
    success: (title: string, content: string, timeOnScreen = 1500) => {
      toast.add({ severity: 'success', summary: title, detail: content, life: timeOnScreen, group: "main" });
    },
    done: (content: string) => {
      toast.add({ severity: 'success', summary: 'Done', detail: content, life: 1500, group: "main" });
    },
    confirmDelete: (msg: string, onConfirm: CallableFunction, onReject: CallableFunction = () => null, title = "Delete") => {
      confirm.require({
        message: msg,
        header: title,
        icon: 'pi pi-info-circle',
        acceptClass: 'p-button-danger',
        accept: () => onConfirm(),
        reject: () => onReject(),
      });
    },
  }
}

export default useNotify;