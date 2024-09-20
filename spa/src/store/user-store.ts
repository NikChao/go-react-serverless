import { makeAutoObservable } from "mobx";
import { HouseholdService } from "../services/household-service";
import { UserService } from "../services/user-service";

export class UserStore {
  private static readonly USER_ID_LOCALSTORAGE_KEY = "GROCERY_USER_ID";

  userId?: string;
  householdId?: string;
  isLoading = false;

  constructor(
    private readonly userService: UserService,
    private readonly householdService: HouseholdService
  ) {
    makeAutoObservable(this);
    this.getOrCreateUser();
  }

  public get effectiveHouseholdId() {
    if (this.householdId) {
      return this.householdId;
    }

    if (this.userId) {
      return `USER-${this.userId}`;
    }

    return "1";
  }

  public createAndJoinHousehold = async () => {
    this.isLoading = true;

    const userId = this.userId;
    if (!userId) {
      throw new Error("User not yet created");
    }

    const household = await this.householdService.createHousehold();
    await this.householdService.joinHousehold(userId, household.id);
    this.householdId = household.id;
    this.isLoading = false;
  };

  public joinHousehold = async (householdId: string) => {
    const userId = await this.getUserId();

    await this.householdService.joinHousehold(userId, householdId);
    this.householdId = householdId;
  };

  public leaveHousehold = async () => {
    if (!this.householdId) {
      return;
    }

    const userId = await this.getUserId();
    await this.householdService.leaveHousehold(userId, this.householdId);
    this.householdId = undefined;
    this.getOrCreateUser();
  };

  private createUser = async () => {
    this.isLoading = true;
    const user = await this.userService.createUser();

    if (!user.id) {
      throw new Error("UserId not created successfully");
    }

    localStorage.setItem(UserStore.USER_ID_LOCALSTORAGE_KEY, user.id);
    this.userId = user.id;
    this.isLoading = false;
  };

  private getOrCreateUser = async () => {
    const userId = localStorage.getItem(UserStore.USER_ID_LOCALSTORAGE_KEY);

    if (userId) {
      this.userId = userId;
      const user = await this.userService.getUser(userId);

      if (user.householdIds.length) {
        this.householdId = user.householdIds[0];
      }
    } else {
      this.createUser();
    }
  };

  private getUserId = async () => {
    let userId = this.userId;

    if (userId) {
      return userId;
    }

    await this.getOrCreateUser();
    userId = this.userId;

    if (!userId) {
      throw new Error("Failed to create user");
    }

    return userId;
  };
}
