import { ApiService } from "./api-service";
import { GroceryItemData } from "./grocery-service";

export interface Catalog {
  data: GroceryItemData[];
}

export class CatalogService {
  constructor(private readonly apiService: ApiService) {}

  public getCatalog() {
    return this.apiService.get<Catalog>("/catalog");
  }
}
