import crypto from "crypto";

export const generateUUID = (): string => {
  return crypto.randomBytes(16).toString("hex");
};
