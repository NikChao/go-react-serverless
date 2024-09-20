export function chunk(fileList: FileList, chunkSize: number): File[][] {
  const files = [...fileList];

  return Array.from(
    { length: Math.ceil(files.length / chunkSize) },
    (_: any, i: number) => files.slice(i * chunkSize, i * chunkSize + chunkSize)
  );
}
