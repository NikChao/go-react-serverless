import { useEffect, ChangeEvent } from "react";
import { Stack, Typography, Box, Button, Input } from "@mui/joy";
import { FormEvent } from "react";
import {
  GroceryList as GroceryListModel,
  StoreName,
} from "../../services/grocery-service";
import { poll } from "../../utils/poll";
import { GroceryList } from "./grocery-list";
import { MoreInfoSheet } from "../more-info/more-info-sheet";

const POLL_INTERVAL = 8000;

export function GroceryScreen({
  groceryItemText,
  groceryList,
  isUpdating,
  endIcons,
  householdId,
  selectedStores,
  handleGroceryItemTextChange,
  checkGroceryItem,
  clearCheckedItems,
  createGroceryItem,
  initializeGroceryList,
  fetchGroceries,
  toggleStore,
}: {
  groceryItemText: string;
  isUpdating: boolean;
  householdId: string;
  endIcons: React.ReactNode;
  groceryList: GroceryListModel;
  selectedStores: StoreName[];
  handleGroceryItemTextChange(e: ChangeEvent<HTMLInputElement>): void;
  checkGroceryItem(id: string): void;
  clearCheckedItems(householdId: string): void;
  createGroceryItem(
    householdId: string
  ): (e: FormEvent<HTMLFormElement>) => void;
  initializeGroceryList(householdId: string): void;
  fetchGroceries(householdId: string): void;
  toggleStore(storeName: StoreName): void;
}) {
  useEffect(() => {
    initializeGroceryList(householdId);
    const interval = poll(() => fetchGroceries(householdId), POLL_INTERVAL);

    return () => {
      clearInterval(interval);
    };
  }, [householdId, initializeGroceryList, fetchGroceries]);

  return (
    <Stack
      height="100%"
      display="flex"
      justifyContent="space-between"
      padding="24px"
      boxSizing="border-box"
    >
      <Stack maxHeight="90%">
        <Stack>
          <Box
            display="flex"
            justifyContent="space-between"
            alignItems="center"
          >
            <Typography level="h3">Groceries</Typography>
            {endIcons}
          </Box>
          <Box
            component="form"
            onSubmit={createGroceryItem(householdId)}
            pt="8px"
            boxSizing="border-box"
          >
            <Input
              value={groceryItemText}
              onChange={handleGroceryItemTextChange}
              placeholder="Add a grocery item or recipe url!"
            />
          </Box>
        </Stack>
        <Box overflow="scroll" width="100%">
          <GroceryList
            groceries={groceryList.items}
            layout={groceryList.layout}
            checkGroceryItem={checkGroceryItem}
          />
        </Box>
      </Stack>
      <Stack>
        <Button
          color="danger"
          variant="soft"
          fullWidth
          onClick={() => clearCheckedItems(householdId)}
          disabled={isUpdating}
          loading={isUpdating}
        >
          delete checked
        </Button>
        <MoreInfoSheet
          selectedStores={selectedStores}
          toggleStore={toggleStore}
        />
      </Stack>
    </Stack>
  );
}
