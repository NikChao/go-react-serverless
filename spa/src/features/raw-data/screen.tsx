import {
  Box,
  IconButton,
  Input,
  Sheet,
  Snackbar,
  Stack,
  Table,
  Typography,
} from "@mui/joy";
import { Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Catalog, CatalogService } from "../../services/catalog-service";
import { Checkbox } from "@mui/joy";
import { StoreData, StoreName } from "../../services/grocery-service";
import FileUploadButton from "../../components/file-upload-button";
import { ReceiptService } from "../../services/receipt-service";
import { AddShoppingCart } from "@mui/icons-material";
import { GroceryListStore } from "../grocery/store";
import { observer } from "mobx-react-lite";

interface RawDataScreenProps {
  catalogService: CatalogService;
  receiptService: ReceiptService;
  groceryStore: GroceryListStore;
}

export const RawDataScreen = observer(
  ({ catalogService, receiptService, groceryStore }: RawDataScreenProps) => {
    const [isAddingGroceries, setIsAddingGroceries] = useState(false);
    const [isThankOpen, setIsThankOpen] = useState(false);
    const [isFileTooBigSnackbarOpen, setFileTooBigSnackbar] = useState(false);
    const [isUploadingReceipt, setIsUploadingReceipt] = useState(false);
    const [{ data }, setData] = useState<Catalog>({ data: [] });
    const [searchTerm, setSearchTerm] = useState("");
    const [checked, setChecked] = useState<Record<string, boolean>>({});

    async function fetchCatalog() {
      const catalog = await catalogService.getCatalog();
      setData(catalog);
    }

    function toggle(name: string) {
      const newChecked = {
        ...checked,
        [name]: !checked[name],
      };

      setChecked(newChecked);
    }

    async function handleReceiptUpload(file: File) {
      setIsUploadingReceipt(true);
      try {
        await receiptService.uploadReceiptToPresignedUrl(file);
        setIsThankOpen(true);
      } catch {
        setFileTooBigSnackbar(true);
      } finally {
        setIsUploadingReceipt(false);
      }
    }

    async function addSelectedItemsToStore() {
      setIsAddingGroceries(true);

      const itemsToAdd = Object.entries(checked)
        .filter(([_, isChecked]) => isChecked)
        .map(([key]) => key);
      await groceryStore.addGroceryItems(itemsToAdd);
      setChecked({});
      setIsAddingGroceries(false);
    }

    useEffect(() => {
      fetchCatalog();
      // eslint-disable-next-line
    }, []);

    const filteredData = data.filter(({ name }) =>
      name.toLocaleLowerCase().includes(searchTerm.toLocaleLowerCase())
    );

    const cost = data
      .filter(({ name }) => checked[name])
      .map((item) => {
        const allItemPricesFromSelectedStores = item.storeData
          .filter(
            (storeData) =>
              groceryStore.selectedStores.includes(storeData.storeName) &&
              storeData.price
          )
          .map((storeData) => parseFloat(storeData.price));
        return Math.min(...allItemPricesFromSelectedStores);
      })
      .reduce((previous, current) => previous + current, 0);

    const roundedCost = Math.floor(cost * 100) / 100;
    const stores = groceryStore.selectedStores;

    return (
      <>
        <Snackbar
          open={isThankOpen}
          onClose={() => setIsThankOpen(false)}
          color="success"
          variant="soft"
          autoHideDuration={1500}
          anchorOrigin={{ vertical: "top", horizontal: "center" }}
        >
          <Typography level="body-md">Thank you!</Typography>
        </Snackbar>
        <Snackbar
          open={isFileTooBigSnackbarOpen}
          onClose={() => setFileTooBigSnackbar(false)}
          color="danger"
          variant="soft"
          autoHideDuration={1500}
          anchorOrigin={{ vertical: "top", horizontal: "center" }}
        >
          <Typography level="body-md">File was too big to upload</Typography>
        </Snackbar>
        <Stack height="100%" boxSizing="border-box">
          <Stack p="16px">
            <Box
              display="flex"
              justifyContent="space-between"
              alignItems="center"
              width="100%"
            >
              <Typography level="body-md" fontWeight="bold">
                Grocery pricing data
              </Typography>
              <Link to="/">home</Link>
            </Box>
            <Box
              display="flex"
              justifyContent="space-between"
              alignItems="center"
              flexWrap="wrap"
              py={3}
              gap={3}
            >
              <Box display="flex" gap={4} flexWrap="wrap">
                <Checkbox
                  label="Aldi"
                  checked={stores.includes("aldi")}
                  onChange={() => groceryStore.toggleStore("aldi")}
                />
                <Checkbox
                  label="Coles"
                  checked={stores.includes("coles")}
                  onChange={() => groceryStore.toggleStore("coles")}
                />
                <Checkbox
                  label="Woolies"
                  checked={stores.includes("woolies")}
                  onChange={() => groceryStore.toggleStore("woolies")}
                />
                <Checkbox
                  label="Sam Cocos"
                  checked={stores.includes("sam cocos")}
                  onChange={() => groceryStore.toggleStore("sam cocos")}
                />
              </Box>
              <Box display="flex" alignItems="center" gap={2}>
                <Stack>
                  <Typography level="body-xs">
                    Want to support the project?
                  </Typography>
                  <Typography level="body-xs">Upload a receipt!</Typography>
                </Stack>
                <FileUploadButton
                  loading={isUploadingReceipt}
                  multiple={false}
                  onChange={handleReceiptUpload}
                />
              </Box>
            </Box>
            <Input
              variant="outlined"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search for an ingredient"
            />
          </Stack>

          <Sheet>
            <Table hoverRow stickyFooter>
              <thead>
                <tr>
                  <th style={{ width: "20px" }} />
                  <th colSpan={2}>Item</th>
                  <Header storeName="aldi" stores={stores} />
                  <Header storeName="coles" stores={stores} />
                  <Header storeName="woolies" stores={stores} />
                  <Header storeName="sam cocos" stores={stores} />
                </tr>
              </thead>
              <tbody>
                {filteredData.map((item) => (
                  <tr key={item.name}>
                    <td
                      style={{
                        display: "flex",
                        alignItems: "center",
                        width: "min-content",
                        paddingRight: "8px",
                        paddingLeft: "16px",
                      }}
                    >
                      <Checkbox
                        checked={checked[item.name] ?? false}
                        onChange={() => toggle(item.name)}
                      />
                    </td>
                    <td colSpan={2}>{item.name}</td>
                    <DataColumn
                      storeName="aldi"
                      storeData={item.storeData}
                      stores={stores}
                    />
                    <DataColumn
                      storeName="coles"
                      storeData={item.storeData}
                      stores={stores}
                    />
                    <DataColumn
                      storeName="woolies"
                      storeData={item.storeData}
                      stores={stores}
                    />
                    <DataColumn
                      storeName="sam cocos"
                      storeData={item.storeData}
                      stores={groceryStore.selectedStores}
                    />
                  </tr>
                ))}
              </tbody>
              <tfoot>
                <tr>
                  <td colSpan={2 + groceryStore.selectedStores.length}>
                    <Typography fontWeight="600" level="body-sm">
                      Total Cost:&nbsp;
                      <Typography fontWeight="400" level="body-sm">
                        {roundedCost}
                      </Typography>
                    </Typography>
                  </td>
                  <td>
                    <IconButton
                      onClick={addSelectedItemsToStore}
                      loading={isAddingGroceries}
                    >
                      <AddShoppingCart />
                    </IconButton>
                  </td>
                </tr>
              </tfoot>
            </Table>
          </Sheet>
        </Stack>
      </>
    );
  }
);

