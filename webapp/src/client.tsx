export type DrinkSize = 0 | 1;

export type OrderList = {
  id: string;
  created: Date;
  orders?: Order[];
};

export type StoreItem = {
  id: string;
  variants?: string[];
  dips?: string[];
};

export type Drink = {
  name: string;
  size: DrinkSize;
};

export type Order = CreateOrder & {
  id: string;
  created: Date;
};

export type CreateOrder = {
  creator: string;
  storeItem: StoreItem;
  drink?: Drink;
};

export type ShopData = {
  categories: ShopCategory[];
  drinks: ShopDrinkItem[];
};

export type ShopCategory = {
  id: string;
  name: string;
  items?: ShopStoreItem[];
};

export type ShopVariant = {
  name: string;
  description: string;
};

export type ShopStoreItem = {
  id: string;
  title: string;
  description: string;
  price: string;
  variants?: ShopVariant[];
  dips?: string[];
};

export type ShopDrinkItem = {
  name: string;
  description: string;
  price: string;
};

export class Client {
  constructor(private endpoint: string) {}

  private async req<R>(
    method: string,
    path: string,
    body?: object,
  ): Promise<R> {
    const res = await window.fetch(`${this.endpoint}${path}`, {
      method,
      body: body ? JSON.stringify(body) : undefined,
    });
    if (!res.ok) {
      throw new Error(`request failed with status code ${res.status}`);
    }
    if (res.status === 204) {
      return {} as R;
    }
    return await res.json();
  }

  async getShopData(): Promise<ShopData> {
    return this.req("GET", "/items");
  }

  async createList(): Promise<OrderList> {
    return this.req("POST", "/lists");
  }

  async getList(id: string): Promise<OrderList> {
    return this.req("GET", `/lists/${id}`);
  }

  async deleteList(id: string): Promise<OrderList> {
    return this.req("DELETE", `/lists/${id}`);
  }

  async createOrder(listId: string, order: CreateOrder): Promise<Order> {
    return this.req("POST", `/lists/${listId}/orders`, order);
  }
}

const ROOT_URL =
  import.meta.env.VITE_API_ROOT_URL ??
  (import.meta.env.PROD ? "/api" : "http://localhost:8080/api");

export const CLIENT = new Client(ROOT_URL); // TODO: Replace with env var
