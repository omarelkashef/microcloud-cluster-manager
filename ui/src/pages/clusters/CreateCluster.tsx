import { FC } from "react";
import {
  ActionButton,
  Button,
  Col,
  Form,
  Row,
} from "@canonical/react-components";
import BaseLayout from "components/BaseLayout";
import { useNavigate } from "react-router-dom";
import { useFormik } from "formik";
import ClusterCreateDetailsForm, {
  CreateClusterFormValues,
  newTokenPayload,
} from "./ClusterCreateDetailsForm";
import FormFooterLayout from "components/forms/FormFooterLayout";
import { createToken } from "api/tokens";
import * as Yup from "yup";
import { queryKeys } from "util/queryKeys";
import { useQueryClient } from "@tanstack/react-query";
import { getDefaultExpiryDate } from "util/createCluster";

const CreateCluster: FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const ClusterSchema = Yup.object().shape({
    site_name: Yup.string().required("Site name is required"),
    expiry: Yup.date()
      .optional()
      .min(new Date(), "Expiry date cannot be in the past"),
  });

  const submit = (values: CreateClusterFormValues) => {
    const tokenPayload = newTokenPayload(values);

    createToken(JSON.stringify(tokenPayload))
      .then(() => {
        navigate(`/ui/sites/tokens`);
      })
      .catch((e: Error) => {
        if (e.message === "Unable to create new Token") {
          return;
        }
      })
      .finally(() => {
        void queryClient.invalidateQueries({
          queryKey: [queryKeys.tokens],
        });
      });
  };

  const formik = useFormik<CreateClusterFormValues>({
    initialValues: {
      site_name: "",
      expiry: getDefaultExpiryDate(),
    },
    validationSchema: ClusterSchema,
    onSubmit: submit,
  });

  return (
    <BaseLayout
      title="Create a new join Token"
      contentClassName="create-cluster"
    >
      <Form onSubmit={formik.handleSubmit} className="form">
        <Row className="form-contents">
          <Col size={12}>
            <ClusterCreateDetailsForm formik={formik} />
          </Col>
        </Row>
      </Form>
      <FormFooterLayout>
        <Button appearance="base" onClick={() => navigate("/")}>
          Cancel
        </Button>
        <ActionButton
          appearance="positive"
          loading={formik.isSubmitting}
          disabled={!formik.isValid}
          onClick={() => submit(formik.values)}
        >
          Create
        </ActionButton>
      </FormFooterLayout>
    </BaseLayout>
  );
};

export default CreateCluster;
