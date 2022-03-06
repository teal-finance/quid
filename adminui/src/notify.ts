import { ToastServiceMethods } from "primevue/toastservice";


const useNotify = function (toast: ToastServiceMethods) {
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
    done(content: string) {
      toast.add({ severity: 'success', summary: 'Done', detail: content, life: 1500, group: "main" });
    }
  }
}

export default useNotify;