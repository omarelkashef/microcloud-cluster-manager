import { FC, useState } from "react";
import { Button, Form, Input } from "@canonical/react-components";
import { ConfigField } from "types/config";
import ConfigFieldDescription from "./ConfigFieldDescription";

interface Props {
  initialValue: string;
  configField: ConfigField;
  onSubmit: (newValue: string | boolean) => void;
  onCancel: () => void;
}

const getConfigId = (key: string) => {
  return key.replace(".", "___");
};

const SettingFormInput: FC<Props> = ({
  initialValue,
  configField,
  onSubmit,
  onCancel,
}) => {
  const [value, setValue] = useState(initialValue);

  const getInputType = () => {
    return "text";
  };

  return (
    <Form
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit(value);
      }}
    >
      {/* hidden submit to enable enter key in inputs */}
      <Input type="submit" hidden value="Hidden input" />
      <Input
        aria-label={configField.key}
        id={getConfigId(configField.key)}
        wrapperClassName="input-wrapper"
        type={getInputType()}
        value={String(value)}
        onChange={(e) => setValue(e.target.value)}
        help={
          <ConfigFieldDescription
            description={configField.longdesc}
            className="p-form-help-text"
          />
        }
      />
      <Button appearance="base" onClick={onCancel}>
        Cancel
      </Button>
      <Button appearance="positive" type="submit">
        Save
      </Button>
    </Form>
  );
};

export default SettingFormInput;
