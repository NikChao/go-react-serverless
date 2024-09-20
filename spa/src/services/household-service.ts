import { ApiService } from "./api-service";

export interface Household {
  id: string;
}

export class HouseholdService {
  constructor(private readonly apiService: ApiService) {}

  public createHousehold(): Promise<Household> {
    return this.apiService.put("/households");
  }

  public joinHousehold(userId: string, householdId: string): Promise<void> {
    return this.apiService.post(`/households/join/${householdId}/${userId}`);
  }

  public leaveHousehold(userId: string, householdId: string): Promise<void> {
    return this.apiService.post(`/households/leave/${householdId}/${userId}`);
  }
}
