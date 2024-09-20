import { observer } from "mobx-react-lite";
import { GroceryScreen } from "./screen";
import { UserStore } from "../../store/user-store";
import { GroceryListStore } from "./store";

export function createGroceryScreen(
  store: GroceryListStore,
  userStore: UserStore,
  endIcons: React.ReactNode
) {
  return observer(() => (
    <GroceryScreen
      selectedStores={store.selectedStores}
      toggleStore={store.toggleStore}
      endIcons={endIcons}
      groceryItemText={store.groceryItemText}
      isUpdating={store.isUpdating}
      groceryList={store.groceryList}
      householdId={userStore.effectiveHouseholdId}
      handleGroceryItemTextChange={store.handleGroceryItemTextChange}
      checkGroceryItem={store.checkGroceryItem}
      clearCheckedItems={store.clearCheckedItems}
      createGroceryItem={store.createGroceryItem}
      fetchGroceries={store.fetchGroceriesIfTabFocussed}
      initializeGroceryList={store.initializeGroceryList}
    />
  ));
}
