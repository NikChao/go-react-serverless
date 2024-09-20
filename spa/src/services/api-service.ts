export class ApiService {
  private static readonly BASE_URL = process.env.REACT_APP_API_BASE_URL;

  constructor(private readonly fetch: typeof window.fetch) {}

  public get<T>(url: string): Promise<T> {
    return this.makeRequest<T>(url, "GET");
  }

  public delete<T>(url: string): Promise<T> {
    return this.makeRequest<T>(url, "DELETE");
  }

  public post<T = void>(url: string, body?: any): Promise<T> {
    return this.makeRequest(url, "POST", body);
  }

  public put<T = void>(url: string, body?: any): Promise<T> {
    return this.makeRequest(url, "PUT", body);
  }

  private async makeRequest<T = void>(
    url: string,
    method: string,
    body?: any
  ): Promise<T> {
    const options: RequestInit = {
      method,
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "same-origin",
    };

    if (body) {
      options.body = JSON.stringify(body);
    }

    const res = await fetch(`${ApiService.BASE_URL}${url}`, options);

    return res.json();
  }
}
