import { ApiService } from "./api-service";

export interface User {
  id: string;
  householdIds: string[];
}

export class UserService {
  constructor(private readonly apiService: ApiService) {}

  public getUser(id: string): Promise<User> {
    return this.apiService.get(`/users/${id}`);
  }

  public createUser(): Promise<User> {
    return this.apiService.put("/users");
  }
}
