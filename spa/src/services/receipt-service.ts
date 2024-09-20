import { ApiService } from "./api-service";

const MAX_CONTENT_LENGTH = 10 * 1024 * 1024; // 10MB

export type ReceiptContext = "Business" | "Retail";

export class ReceiptService {
  constructor(private readonly apiService: ApiService) {}

  async uploadBusinessReceiptToPresignedUrl(file: File) {
    return this.uploadReceiptToPresignedUrl(file, "Business");
  }

  async uploadReceiptToPresignedUrl(
    file: File,
    receiptContext: ReceiptContext = "Retail"
  ) {
    const contentLength = file.size;

    if (contentLength > MAX_CONTENT_LENGTH) {
      throw new Error("File too big");
    }

    const { url } = await this.apiService.post<{ url: string }>(
      "/receipt/upload",
      {
        fileName: file.name,
        contentLength,
        receiptContext,
      }
    );

    if (!url) {
      return;
    }

    const uploadResponse = await fetch(url, {
      method: "PUT",
      headers: {
        "Content-Length": contentLength.toString(),
      },
      body: file,
    });

    if (!uploadResponse.ok) {
      throw new Error(`Failed to upload file ${file.name}`);
    }
  }
}
