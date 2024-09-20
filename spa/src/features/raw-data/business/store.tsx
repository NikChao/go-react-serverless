import { makeAutoObservable } from "mobx";
import { Snack } from "../../../models/snack";
import { CatalogService } from "../../../services/catalog-service";
import { ReceiptService } from "../../../services/receipt-service";
import { chunk } from "../../../utils/chunk";
import { wait } from "@testing-library/user-event/dist/utils";

const FILE_CHUNK_SIZE = 4;

export class RawDataForBusinessStore {
  isUploadingReceipt: boolean = false;
  snacks: Snack[] = [];

  constructor(
    private readonly catalogService: CatalogService,
    private readonly receiptService: ReceiptService
  ) {
    makeAutoObservable(this);
  }

  public handleReceiptUpload = async (fileList: FileList) => {
    this.isUploadingReceipt = true;

    try {
      // 4 at a time to not get throttled
      const chunks = chunk(fileList, FILE_CHUNK_SIZE);

      for (const chunk of chunks) {
        await Promise.all(
          chunk.map((file) =>
            this.receiptService.uploadBusinessReceiptToPresignedUrl(file)
          )
        );

        // Wait for throttling reasons
        await wait(200);
      }

      const id = Math.floor(Math.random() * 100_000_000).toString();
      this.snacks.push({
        id,
        text: "Thank you!",
        open: true,
        timeout: 1200,
        color: "success",
        onClose: this.removeSnackById(id),
      });
    } catch {
      const id = Math.floor(Math.random() * 100_000_000).toString();
      this.snacks.push({
        id,
        text: "Receipt is too big!",
        open: true,
        timeout: 1200,
        color: "danger",
        onClose: this.removeSnackById(id),
      });
    } finally {
      this.isUploadingReceipt = false;
    }
  };

  private removeSnackById(id: string) {
    return () => {
      this.snacks = this.snacks.map((snack) => {
        if (snack.id === id) {
          snack.open = false;
          setTimeout(() => {
            this.snacks = this.snacks.filter((snack) => snack.id !== id);
          }, 250);
        }

        return snack;
      });
    };
  }
}
