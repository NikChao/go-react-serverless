export interface Snack {
  id: string;
  open: boolean;
  timeout: number;
  color: "danger" | "success";
  text: string;
  onClose(): void;
}