function Header({
  storeName,
  stores,
}: {
  storeName: StoreName;
  stores: StoreName[];
}) {
  if (!stores.includes(storeName)) {
    return null;
  }

  return <th>{storeName[0].toUpperCase() + storeName.slice(1)}</th>;
}

function DataColumn({
  storeName,
  stores,
  storeData,
}: {
  storeName: StoreName;
  stores: StoreName[];
  storeData: StoreData[];
}) {
  if (!stores.includes(storeName)) {
    return null;
  }

  const selectedStoreData = storeData.find((d) => storeName === d.storeName);

  if (!selectedStoreData) {
    return <td />;
  }

  const price = parseFloat(selectedStoreData.price);
  const allPrices = storeData
    .filter(({ storeName }) => stores.includes(storeName))
    .map(({ price }) => parseFloat(price));

  const isCheapest = price === Math.min(...allPrices);

  return (
    <td>
      <Stack
        justifyContent="flex-start"
        height="100%"
        overflow="hidden"
        pr="16px"
      >
        <Typography
          fontWeight={isCheapest ? 600 : undefined}
          textColor={isCheapest ? "success.400" : undefined}
        >
          {selectedStoreData?.price}
        </Typography>
        <Typography level="body-xs" overflow="ellipsis">
          {selectedStoreData?.itemName}
        </Typography>
      </Stack>
    </td>
  );
}
