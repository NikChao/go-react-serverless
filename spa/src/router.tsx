import {
  Outlet,
  createRootRoute,
  createRoute,
  createRouter,
  useNavigate,
} from "@tanstack/react-router";
import { createGroceryScreen } from "./features/grocery/create";
import { UserStore } from "./store/user-store";
import { CreateAndInviteToHousehold as CreateAndInviteToHouseholdImpl } from "./features/household/invite-to-household";
import { UserService } from "./services/user-service";
import { ApiService } from "./services/api-service";
import { GroceryService } from "./services/grocery-service";
import { HouseholdService } from "./services/household-service";
import { observer } from "mobx-react-lite";
import { Box, CircularProgress } from "@mui/joy";
import { AutoCompleteGroceries } from "./features/end-buttons/autocomplete-groceries";
import { GroceryListStore } from "./features/grocery/store";
import { RawDataScreen } from "./features/raw-data/screen";
import { CatalogService } from "./services/catalog-service";
import { ReceiptService } from "./services/receipt-service";
import { createRawDataForBusiness } from "./features/raw-data/business/create";

const apiService = new ApiService(window.fetch);
const groceryService = new GroceryService(apiService);
const catalogService = new CatalogService(apiService);
const receiptService = new ReceiptService(apiService);
const userService = new UserService(apiService);
const householdService = new HouseholdService(apiService);
const userStore = new UserStore(userService, householdService);
const groceryStore = new GroceryListStore(groceryService, userStore);

const CreateAndInviteToHousehold = observer(() => {
  if (!userStore.userId) {
    return <CircularProgress size="md" />;
  }

  return (
    <CreateAndInviteToHouseholdImpl
      householdId={userStore.effectiveHouseholdId}
      isLoading={userStore.isLoading}
      leaveHousehold={userStore.leaveHousehold}
      joinHousehold={userStore.joinHousehold}
      createAndJoinHousehold={userStore.createAndJoinHousehold}
    />
  );
});

const EndIcons = observer(() => (
  <Box display="flex" alignItems="center">
    <AutoCompleteGroceries
      isLoading={groceryStore.isFetching}
      householdId={userStore.effectiveHouseholdId}
      isMagicEnabled={groceryStore.magicEnabled}
      magic={groceryStore.magic}
    />
    <CreateAndInviteToHousehold />
  </Box>
));

const endIcons = <EndIcons />;

const GroceryScreen = createGroceryScreen(groceryStore, userStore, endIcons);

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: GroceryScreen,
});

const joinHouseholdRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/households/join/$householdId",
  component: function JoinHousehold() {
    /** @ts-ignore */
    const { householdId } = joinHouseholdRoute.useParams();
    const navigate = useNavigate();

    userStore.joinHousehold(householdId).then(() => navigate({ to: "/" }));

    return null;
  },
});

const rawDataRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/data",
  component: () => (
    <RawDataScreen
      catalogService={catalogService}
      receiptService={receiptService}
      groceryStore={groceryStore}
    />
  ),
});

const RawDataForBusiness = createRawDataForBusiness(
  catalogService,
  receiptService
);

const rawDataForBusinessRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/business/data",
  component: RawDataForBusiness,
});

/** @ts-ignore */
const routeTree = rootRoute.addChildren([
  indexRoute,
  joinHouseholdRoute,
  rawDataRoute,
  rawDataForBusinessRoute,
]);

export const router = createRouter({ routeTree });
