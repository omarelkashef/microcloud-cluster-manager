import {
  CheckboxInput,
  ConfirmationButton,
  useToastNotification,
} from "@canonical/react-components";
import { useQueryClient } from "@tanstack/react-query";
import { deleteClusterLinks } from "api/clusters";
import type { FC } from "react";
import { useState } from "react";
import DsIcon from "components/DsIcon";
import type { ClusterLink } from "types/cluster";
import type { ClusterConnection } from "util/clusterLinks";
import {
  canRemoveTargetEnd,
  hasTargetEnd,
  isBidirectional,
} from "util/clusterLinks";
import { pluralize } from "util/helpers";
import { queryKeys } from "util/queryKeys";

interface Props {
  cluster: string;
  links: ClusterConnection[];
  hasLabel?: boolean;
  disabled?: boolean;
  onStart?: () => void;
  onFinish?: () => void;
}

interface LinkDeletion {
  link: ClusterLink;
  targetToClean?: string;
}

const DeleteClusterLinksButton: FC<Props> = ({
  cluster,
  links,
  hasLabel = true,
  disabled = false,
  onStart,
  onFinish,
}) => {
  const queryClient = useQueryClient();
  const toastNotification = useToastNotification();
  const [isLoading, setLoading] = useState(false);
  const [deleteBothEnds, setDeleteBothEnds] = useState(true);

  const isSingle = links.length === 1;
  const linksWithTargetEnd = links.filter(({ link }) => hasTargetEnd(link));
  const cleanableLinks = links.filter(({ link, target }) =>
    canRemoveTargetEnd(link, target),
  );
  const canCleanTargets = cleanableLinks.length > 0;
  const uncleanableCount = linksWithTargetEnd.length - cleanableLinks.length;
  const removesTargetEnds = deleteBothEnds && canCleanTargets;
  const checkboxLabel = linksWithTargetEnd.some(({ link }) =>
    isBidirectional(link),
  )
    ? "Delete links on both source and target cluster"
    : "Also delete identity on the target cluster";
  const title = isSingle ? "Delete cluster link" : "Delete cluster links";
  const iconLabel = isSingle
    ? `Delete cluster link ${links[0].link.name}`
    : title;
  const linkCount = isSingle ? (
    <>
      cluster link <strong>{links[0].link.name}</strong>
    </>
  ) : (
    <strong>
      {links.length} {pluralize("cluster link", links.length)}
    </strong>
  );

  const getSuccessMessage = () => {
    if (!isSingle) {
      return <>Deleted {linkCount}.</>;
    }
    const { link, target } = links[0];
    if (removesTargetEnds && target) {
      return (
        <>
          Deleted cluster link between <strong>{cluster}</strong> and{" "}
          <strong>{target.name}</strong>.
        </>
      );
    }
    return (
      <>
        Deleted cluster link <strong>{link.name}</strong>.
      </>
    );
  };

  const invalidateAffectedQueries = async (
    deletions: LinkDeletion[],
  ): Promise<void> => {
    const linkClusters = new Set([cluster]);
    const identityClusters = new Set<string>();

    for (const { link, targetToClean } of deletions) {
      // Only bidirectional links have an identity on the source cluster.
      if (isBidirectional(link)) {
        identityClusters.add(cluster);
      }
      if (targetToClean !== undefined) {
        identityClusters.add(targetToClean);
        if (isBidirectional(link)) {
          linkClusters.add(targetToClean);
        }
      }
    }

    await Promise.all([
      ...Array.from(linkClusters, async (clusterName) =>
        queryClient.invalidateQueries({
          queryKey: [queryKeys.clusters, clusterName, queryKeys.links],
        }),
      ),
      ...Array.from(identityClusters, async (clusterName) =>
        queryClient.invalidateQueries({
          queryKey: [queryKeys.identities, clusterName],
        }),
      ),
    ]);
  };

  const handleDelete = async (): Promise<void> => {
    setLoading(true);
    onStart?.();
    const deletions = links.map(({ link, target }) => ({
      link,
      targetToClean:
        deleteBothEnds && canRemoveTargetEnd(link, target)
          ? target?.name
          : undefined,
    }));

    try {
      await deleteClusterLinks(
        cluster,
        deletions.map(({ link, targetToClean }) => ({
          name: link.name,
          removeReverseEdge: targetToClean !== undefined,
        })),
      );
      toastNotification.success(getSuccessMessage());
    } catch (error) {
      toastNotification.failure(
        isSingle
          ? `Unable to delete cluster link ${links[0].link.name}.`
          : "Cluster link deletion failed.",
        error,
      );
    } finally {
      setDeleteBothEnds(true);
      setLoading(false);
      onFinish?.();
      void invalidateAffectedQueries(deletions);
    }
  };

  const getUncleanableNote = (): string | null => {
    if (uncleanableCount === 0) {
      return null;
    }
    if (isSingle) {
      return "The target cluster is not managed by cluster manager or has no registered tunnel. Only the source link is deleted.";
    }
    if (uncleanableCount === 1) {
      return "1 of the selected links has a target cluster that is not managed by cluster manager or has no registered tunnel. Only its source link is deleted.";
    }
    return `${uncleanableCount} of the selected links have target clusters that are not managed by cluster manager or have no registered tunnel. Only their source links are deleted.`;
  };
  const uncleanableNote = getUncleanableNote();

  return (
    <ConfirmationButton
      className="u-no-margin--bottom has-icon"
      appearance={hasLabel ? "default" : "base"}
      loading={isLoading}
      disabled={disabled || undefined}
      aria-label={hasLabel ? undefined : iconLabel}
      onHoverText={hasLabel ? undefined : "Delete"}
      confirmationModalProps={{
        title,
        children: (
          <>
            <p>
              This will permanently delete {linkCount}. This action cannot be
              undone.
            </p>
            {linksWithTargetEnd.length > 0 && (
              <CheckboxInput
                id="delete-cluster-links-both-ends"
                label={checkboxLabel}
                checked={removesTargetEnds}
                disabled={!canCleanTargets}
                onChange={() => {
                  setDeleteBothEnds((value) => !value);
                }}
              />
            )}
            {uncleanableNote && <p>{uncleanableNote}</p>}
          </>
        ),
        confirmButtonLabel: "Delete",
        close: () => {
          setDeleteBothEnds(true);
        },
        onConfirm: () => void handleDelete(),
      }}
      shiftClickEnabled
    >
      <DsIcon icon="delete" />
      {hasLabel && (
        <span>Delete {pluralize("cluster link", links.length)}</span>
      )}
    </ConfirmationButton>
  );
};

export default DeleteClusterLinksButton;
