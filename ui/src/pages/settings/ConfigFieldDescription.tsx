import { FC } from "react";
import { configDescriptionToHtml } from "util/config";

interface Props {
  description?: string;
  className?: string;
}

const ConfigFieldDescription: FC<Props> = ({ description, className }) => {
  return description ? (
    <p
      className={className}
      dangerouslySetInnerHTML={{
        __html: configDescriptionToHtml(description),
      }}
    ></p>
  ) : null;
};

export default ConfigFieldDescription;
