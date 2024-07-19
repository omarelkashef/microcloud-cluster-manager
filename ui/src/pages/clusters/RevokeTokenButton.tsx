import { FC } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { ActionButton } from "@canonical/react-components";
import { deleteToken } from "api/tokens";
import { Token } from "types/token";
import { queryKeys } from "util/queryKeys";

interface Props {
  token: Token;
}

const RevokeTokenButton: FC<Props> = ({ token }) => {
  const queryClient = useQueryClient();

  const handleDeleteToken = async () => {
    await deleteToken(token.site_name)
      .then(() => {
        void queryClient.invalidateQueries({
          queryKey: [queryKeys.tokens],
        });
      })
      .catch((e: Error) => {
        if (e.message === "Unable to delete Token") {
          return;
        }
      });
  };

  return (
    <ActionButton
      className="u-no-margin"
      appearance="negative"
      onClick={() => void handleDeleteToken()}
    >
      Revoke
    </ActionButton>
  );
};

export default RevokeTokenButton;
