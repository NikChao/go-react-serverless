import { makeAutoObservable } from "mobx";
import {
  ALL_STORES,
  GroceryItem,
  GroceryList,
  GroceryService,
  StoreName,
} from "../../services/grocery-service";
import { ChangeEvent, FormEvent } from "react";
import { UserStore } from "../../store/user-store";

const MAGIC_ENABLED_STORAGE_KEY = "MAGIC_ENABLED";
export const RAW_DATA_SELECTED_STORES_KEY = "RAW_DATA_SELECTED_STORES";

export class GroceryListStore {
  isFetching = false;
  isUpdating = false;
  magicEnabled = true;
  groceryItemText: string = "";
  groceryList: GroceryList = { items: [], layout: [] };
  selectedStores: StoreName[] = [];

  constructor(
    private readonly groceryService: GroceryService,
    private readonly userStore: UserStore,
    private readonly storage: Storage = window.localStorage
  ) {
    makeAutoObservable(this);

    const storedMagicEnabled = this.storage.getItem(MAGIC_ENABLED_STORAGE_KEY);
    if (storedMagicEnabled) {
      this.magicEnabled = Boolean(storedMagicEnabled);
    }

    this.selectedStores = this.getStoredSelectedStores();
  }

  initializeGroceryList = (householdId: string) => {
    this.fetchGroceryList(householdId);
  };

  fetchGroceryList = async (householdId: string) => {
    if (this.isFetching) {
      return;
    }

    try {
      this.isFetching = true;
      let groceryList = await this.groceryService.getGroceryList(householdId);

      if (this.magicEnabled) {
        groceryList = (
          await this.groceryService.magic(
            householdId,
            groceryList,
            this.selectedStores
          )
        ).groceryList;
      }

      this.groceryList = groceryList;
    } finally {
      this.isFetching = false;
    }
  };

  fetchGroceriesIfTabFocussed = async (householdId: string) => {
    if (document.hasFocus()) {
      await this.fetchGroceryList(householdId);
    }
  };

  magic = async (householdId: string) => {
    this.magicEnabled = !this.magicEnabled;
    this.storage.setItem(
      MAGIC_ENABLED_STORAGE_KEY,
      this.magicEnabled.toString()
    );

    this.fetchGroceryList(householdId);
  };

  createGroceryItem =
    (householdId: string) => async (e: FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      if (!this.groceryItemText) {
        return;
      }

      await this.groceryService.createGroceryItem(
        this.groceryItemText,
        householdId
      );
      await this.fetchGroceryList(householdId);
      this.groceryItemText = "";
    };

  addGroceryItems = async (groceryItems: string[]) => {
    const householdId = this.userStore.effectiveHouseholdId;

    const requests = groceryItems.map((itemText) =>
      this.groceryService.createGroceryItem(itemText, householdId)
    );
    return Promise.all(requests);
  };

  checkGroceryItem = async (id: string) => {
    this.groceryList = {
      ...this.groceryList,
      items: this.groceryList.items.map((item: GroceryItem) => {
        if (item.id === id) {
          item.checked = !item.checked;
          this.groceryService.updateGroceryItem(item);
        }

        return item;
      }),
    };
  };

  clearCheckedItems = async (householdId: string) => {
    this.isUpdating = true;
    await this.groceryService.clearCheckedGroceryItems(this.groceryList.items);
    await this.fetchGroceryList(householdId);
    this.isUpdating = false;
  };

  handleGroceryItemTextChange = (e: ChangeEvent<HTMLInputElement>) => {
    this.groceryItemText = e.target.value;
  };

  getStoredSelectedStores = (): StoreName[] => {
    const storedSelectedStores = localStorage.getItem(
      RAW_DATA_SELECTED_STORES_KEY
    );
    return storedSelectedStores ? JSON.parse(storedSelectedStores) : ALL_STORES;
  };

  toggleStore = (store: StoreName): void => {
    if (this.selectedStores.includes(store)) {
      this.selectedStores = this.selectedStores.filter(
        (storeName) => storeName != store
      );
    } else {
      this.selectedStores.push(store);
    }

    localStorage.setItem(
      RAW_DATA_SELECTED_STORES_KEY,
      JSON.stringify(this.selectedStores)
    );

    this.fetchGroceryList(this.userStore.effectiveHouseholdId);
  };
}
