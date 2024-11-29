import { FC } from "react";
import { MainTable, Row } from "@canonical/react-components";
import BaseLayout from "components/BaseLayout";
import { queryKeys } from "util/queryKeys";
import { useQuery } from "@tanstack/react-query";
import NotificationRow from "components/NotificationRow";
import Loader from "components/Loader";
import { fetchConfigurations } from "api/settings";
import { ConfigData } from "types/config";

const Settings: FC = () => {
  const { data: configurations, isLoading } = useQuery({
    queryKey: [queryKeys.configuration],
    queryFn: fetchConfigurations,
  });

  const headers = [
    { content: "Configuration", classNames: "title" },
    { content: "Description", classNames: "description" },
    { content: "Value", classNames: "value" },
  ];

  const configKeys = Object.keys(configurations || {});
  const rows = configKeys.map((key) => {
    const config = configurations?.[key] as ConfigData;
    return {
      columns: [
        {
          content: config.title,
          role: "cell",
          title: config.title,
          className: "u-truncate title",
          "aria-label": "Configuration",
        },
        {
          content: config.description,
          role: "cell",
          title: config.description,
          className: "u-truncate description",
          "aria-label": "Description",
        },
        {
          content: config.value,
          role: "cell",
          title: config.value,
          className: "u-truncate value",
          "aria-label": "Value",
        },
      ],
    };
  });

  if (isLoading) {
    return <Loader />;
  }

  return (
    <BaseLayout title="Settings">
      <Row>
        <NotificationRow />
        <div className="settings">
          <MainTable
            id="settings-table"
            headers={headers}
            rows={rows}
            emptyStateMsg="No settings to display"
          />
        </div>
      </Row>
    </BaseLayout>
  );
};

export default Settings;
