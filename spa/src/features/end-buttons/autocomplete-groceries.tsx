import { AutoFixHigh, AutoFixOff } from "@mui/icons-material";
import { IconButton } from "@mui/joy";

export function AutoCompleteGroceries({
  householdId,
  isLoading,
  isMagicEnabled,
  magic,
}: {
  householdId: string;
  isLoading: boolean;
  isMagicEnabled: boolean;
  magic(householdId: string): Promise<void>;
}) {
  return (
    <IconButton onClick={() => magic(householdId)} loading={isLoading}>
      {isMagicEnabled ? <AutoFixOff /> : <AutoFixHigh />}
    </IconButton>
  );
}
