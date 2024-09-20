import { CloudUploadOutlined } from "@mui/icons-material";
import { Button, styled } from "@mui/joy";
import { ChangeEvent } from "react";

type FileUploadButtonProps =
  | {
      loading?: boolean;
      multiple: false | undefined;
      onChange(file: File): void;
    }
  | {
      loading?: boolean;
      multiple: true;
      onChange(fileList: FileList): void;
    };

const VisuallyHiddenInput = styled("input")`
  clip: rect(0 0 0 0);
  clip-path: inset(50%);
  height: 1px;
  overflow: hidden;
  position: absolute;
  bottom: 0;
  left: 0;
  white-space: nowrap;
  width: 1px;
`;

export default function FileUploadButton({
  loading,
  onChange,
  multiple,
}: FileUploadButtonProps) {
  function handleFileChange(event: ChangeEvent<HTMLInputElement>) {
    const files = event.target.files;

    if (!files?.length) {
      return;
    }

    if (multiple) {
      onChange(files);
    } else {
      onChange(files[0]);
    }
  }

  return (
    <Button
      component="label"
      role={undefined}
      tabIndex={-1}
      variant="outlined"
      color="neutral"
      loading={loading}
    >
      <CloudUploadOutlined />
      <VisuallyHiddenInput
        type="file"
        multiple={multiple}
        onChange={handleFileChange}
      />
    </Button>
  );
}
