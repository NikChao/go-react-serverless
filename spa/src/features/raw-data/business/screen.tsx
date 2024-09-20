import { Typography, Stack, Snackbar, Box } from "@mui/joy";
import { Snack } from "../../../models/snack";
import FileUploadButton from "../../../components/file-upload-button";

interface RawDataForBusinessProps {
  isUploadingReceipt: boolean;
  handleReceiptUpload(fileList: FileList): void;
  snacks: Snack[];
}

export function RawDataForBusiness({
  snacks,
  isUploadingReceipt,
  handleReceiptUpload,
}: RawDataForBusinessProps) {
  return (
    <>
      {snacks.map(toSnackBar)}
      <Stack padding="16px">
        <Typography fontStyle="italic" fontWeight="600" level="body-md">
          TaskTote for business coming soon...
        </Typography>
        <Box display="flex" alignItems="center" gap={2}>
          <Stack>
            <Typography level="body-xs">
              Want to support the project?
            </Typography>
            <Typography level="body-xs">Upload a receipt!</Typography>
          </Stack>
          <FileUploadButton
            multiple
            loading={isUploadingReceipt}
            onChange={handleReceiptUpload}
          />
        </Box>
      </Stack>
    </>
  );
}

function toSnackBar(snack: Snack) {
  return (
    <Snackbar
      key={snack.id}
      open={snack.open}
      onClose={snack.onClose}
      color={snack.color}
      variant="soft"
      autoHideDuration={1500}
      anchorOrigin={{ vertical: "top", horizontal: "center" }}
    >
      <Typography level="body-md">{snack.text}</Typography>
    </Snackbar>
  );
}
