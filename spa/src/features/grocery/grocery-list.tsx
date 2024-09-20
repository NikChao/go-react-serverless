import {
  List,
  ListItem,
  ListItemButton,
  Checkbox,
  Typography,
  Sheet,
} from "@mui/joy";
import { GroceryItem, LayoutBlock } from "../../services/grocery-service";

const containerStyles = {
  width: "100%",
  bgcolor: "background.paper",
};

export function GroceryList({
  groceries,
  layout,
  checkGroceryItem,
}: {
  groceries: GroceryItem[];
  layout: LayoutBlock[];
  checkGroceryItem(id: string): void;
}) {
  if (!groceries?.length) {
    return (
      <Sheet sx={{ backgroundColor: "white", paddingTop: "16px" }}>
        <Typography
          level="body-sm"
          textColor="neutral.400"
          fontStyle="italic"
          textAlign="center"
        >
          No items yet...
        </Typography>
      </Sheet>
    );
  }

  return (
    <List sx={containerStyles}>
      {layout.map(({ type, value }) => {
        return (
          <ListItem
            sx={{ width: "100%", height: "48px" }}
            key={value}
            onClick={() => {
              if (type !== "GroceryItemId") {
                return;
              }

              checkGroceryItem(value);
            }}
          >
            {type === "GroceryItemId" ? (
              <ListItemButton
                sx={{ display: "flex", justifyContent: "space-between" }}
              >
                <Typography>
                  {groceries.find(({ id }) => id === value)?.name}
                </Typography>
                <Checkbox
                  checked={
                    groceries.find(({ id }) => id === value)?.checked ?? false
                  }
                />
              </ListItemButton>
            ) : (
              <Typography level="h4">{value}</Typography>
            )}
          </ListItem>
        );
      })}
    </List>
  );
}
