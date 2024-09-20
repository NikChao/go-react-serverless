import { observer } from "mobx-react-lite";
import { CatalogService } from "../../../services/catalog-service";
import { ReceiptService } from "../../../services/receipt-service";
import { RawDataForBusiness } from "./screen";
import { RawDataForBusinessStore } from "./store";

export function createRawDataForBusiness(
  catalogService: CatalogService,
  receiptService: ReceiptService
) {
  const store = new RawDataForBusinessStore(catalogService, receiptService);

  return observer(() => {
    return (
      <RawDataForBusiness
        snacks={store.snacks}
        isUploadingReceipt={store.isUploadingReceipt}
        handleReceiptUpload={store.handleReceiptUpload}
      />
    );
  });
}
