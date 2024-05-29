import { Button, CircularProgress } from "@chakra-ui/react";

interface MineBlockButtonProps {
  isPending: boolean;
  onClick: () => void;
}

export default function MineBlockButton({
  isPending,
  onClick,
}: MineBlockButtonProps) {
  return isPending ? (
    <Button>
      Mining block...
      <CircularProgress
        ml={"10px"}
        size={"20px"}
        isIndeterminate
        color="green.300"
      />
    </Button>
  ) : (
    <Button onClick={onClick}>Mine block</Button>
  );
}
