type ColorVariant = "primary" | "secondary" | "neutral" | "light" | "success" | "warning" | "danger";

type PopToast = (msg: string, type: ColorVariant, delay?: number | undefined) => Promise<void>;

export { PopToast, ColorVariant }