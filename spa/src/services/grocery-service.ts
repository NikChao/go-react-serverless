import { ApiService } from "./api-service";

export interface GroceryItem {
  householdId: string;
  id: string;
  name: string;
  checked: boolean;
}

type LayoutBlockType = "GroceryItemId" | "Text";
export interface LayoutBlock {
  value: string;
  type: LayoutBlockType;
}

export interface GroceryList {
  items: GroceryItem[];
  layout: LayoutBlock[];
}

export interface GroceryMagicRequest {
  groceryList: GroceryList;
  householdId: string;
  preferredStores: StoreName[];
}

export interface GroceryMagicResponse {
  groceryList: GroceryList;
}

export interface BatchDeleteGroceryItemsRequest {
  itemsToDelete: GroceryItem[];
}

export type StoreName = "coles" | "aldi" | "woolies" | "sam cocos";
export const ALL_STORES: StoreName[] = [
  "aldi",
  "coles",
  "woolies",
  "sam cocos",
];
export interface StoreData {
  itemName: string;
  price: string;
  lastUpdated: string;
  storeName: StoreName;
}

export interface GroceryItemData {
  name: string;
  storeData: StoreData[];
}

export class GroceryService {
  constructor(private readonly apiService: ApiService) {}

  public getGroceryList(householdId: string): Promise<GroceryList> {
    return this.apiService.get<GroceryList>(`/groceries/${householdId}`);
  }

  public createGroceryItem(name: string, householdId: string): Promise<void> {
    const groceryItem: GroceryItem = {
      id: "",
      name,
      householdId,
      checked: false,
    };

    return this.apiService.put("/groceries", groceryItem);
  }

  public updateGroceryItem(groceryItem: GroceryItem): Promise<void> {
    return this.apiService.post("/groceries", groceryItem);
  }

  public magic(
    householdId: string,
    groceryList: GroceryList,
    preferredStores: StoreName[]
  ): Promise<GroceryMagicResponse> {
    const request: GroceryMagicRequest = {
      householdId,
      groceryList,
      preferredStores,
    };

    return this.apiService.post<GroceryMagicResponse>(
      "/groceries/magic",
      request
    );
  }

  public async clearCheckedGroceryItems(
    groceries: GroceryItem[]
  ): Promise<void> {
    const itemsToDelete = groceries.filter(({ checked }) => checked);

    if (!itemsToDelete.length) {
      return;
    }

    const request: BatchDeleteGroceryItemsRequest = {
      itemsToDelete,
    };

    return this.apiService.post("/groceries/batchDelete", request);
  }
}
