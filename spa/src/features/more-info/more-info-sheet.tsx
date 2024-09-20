import { Button, Drawer, Stack, Typography } from "@mui/joy";
import { useState } from "react";
import { Link } from "@tanstack/react-router";
import { ALL_STORES, StoreName } from "../../services/grocery-service";
import { Checkbox } from "@mui/joy";
import { observer } from "mobx-react-lite";

interface MoreInfoSheetProps {
  selectedStores: StoreName[];
  toggleStore(storeName: StoreName): void;
}

export const MoreInfoSheet = observer(
  ({ selectedStores, toggleStore }: MoreInfoSheetProps) => {
    const [isOpen, setIsOpen] = useState(false);

    function open() {
      setIsOpen(true);
    }

    function close() {
      setIsOpen(false);
    }

    return (
      <>
        <Button color="primary" variant="plain" onClick={open}>
          more info
        </Button>
        <Drawer open={isOpen} onClose={close} anchor="bottom">
          <Stack p="16px" height="70vh" gap="16px">
            <Stack>
              <Typography level="body-md" fontWeight="600">
                More info
              </Typography>
              <Typography level="body-sm">
                Store recommendations are made firstly based on lowest price,
                and if no pricing data exists, then on my personal preference
                (i.e. I like the baked goods at coles). Price data is collected
                weekly online and through receipts sent in by email.{" "}
                <Link to="/data" size={12}>
                  raw price data can be seen here.
                </Link>
              </Typography>
            </Stack>
            <Stack gap={2}>
              <Typography level="body-sm" fontWeight="600">
                What stores do you shop at?
              </Typography>
              {ALL_STORES.map((store) => (
                <Checkbox
                  key={store}
                  name={store}
                  label={store}
                  onChange={() => toggleStore(store)}
                  checked={selectedStores.includes(store)}
                />
              ))}
            </Stack>
          </Stack>
        </Drawer>
      </>
    );
  }
);
