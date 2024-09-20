import { useState } from "react";
import { InsertLink, AddHome, ArrowForward } from "@mui/icons-material";
import {
  IconButton,
  CircularProgress,
  Drawer,
  Box,
  Typography,
  Stack,
  Tooltip,
  Input,
  Button,
} from "@mui/joy";

interface InviteToHouseholdProps {
  householdId?: string;
  isLoading: boolean;
  leaveHousehold(): Promise<void>;
  joinHousehold(householdId: string): Promise<void>;
  createAndJoinHousehold(): Promise<void>;
}

export function CreateAndInviteToHousehold({
  householdId,
  isLoading,
  leaveHousehold,
  createAndJoinHousehold,
  joinHousehold,
}: InviteToHouseholdProps) {
  const [groceryListId, setGroceryListId] = useState("");
  const [groceryListIdHasError, setGroceryListIdHasError] = useState(false);
  const [isTooltipOpen, setIsTooltipOpen] = useState(false);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  let joinUrl = `${window.location.protocol}//${window.location.hostname}`;
  if (window.location.port) {
    joinUrl += `:${window.location.port}`;
  }
  joinUrl += `/households/join/${householdId}`;

  function copyJoinUrl() {
    setIsTooltipOpen(true);
    navigator.clipboard.writeText(joinUrl);

    setTimeout(() => {
      setIsTooltipOpen(false);
    }, 500);
  }

  function closeDrawer() {
    setIsDrawerOpen(false);
    setGroceryListId("");
    setGroceryListIdHasError(false);
  }

  function joinGroceryList() {
    const pattern = /\/households\/join\/([0-9a-fA-F-]{36})/;
    const parsedId = groceryListId.match(pattern)?.[1] ?? groceryListId;

    const uuidPattern =
      /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[4][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$/;
    const isValid = uuidPattern.test(parsedId);

    if (!isValid) {
      setGroceryListIdHasError(true);
      return;
    }

    joinHousehold(parsedId);
    closeDrawer();
  }

  function handleLeaveHousehold() {
    leaveHousehold();
    closeDrawer();
  }

  if (
    householdId !== null &&
    householdId !== "1" &&
    !householdId?.startsWith("USER")
  ) {
    return (
      <>
        <IconButton onClick={() => setIsDrawerOpen(true)}>
          <InsertLink />
        </IconButton>
        <Drawer anchor="bottom" open={isDrawerOpen} onClose={closeDrawer}>
          <Box
            height="100%"
            padding="24px"
            display="flex"
            flexDirection="column"
            justifyContent="center"
            alignItems="center"
          >
            <Stack
              height="100%"
              width="100%"
              px="24px"
              display="flex"
              flexDirection="column"
              justifyContent="space-between"
            >
              <Box />
              <Stack width="100%" gap={4}>
                <Stack>
                  <Typography level="body-sm">
                    Share this link to share your grocery list
                  </Typography>
                  <Box display="flex" alignItems="center">
                    <Input value={joinUrl} fullWidth />
                    <Tooltip
                      open={isTooltipOpen}
                      onClose={() => {}}
                      title="Copied!"
                      disableHoverListener
                    >
                      <IconButton onClick={copyJoinUrl}>
                        <InsertLink />
                      </IconButton>
                    </Tooltip>
                  </Box>
                </Stack>

                <Stack>
                  <Typography level="body-sm">
                    Or join someone elses!
                  </Typography>
                  <Box display="flex" alignItems="center">
                    <Input
                      placeholder="grocery list id"
                      value={groceryListId}
                      onChange={(e) => setGroceryListId(e.target.value)}
                      fullWidth
                      error={groceryListIdHasError}
                    />
                    <IconButton onClick={joinGroceryList}>
                      <ArrowForward />
                    </IconButton>
                  </Box>
                </Stack>
              </Stack>
              <Button
                variant="plain"
                color="danger"
                onClick={handleLeaveHousehold}
              >
                leave grocery list
              </Button>
            </Stack>
          </Box>
        </Drawer>
      </>
    );
  }

  if (isLoading) {
    return <CircularProgress size="sm" />;
  }

  return (
    <IconButton onClick={createAndJoinHousehold}>
      <AddHome />
    </IconButton>
  );
}
